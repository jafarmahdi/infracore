package crypto

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestJWTManagerRoundTrip(t *testing.T) {
	manager := NewJWTManager("access-secret", "refresh-secret", time.Minute, time.Hour)
	userID := uuid.New()
	tenantID := uuid.New()

	accessToken, err := manager.GenerateAccessToken(AccessTokenClaims{
		UserID:   userID,
		TenantID: tenantID,
		Email:    "admin@infracore.io",
		Roles:    []string{"admin"},
	})
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	accessClaims, err := manager.ParseAccessToken(accessToken)
	if err != nil {
		t.Fatalf("ParseAccessToken() error = %v", err)
	}
	if accessClaims.UserID != userID || accessClaims.TenantID != tenantID {
		t.Fatalf("unexpected access claims: %+v", accessClaims)
	}

	refreshToken, err := manager.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}
	refreshClaims, err := manager.ParseRefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("ParseRefreshToken() error = %v", err)
	}
	if refreshClaims.UserID != userID {
		t.Fatalf("refresh UserID = %s, want %s", refreshClaims.UserID, userID)
	}
}

func TestJWTManagerRejectsDifferentHMACAlgorithm(t *testing.T) {
	manager := NewJWTManager("access-secret", "refresh-secret", time.Minute, time.Hour)
	claims := AccessTokenClaims{
		UserID: uuid.New(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signed, err := token.SignedString([]byte("access-secret"))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}

	if _, err := manager.ParseAccessToken(signed); err == nil {
		t.Fatal("ParseAccessToken() accepted an HS512 token")
	}
}
