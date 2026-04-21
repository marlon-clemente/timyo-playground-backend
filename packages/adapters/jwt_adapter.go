package sharedadapters

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenExpiry = 24 * time.Hour

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type JWTAdapterPort interface {
	GenerateToken(ctx context.Context, userID, role string) (string, error)
	ValidateToken(ctx context.Context, token string) (*CustomClaims, error)
}

type JWTAdapter struct {
	secretKey []byte
}

func NewJWTAdapter(secretKey string) JWTAdapterPort {
	return &JWTAdapter{
		secretKey: []byte(secretKey),
	}
}

func (a *JWTAdapter) GenerateToken(ctx context.Context, userID, role string) (string, error) {
	claims := CustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secretKey)
}

func (a *JWTAdapter) ValidateToken(ctx context.Context, token string) (*CustomClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return a.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*CustomClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
