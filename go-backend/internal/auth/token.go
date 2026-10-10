package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenLifetime = 7 * 24 * time.Hour

// TokenManager signs authentication tokens with an application secret.
type TokenManager struct {
	secret []byte
}

// NewTokenManager creates a token manager from a secret of at least 32 bytes.
func NewTokenManager(secret string) (*TokenManager, error) {
	if len(secret) < 32 {
		return nil, errors.New("JWT_SECRET must contain at least 32 bytes")
	}
	return &TokenManager{secret: []byte(secret)}, nil
}

// Generate creates a signed JWT identifying the user by username.
func (m *TokenManager) Generate(username string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    "PayFlow API",
		Subject:   username,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenLifetime)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign JWT: %w", err)
	}
	return signedToken, nil
}
