package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateCreatesSevenDayHS256Token(t *testing.T) {
	secret := "0123456789abcdef0123456789abcdef"
	manager, err := NewTokenManager(secret)
	if err != nil {
		t.Fatalf("NewTokenManager() error = %v", err)
	}

	tokenString, err := manager.Generate("user@example.com")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			t.Fatalf("signing method = %v, want HS256", token.Method.Alg())
		}
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("ParseWithClaims() error = %v", err)
	}
	if !token.Valid {
		t.Fatal("generated token is invalid")
	}
	if claims.Subject != "user@example.com" {
		t.Errorf("subject = %q, want user@example.com", claims.Subject)
	}
	if claims.Issuer != "PayFlow API" {
		t.Errorf("issuer = %q, want PayFlow API", claims.Issuer)
	}
	if claims.ExpiresAt == nil {
		t.Fatal("token has no expiration")
	}
	wantExpiration := time.Now().Add(tokenLifetime)
	if difference := claims.ExpiresAt.Time.Sub(wantExpiration); difference < -time.Second || difference > time.Second {
		t.Errorf("expiration differs from seven-day lifetime by %v", difference)
	}
}

func TestNewTokenManagerRejectsShortSecret(t *testing.T) {
	if _, err := NewTokenManager("too-short"); err == nil {
		t.Fatal("NewTokenManager() succeeded with a short secret")
	}
}
