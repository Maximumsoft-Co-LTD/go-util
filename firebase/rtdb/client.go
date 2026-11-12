// Package rtdb is a minimal Firebase Realtime Database REST client that
// attaches a Google OAuth2 access token (ADC: GKE Workload Identity,
// Workload Identity Federation external_account config, or a key file)
// as an Authorization: Bearer header.
//
// Rollout is controlled by FIREBASE_RTDB_AUTH:
//
//	off     (default) send requests exactly as before, no Authorization header
//	prefer  attach a token when one is available, otherwise send unauthenticated and log
//	require never send without a token (returns ErrNoToken)
//
// Source of truth: shared/libs-templates/go-util/firebase/rtdb.
// Copied into service repos as internal/rtdb (Firebase leak 2026 / impl-p1-withdraw-rtdb).
// Keep copies identical; change the canonical package first.
//
// The package never logs tokens, request bodies or response bodies.
package rtdb

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuth scopes required by the Firebase RTDB REST API.
const (
	ScopeEmail = "https://www.googleapis.com/auth/userinfo.email"
	ScopeRTDB  = "https://www.googleapis.com/auth/firebase.database"
)

// EnvAuthMode is the environment variable read by FromEnv / ModeFromEnv.
const EnvAuthMode = "FIREBASE_RTDB_AUTH"

const (
	defaultTimeout = 10 * time.Second
	maxRespBytes   = 1 << 20
	maxSegmentLen  = 768
	logEvery       = time.Minute
)

// ErrNoToken is returned in Require mode when no access token is available.
var ErrNoToken = errors.New("rtdb: no access token (FIREBASE_RTDB_AUTH=require)")

// Mode selects how the client authenticates.
type Mode int

const (
	Off     Mode = iota // legacy: no Authorization header
	Prefer              // attach token when available, else send legacy + log
	Require             // refuse to send without a token
)

func (m Mode) String() string {
	switch m {
	case Prefer:
		return "prefer"
	case Require:
		return "require"
	default:
		return "off"
	}
}

// ParseMode parses off|prefer|require (case-insensitive). Empty means Off.
func ParseMode(s string) (Mode, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "off":
		return Off, nil
	case "prefer":
		return Prefer, nil
	case "require":
		return Require, nil
	}
	return Off, fmt.Errorf("rtdb: invalid %s value %q (want off|prefer|require)", EnvAuthMode, s)
}

// Client talks to one RTDB instance.
type Client struct {
	base string
	hc   *http.Client
	ts   oauth2.TokenSource
	mode Mode
	logf func(string, ...interface{})

	mu      sync.Mutex
	lastLog map[string]time.Time
}

// Option customises a Client (mainly for tests).
type Option func(*Client)

// WithHTTPClient overrides the HTTP client used for RTDB requests.
func WithHTTPClient(hc *http.Client) Option { return func(c *Client) { c.hc = hc } }

// WithTransport replaces only the RoundTripper, keeping the default timeout
// and redirect policy (e.g. for tests or an egress proxy).
func WithTransport(rt http.RoundTripper) Option {
	return func(c *Client) {
		hc := *c.hc
		hc.Transport = rt
		c.hc = &hc
	}
}

// WithTokenSource uses ts instead of Application Default Credentials.
func WithTokenSource(ts oauth2.TokenSource) Option { return func(c *Client) { c.ts = ts } }

// WithLogger overrides the logger (default log.Printf). Messages never contain tokens or bodies.
func WithLogger(f func(string, ...interface{})) Option { return func(c *Client) { c.logf = f } }

// New builds a client for base (https://<db>.firebaseio.com or
// https://<db>.<region>.firebasedatabase.app). In Prefer/Require mode it
// resolves Application Default Credentials unless WithTokenSource is given.
func New(ctx context.Context, base string, mode Mode, opts ...Option) (*Client, error) {
	b, err := validateBase(base)
	if err != nil {
		return nil, err
	}
	c := &Client{
		base:    b,
		hc:      &http.Client{Timeout: defaultTimeout, CheckRedirect: safeRedirect},
		mode:    mode,
		logf:    log.Printf,
		lastLog: map[string]time.Time{},
	}
	for _, o := range opts {
		o(c)
	}
	if mode == Off {
		c.ts = nil
		return c, nil
	}
	if c.ts == nil {
		// Bound the token fetch (metadata / STS / IAM credentials) so a hung
		// endpoint cannot block callers forever.
		tctx := context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Timeout: defaultTimeout})
		creds, cerr := google.FindDefaultCredentials(tctx, ScopeEmail, ScopeRTDB)
		if cerr != nil {
			if mode == Require {
				return nil, fmt.Errorf("rtdb: no credentials: %w", cerr)
			}
			c.logf("rtdb_auth_mode=%s credentials_unavailable", mode)
			return c, nil
		}
		c.ts = creds.TokenSource
	}
	c.ts = oauth2.ReuseTokenSource(nil, c.ts)
	c.logf("rtdb_auth_mode=%s", mode)
	return c, nil
}

// ModeFromEnv reads FIREBASE_RTDB_AUTH.
func ModeFromEnv() (Mode, error) { return ParseMode(os.Getenv(EnvAuthMode)) }

// FromEnv builds a client with base = $urlEnv (or def when unset) and mode = $FIREBASE_RTDB_AUTH.
func FromEnv(ctx context.Context, urlEnv, def string, opts ...Option) (*Client, error) {
	base := strings.TrimSpace(os.Getenv(urlEnv))
	if base == "" {
		base = def
	}
	mode, err := ModeFromEnv()
	if err != nil {
		return nil, err
	}
	return New(ctx, base, mode, opts...)
}

// Mode reports the configured auth mode.
func (c *Client) Mode() Mode { return c.mode }

// Do sends method to <base>/<segs...>.json. segs must not include ".json".
// body (if non-nil) is JSON-encoded with encoding/json. It returns the HTTP
// status and the response body (also for non-2xx, capped at 1 MiB); err is
// non-nil for transport errors, invalid input, missing token in Require mode
// and non-2xx responses.
func (c *Client) Do(ctx context.Context, method string, segs []string, body interface{}) (int, []byte, error) {
	if len(segs) == 0 {
		return 0, nil, errors.New("rtdb: empty path")
	}
	for _, s := range segs {
		if err := validSegment(s); err != nil {
			return 0, nil, err
		}
	}
	var rd io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return 0, nil, fmt.Errorf("rtdb: encode body: %w", err)
		}
		rd = bytes.NewReader(b)
	}
	esc := make([]string, len(segs))
	for i, s := range segs {
		esc[i] = url.PathEscape(s)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.base+"/"+strings.Join(esc, "/")+".json", rd)
	if err != nil {
		return 0, nil, fmt.Errorf("rtdb: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	if c.ts != nil {
		tok, terr := c.ts.Token()
		switch {
		case terr == nil && tok.AccessToken != "":
			tok.SetAuthHeader(req)
		case c.mode == Require:
			return 0, nil, ErrNoToken
		default:
			c.logOnce("rtdb_token_unavailable", "rtdb_token_unavailable mode=%s", c.mode)
		}
	} else if c.mode == Require {
		return 0, nil, ErrNoToken
	}

	res, err := c.hc.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("rtdb: %s request failed: %w", method, redactURLErr(err))
	}
	defer res.Body.Close()
	resp, _ := io.ReadAll(io.LimitReader(res.Body, maxRespBytes))
	if res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusForbidden {
		c.logf("rtdb_permission_denied status=%d method=%s path_depth=%d", res.StatusCode, method, len(segs))
	}
	if res.StatusCode/100 != 2 {
		return res.StatusCode, resp, fmt.Errorf("rtdb: HTTP %d", res.StatusCode)
	}
	return res.StatusCode, resp, nil
}

// SplitPath converts a legacy REST path such as "svc/svc_123.json" into
// segments for Do. A leading "/" and a trailing ".json" are removed. Query
// strings are not supported.
func SplitPath(p string) ([]string, error) {
	if strings.ContainsAny(p, "?#") {
		return nil, errors.New("rtdb: query/fragment not supported in path")
	}
	p = strings.TrimSuffix(strings.TrimPrefix(p, "/"), ".json")
	if p == "" {
		return nil, errors.New("rtdb: empty path")
	}
	segs := strings.Split(p, "/")
	for _, s := range segs {
		if err := validSegment(s); err != nil {
			return nil, err
		}
	}
	return segs, nil
}

func (c *Client) logOnce(key, format string, args ...interface{}) {
	c.mu.Lock()
	now := time.Now()
	last, ok := c.lastLog[key]
	if ok && now.Sub(last) < logEvery {
		c.mu.Unlock()
		return
	}
	c.lastLog[key] = now
	c.mu.Unlock()
	c.logf(format, args...)
}

// validateBase only allows https Firebase RTDB hosts so a bearer token can
// never be sent to another host.
func validateBase(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("rtdb: invalid base URL: %w", err)
	}
	if u.Scheme != "https" {
		return "", errors.New("rtdb: base URL must use https")
	}
	if u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
		return "", errors.New("rtdb: base URL must be scheme://host only")
	}
	host := strings.ToLower(u.Hostname())
	if u.Port() != "" && u.Port() != "443" {
		return "", errors.New("rtdb: base URL must not set a custom port")
	}
	if !isRTDBHost(host) {
		return "", fmt.Errorf("rtdb: host %q is not a Firebase RTDB host", host)
	}
	return "https://" + host, nil
}

func isRTDBHost(h string) bool {
	for _, suf := range []string{".firebaseio.com", ".firebasedatabase.app"} {
		if strings.HasSuffix(h, suf) && len(h) > len(suf) && !strings.HasPrefix(h, ".") {
			return true
		}
	}
	return false
}

// validSegment enforces RTDB key rules: non-empty, ≤768 bytes, none of . $ # [ ] /
// and no ASCII control characters.
func validSegment(s string) error {
	if s == "" {
		return errors.New("rtdb: empty path segment")
	}
	if len(s) > maxSegmentLen {
		return errors.New("rtdb: path segment too long")
	}
	if strings.ContainsAny(s, ".$#[]/") {
		return errors.New("rtdb: path segment contains forbidden character")
	}
	for _, r := range s {
		if r < 0x20 || r == 0x7f {
			return errors.New("rtdb: path segment contains control character")
		}
	}
	return nil
}

// safeRedirect follows at most 3 redirects and only to https RTDB hosts, so a
// bearer token is never forwarded elsewhere (net/http also drops
// Authorization on cross-domain redirects). A redirect that would change the
// method (301/302/303 turn PATCH/PUT into a body-less GET) is not followed:
// the 3xx is returned and Do reports it as an error instead of a false success.
func safeRedirect(req *http.Request, via []*http.Request) error {
	if len(via) > 0 && req.Method != via[0].Method {
		return http.ErrUseLastResponse
	}
	if len(via) >= 3 {
		return errors.New("rtdb: too many redirects")
	}
	if req.URL.Scheme != "https" || !isRTDBHost(strings.ToLower(req.URL.Hostname())) {
		return errors.New("rtdb: refusing redirect to non-RTDB host")
	}
	return nil
}

// redactURLErr drops the URL (which contains tenant/account path) from *url.Error.
func redactURLErr(err error) error {
	var ue *url.Error
	if errors.As(err, &ue) {
		return ue.Err
	}
	return err
}
