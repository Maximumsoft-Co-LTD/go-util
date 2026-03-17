package goutil

import (
	"errors"
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

// JWTService defines the interface for JWT operations.
type JWTService interface {
	GenerateToken(username string, sessionID string, ttl time.Duration) (string, error)
	ValidateToken(token string) (*gojwt.Token, error)
}

// JWTClaims holds the JWT payload fields.
type JWTClaims struct {
	Username  string `json:"username"`
	SessionID string `json:"session_id"`
	gojwt.RegisteredClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

// NewJWTService creates a new JWTService with the given secret and issuer.
// Returns an error if secret is empty.
func NewJWTService(secret string, issuer string) (JWTService, error) {
	if secret == "" {
		return nil, errors.New("jwt secret must not be empty")
	}
	if issuer == "" {
		issuer = "go-util"
	}
	return &jwtService{
		secretKey: secret,
		issuer:    issuer,
	}, nil
}

func (s *jwtService) GenerateToken(username string, sessionID string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := &JWTClaims{
		Username:  username,
		SessionID: sessionID,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(now.Add(ttl)),
			Issuer:    s.issuer,
			IssuedAt:  gojwt.NewNumericDate(now),
		},
	}
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}
	return signed, nil
}

func (s *jwtService) ValidateToken(encodedToken string) (*gojwt.Token, error) {
	return gojwt.Parse(encodedToken, func(token *gojwt.Token) (interface{}, error) {
		if _, valid := token.Method.(*gojwt.SigningMethodHMAC); !valid {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
}
