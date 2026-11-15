package rtdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"golang.org/x/oauth2"
)

const testBase = "https://withdraw-monitor-new.firebaseio.com"

type captured struct {
	mu     sync.Mutex
	n      int
	method string
	path   string
	auth   string
	ctype  string
	body   []byte
}

func (c *captured) get() captured {
	c.mu.Lock()
	defer c.mu.Unlock()
	return captured{n: c.n, method: c.method, path: c.path, auth: c.auth, ctype: c.ctype, body: c.body}
}

func fakeRTDB(t *testing.T, status int, resp string) (*httptest.Server, *captured) {
	t.Helper()
	cap := &captured{}
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		cap.mu.Lock()
		cap.n++
		cap.method, cap.path = r.Method, r.URL.EscapedPath()
		cap.auth, cap.ctype, cap.body = r.Header.Get("Authorization"), r.Header.Get("Content-Type"), b
		cap.mu.Unlock()
		w.WriteHeader(status)
		_, _ = io.WriteString(w, resp)
	}))
	t.Cleanup(srv.Close)
	return srv, cap
}

type logBuf struct {
	mu   sync.Mutex
	msgs []string
}

func (l *logBuf) logf(f string, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.msgs = append(l.msgs, fmt.Sprintf(f, a...))
}

func (l *logBuf) all() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return strings.Join(l.msgs, "\n")
}

type errTS struct{}

func (errTS) Token() (*oauth2.Token, error) { return nil, errors.New("metadata unreachable") }

// newTestClient builds a client against the fake server (bypassing the host
// allowlist, which is covered by TestValidateBase).
func newTestClient(t *testing.T, srv *httptest.Server, mode Mode, ts oauth2.TokenSource, lb *logBuf) *Client {
	t.Helper()
	opts := []Option{WithHTTPClient(srv.Client()), WithLogger(lb.logf)}
	if ts != nil {
		opts = append(opts, WithTokenSource(ts))
	}
	c, err := New(context.Background(), testBase, mode, opts...)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	c.base = srv.URL
	return c
}

func TestDo_AttachesBearer(t *testing.T) {
	srv, cap := fakeRTDB(t, 200, `{"ok":true}`)
	lb := &logBuf{}
	c := newTestClient(t, srv, Prefer, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "t"}), lb)

	body := map[string]string{"status": "1", "updated_at": "2026-11-13 10:00:00"}
	st, resp, err := c.Do(context.Background(), http.MethodPatch, []string{"svc", "svc_123"}, body)
	if err != nil || st != 200 || string(resp) != `{"ok":true}` {
		t.Fatalf("Do = %d %q %v", st, resp, err)
	}
	g := cap.get()
	if g.auth != "Bearer t" {
		t.Errorf("Authorization = %q, want Bearer t", g.auth)
	}
	if g.method != http.MethodPatch || g.path != "/svc/svc_123.json" || g.ctype != "application/json" {
		t.Errorf("got %s %s ctype=%q", g.method, g.path, g.ctype)
	}
	var got map[string]string
	if err := json.Unmarshal(g.body, &got); err != nil || got["status"] != "1" || got["updated_at"] != "2026-11-13 10:00:00" || len(got) != 2 {
		t.Errorf("body = %s (%v)", g.body, err)
	}
	if strings.Contains(lb.all(), "Bearer") || strings.Contains(lb.all(), "updated_at") {
		t.Errorf("log leaks token/body: %s", lb.all())
	}
}

func TestDo_OffMode_NoHeader(t *testing.T) {
	srv, cap := fakeRTDB(t, 200, `{}`)
	// even if a token source is supplied, Off must not use it
	c := newTestClient(t, srv, Off, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "t"}), &logBuf{})
	if _, _, err := c.Do(context.Background(), http.MethodPatch, []string{"svc", "svc_1"}, map[string]string{"a": "b"}); err != nil {
		t.Fatal(err)
	}
	if g := cap.get(); g.auth != "" || g.n != 1 {
		t.Errorf("Off mode: auth=%q n=%d", g.auth, g.n)
	}
}

func TestDo_PreferNoToken_SendsLegacy(t *testing.T) {
	srv, cap := fakeRTDB(t, 200, `{}`)
	lb := &logBuf{}
	c := newTestClient(t, srv, Prefer, errTS{}, lb)
	for i := 0; i < 3; i++ {
		if _, _, err := c.Do(context.Background(), http.MethodPatch, []string{"svc", "svc_1"}, map[string]string{"a": "b"}); err != nil {
			t.Fatal(err)
		}
	}
	if g := cap.get(); g.auth != "" || g.n != 3 {
		t.Errorf("prefer/no token: auth=%q n=%d", g.auth, g.n)
	}
	if n := strings.Count(lb.all(), "rtdb_token_unavailable"); n != 1 {
		t.Errorf("rtdb_token_unavailable logged %d times, want 1 (rate limited)", n)
	}
}

func TestDo_RequireNoToken_Blocks(t *testing.T) {
	srv, cap := fakeRTDB(t, 200, `{}`)
	c := newTestClient(t, srv, Require, errTS{}, &logBuf{})
	_, _, err := c.Do(context.Background(), http.MethodPatch, []string{"svc", "svc_1"}, map[string]string{"a": "b"})
	if !errors.Is(err, ErrNoToken) {
		t.Fatalf("err = %v, want ErrNoToken", err)
	}
	if g := cap.get(); g.n != 0 {
		t.Errorf("request reached server %d times", g.n)
	}

	// Require with no token source at all
	c2 := newTestClient(t, srv, Require, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "t"}), &logBuf{})
	c2.ts = nil
	if _, _, err := c2.Do(context.Background(), http.MethodGet, []string{"svc"}, nil); !errors.Is(err, ErrNoToken) {
		t.Fatalf("nil ts: err = %v", err)
	}
}

func TestDo_401_LogsDenied(t *testing.T) {
	srv, _ := fakeRTDB(t, 401, `{"error":"Permission denied"}`)
	lb := &logBuf{}
	c := newTestClient(t, srv, Prefer, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "secret-token"}), lb)
	st, resp, err := c.Do(context.Background(), http.MethodPatch, []string{"svc", "svc_1"}, map[string]string{"acc_bank_number": "1234567890"})
	if err == nil || st != 401 {
		t.Fatalf("want error + 401, got %d %v", st, err)
	}
	if !strings.Contains(string(resp), "Permission denied") {
		t.Errorf("resp = %q", resp)
	}
	l := lb.all()
	if !strings.Contains(l, "rtdb_permission_denied status=401") {
		t.Errorf("missing denied log: %s", l)
	}
	for _, leak := range []string{"secret-token", "1234567890", "svc_1", "Permission denied"} {
		if strings.Contains(l, leak) || strings.Contains(err.Error(), leak) {
			t.Errorf("log/err leaks %q: %s | %v", leak, l, err)
		}
	}
}

func TestDo_GET_NoBody(t *testing.T) {
	srv, cap := fakeRTDB(t, 200, `{"status":"1"}`)
	c := newTestClient(t, srv, Prefer, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "t"}), &logBuf{})
	_, resp, err := c.Do(context.Background(), http.MethodGet, []string{"svc", "svc_1"}, nil)
	if err != nil || string(resp) != `{"status":"1"}` {
		t.Fatalf("%q %v", resp, err)
	}
	if g := cap.get(); g.method != http.MethodGet || len(g.body) != 0 {
		t.Errorf("got %s body=%q", g.method, g.body)
	}
}

func TestJSONEscaping(t *testing.T) {
	srv, cap := fakeRTDB(t, 200, `{}`)
	c := newTestClient(t, srv, Off, nil, &logBuf{})
	name := `ธนาคาร "A", \ x"}` + "\n"
	if _, _, err := c.Do(context.Background(), http.MethodPatch, []string{"svc", "svc_1"}, map[string]string{"bank_name": name}); err != nil {
		t.Fatal(err)
	}
	var got map[string]string
	if err := json.Unmarshal(cap.get().body, &got); err != nil {
		t.Fatalf("body not valid JSON: %v", err)
	}
	if got["bank_name"] != name || len(got) != 1 {
		t.Errorf("got %#v", got)
	}
}

func TestValidSegment(t *testing.T) {
	for _, s := range []string{"", ".", "..", "a.b", "a$", "a#", "a[", "a]", "a/b", "a\x00", "a\n", strings.Repeat("x", 769)} {
		if validSegment(s) == nil {
			t.Errorf("validSegment(%q) accepted", s)
		}
	}
	for _, s := range []string{"svc", "svc_123", "SERVICE-01", "ทดสอบ", "a b", strings.Repeat("x", 768)} {
		if err := validSegment(s); err != nil {
			t.Errorf("validSegment(%q) = %v", s, err)
		}
	}
	srv, cap := fakeRTDB(t, 200, `{}`)
	c := newTestClient(t, srv, Off, nil, &logBuf{})
	if _, _, err := c.Do(context.Background(), http.MethodPatch, []string{"svc", "../x"}, nil); err == nil {
		t.Error("Do accepted traversal segment")
	}
	if _, _, err := c.Do(context.Background(), http.MethodPatch, nil, nil); err == nil {
		t.Error("Do accepted empty path")
	}
	if cap.get().n != 0 {
		t.Error("invalid path reached server")
	}
}

func TestSplitPath(t *testing.T) {
	segs, err := SplitPath("svc/svc_123.json")
	if err != nil || strings.Join(segs, "|") != "svc|svc_123" {
		t.Fatalf("%v %v", segs, err)
	}
	segs, err = SplitPath("/svc/notification/MEMBER_APPROVE_BANK.json")
	if err != nil || len(segs) != 3 {
		t.Fatalf("%v %v", segs, err)
	}
	for _, p := range []string{"", ".json", "svc//x.json", "svc/a.b.json", "svc/x.json?auth=abc", "svc/../x.json"} {
		if _, err := SplitPath(p); err == nil {
			t.Errorf("SplitPath(%q) accepted", p)
		}
	}
}

func TestValidateBase(t *testing.T) {
	ok := map[string]string{
		"https://withdraw-monitor-new.firebaseio.com":                             "https://withdraw-monitor-new.firebaseio.com",
		"https://withdraw-monitor-new.firebaseio.com/":                            "https://withdraw-monitor-new.firebaseio.com",
		" https://X.firebaseio.com ":                                              "https://x.firebaseio.com",
		"https://go-care-679ce-default-rtdb.asia-southeast1.firebasedatabase.app": "https://go-care-679ce-default-rtdb.asia-southeast1.firebasedatabase.app",
		"https://db.firebaseio.com:443":                                           "https://db.firebaseio.com",
	}
	for in, want := range ok {
		got, err := validateBase(in)
		if err != nil || got != want {
			t.Errorf("validateBase(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, in := range []string{
		"http://withdraw-monitor-new.firebaseio.com",
		"https://evil.com",
		"https://firebaseio.com",
		"https://.firebaseio.com",
		"https://withdraw-monitor-new.firebaseio.com.evil.com",
		"https://evil.com/x.firebaseio.com",
		"https://user:pw@db.firebaseio.com",
		"https://db.firebaseio.com/svc",
		"https://db.firebaseio.com?auth=x",
		"https://db.firebaseio.com:8443",
		"ftp://db.firebaseio.com",
		"",
	} {
		if _, err := validateBase(in); err == nil {
			t.Errorf("validateBase(%q) accepted", in)
		}
	}
}

func TestRedirectPolicy(t *testing.T) {
	mk := func(s string) *http.Request { r, _ := http.NewRequest("GET", s, nil); return r }
	if err := safeRedirect(mk("https://db.firebaseio.com/a.json"), nil); err != nil {
		t.Errorf("RTDB redirect refused: %v", err)
	}
	if err := safeRedirect(mk("https://evil.com/a.json"), nil); err == nil {
		t.Error("redirect to evil.com allowed")
	}
	if err := safeRedirect(mk("http://db.firebaseio.com/a.json"), nil); err == nil {
		t.Error("redirect to http allowed")
	}
	via := []*http.Request{mk("https://a.firebaseio.com"), mk("https://a.firebaseio.com"), mk("https://a.firebaseio.com")}
	if err := safeRedirect(mk("https://db.firebaseio.com/a.json"), via); err == nil {
		t.Error("4th redirect allowed")
	}
	patch, _ := http.NewRequest(http.MethodPatch, "https://db.firebaseio.com/a.json", nil)
	if err := safeRedirect(mk("https://db.firebaseio.com/b.json"), []*http.Request{patch}); err != http.ErrUseLastResponse {
		t.Errorf("method-changing redirect: err = %v, want ErrUseLastResponse", err)
	}
}

// 301 would turn PATCH into a body-less GET: must surface as an error, not a
// false success. 307 keeps the method and is followed.
func TestDo_RedirectKeepsMethod(t *testing.T) {
	for _, tc := range []struct {
		code    int
		wantErr bool
	}{{http.StatusMovedPermanently, true}, {http.StatusTemporaryRedirect, false}} {
		srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/final.json" {
				if r.Method != http.MethodPatch {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				_, _ = io.WriteString(w, "{}")
				return
			}
			http.Redirect(w, r, "/final.json", tc.code)
		}))
		// Real base host + production redirect policy; the transport dials the test server.
		c, err := New(context.Background(), testBase, Off, WithTransport(rewriteTransport{target: srv}))
		if err != nil {
			t.Fatal(err)
		}
		if c.hc.CheckRedirect == nil || c.hc.Timeout != defaultTimeout {
			t.Fatal("WithTransport dropped redirect policy or timeout")
		}
		_, _, err = c.Do(context.Background(), http.MethodPatch, []string{"start"}, map[string]string{"a": "b"})
		if (err != nil) != tc.wantErr {
			t.Errorf("%d: err = %v, wantErr %v", tc.code, err, tc.wantErr)
		}
		srv.Close()
	}
}

func TestParseModeAndFromEnv(t *testing.T) {
	cases := map[string]Mode{"": Off, "off": Off, "OFF": Off, " prefer ": Prefer, "Require": Require}
	for in, want := range cases {
		if got, err := ParseMode(in); err != nil || got != want {
			t.Errorf("ParseMode(%q) = %v, %v", in, got, err)
		}
	}
	if _, err := ParseMode("on"); err == nil {
		t.Error("ParseMode(on) accepted")
	}

	t.Setenv(EnvAuthMode, "")
	t.Setenv("TEST_RTDB_URL", "")
	c, err := FromEnv(context.Background(), "TEST_RTDB_URL", testBase)
	if err != nil || c.Mode() != Off || c.base != testBase || c.ts != nil {
		t.Fatalf("default FromEnv = %+v, %v", c, err)
	}
	t.Setenv("TEST_RTDB_URL", "https://other-db.firebaseio.com")
	if c, err = FromEnv(context.Background(), "TEST_RTDB_URL", testBase); err != nil || c.base != "https://other-db.firebaseio.com" {
		t.Fatalf("env override = %v, %v", c, err)
	}
	t.Setenv("TEST_RTDB_URL", "https://evil.com")
	if _, err = FromEnv(context.Background(), "TEST_RTDB_URL", testBase); err == nil {
		t.Fatal("FromEnv accepted non-RTDB host")
	}
	t.Setenv("TEST_RTDB_URL", "")
	t.Setenv(EnvAuthMode, "bogus")
	if _, err = FromEnv(context.Background(), "TEST_RTDB_URL", testBase); err == nil {
		t.Fatal("FromEnv accepted bogus mode")
	}
}

func TestNew_NoADC(t *testing.T) {
	// Point ADC at a missing file so FindDefaultCredentials fails fast and
	// deterministically (no metadata server probing).
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", t.TempDir()+"/missing.json")
	lb := &logBuf{}
	c, err := New(context.Background(), testBase, Prefer, WithLogger(lb.logf))
	if err != nil || c.ts != nil || !strings.Contains(lb.all(), "credentials_unavailable") {
		t.Fatalf("prefer without ADC: c=%v err=%v log=%s", c, err, lb.all())
	}
	if _, err := New(context.Background(), testBase, Require, WithLogger(lb.logf)); err == nil {
		t.Fatal("require without ADC: want error")
	}
}

// rewriteTransport sends every request to target while keeping the original
// URL (host) visible to the client and its redirect policy.
type rewriteTransport struct{ target *httptest.Server }

func (rt rewriteTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r2 := r.Clone(r.Context())
	u, _ := url.Parse(rt.target.URL)
	r2.URL.Scheme, r2.URL.Host = u.Scheme, u.Host
	return rt.target.Client().Transport.RoundTrip(r2)
}
