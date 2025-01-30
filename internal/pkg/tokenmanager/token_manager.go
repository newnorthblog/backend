package tokenmanager

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenManager provides logic for JWT & Refresh tokens generation and parsing.
var ErrAccessTokenExpired = errors.New("token has invalid claims: token is expired")

type TokenManager interface {
	NewJWT(userID *uuid.UUID) (string, time.Duration, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() (uuid.UUID, time.Duration, error)
	ValidateRefreshToken(refreshToken string) (*uuid.UUID, error)
}

type Manager struct {
	signingKey      string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewManager(signingKey string, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	if accessTokenTTL == 0 {
		return nil, errors.New("empty access token ttl")
	}

	if refreshTokenTTL == 0 {
		return nil, errors.New("empty refresh token ttl")
	}

	return &Manager{
		signingKey:      signingKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}, nil
}

func (m *Manager) NewJWT(userID *uuid.UUID) (string, time.Duration, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTokenTTL)),
		Subject:   userID.String(),
	})

	accessToken, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", 0, fmt.Errorf("sign jwt failed")
	}

	return accessToken, m.accessTokenTTL, nil
}

func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil
}

func (m *Manager) NewRefreshToken() (uuid.UUID, time.Duration, error) {
	refreshToken, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("new refresh token failed: %v", err)
	}
	return refreshToken, m.refreshTokenTTL, nil
}

func (m *Manager) ValidateRefreshToken(refreshToken string) (*uuid.UUID, error) {
	id, err := uuid.Parse(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("refresh token uuid parse: %v", err)
	}

	return &id, nil
}
