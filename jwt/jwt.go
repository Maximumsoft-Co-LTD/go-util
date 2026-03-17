// Package jwt provides JWT token generation and validation using HS256.
package jwt

import (
	"errors"
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
)

// Service defines the interface for JWT operations.
type Service interface {
	GenerateToken(username string, sessionID string, ttl time.Duration) (string, error)
	ValidateToken(token string) (*gojwt.Token, error)
}

// CustomClaims holds the JWT payload fields.
type CustomClaims struct {
	Username  string `json:"username"`
	SessionID string `json:"session_id"`
	gojwt.RegisteredClaims
}

type service struct {
	secretKey string
	issuer    string
}

// NewService creates a new JWT Service with the given secret and issuer.
// Returns an error if secret is empty.
func NewService(secret string, issuer string) (Service, error) {
	if secret == "" {
		return nil, errors.New("jwt secret must not be empty")
	}
	if issuer == "" {
		issuer = "go-util"
	}
	return &service{
		secretKey: secret,
		issuer:    issuer,
	}, nil
}

func (s *service) GenerateToken(username string, sessionID string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := &CustomClaims{
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

func (s *service) ValidateToken(encodedToken string) (*gojwt.Token, error) {
	return gojwt.Parse(encodedToken, func(token *gojwt.Token) (interface{}, error) {
		if _, valid := token.Method.(*gojwt.SigningMethodHMAC); !valid {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
}
