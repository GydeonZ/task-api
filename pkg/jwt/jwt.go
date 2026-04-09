package jwtutil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenType distinguishes between access and refresh tokens.
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims are the custom JWT claims used by this application.
type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

// Manager provides token generation and parsing.
type Manager struct {
	secret []byte
}

// New creates a new Manager.
func New(secret string) *Manager {
	return &Manager{secret: []byte(secret)}
}

// Generate creates a signed JWT for the given user.
func (m *Manager) Generate(userID uuid.UUID, tokenType TokenType, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// Parse validates and parses a JWT, returning the embedded Claims.
func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
