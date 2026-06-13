package crypto

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AccessTokenClaims are embedded in the short-lived access JWT.
type AccessTokenClaims struct {
	UserID      uuid.UUID `json:"uid"`
	TenantID    uuid.UUID `json:"tid"`
	TenantSlug  string    `json:"tslug"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	IsSuperuser bool      `json:"su"`
	Roles       []string  `json:"roles"`
	Permissions []string  `json:"perms"` // "resource:action" strings
	jwt.RegisteredClaims
}

// RefreshTokenClaims are embedded in the long-lived refresh JWT.
type RefreshTokenClaims struct {
	UserID uuid.UUID `json:"uid"`
	jwt.RegisteredClaims
}

// JWTManager handles token generation and parsing.
type JWTManager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJWTManager(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *JWTManager {
	return &JWTManager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

// GenerateAccessToken creates a signed JWT access token.
func (m *JWTManager) GenerateAccessToken(claims AccessTokenClaims) (string, error) {
	now := time.Now()
	claims.RegisteredClaims = jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
		ID:        uuid.NewString(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.accessSecret)
}

// GenerateRefreshToken creates a signed JWT refresh token.
func (m *JWTManager) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	now := time.Now()
	claims := RefreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTTL)),
			ID:        uuid.NewString(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.refreshSecret)
}

// ParseAccessToken validates and parses an access token string.
func (m *JWTManager) ParseAccessToken(tokenStr string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(t *jwt.Token) (any, error) {
		// Pin parsing to the same algorithm used for generation. Accepting every
		// HMAC variant would allow tokens signed under a different algorithm.
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
		}
		return m.accessSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid access token")
	}
	return claims, nil
}

// ParseRefreshToken validates and parses a refresh token string.
func (m *JWTManager) ParseRefreshToken(tokenStr string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshTokenClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
		}
		return m.refreshSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}
	return claims, nil
}

// RefreshTTL exposes the configured refresh token lifetime.
func (m *JWTManager) RefreshTTL() time.Duration {
	return m.refreshTTL
}
