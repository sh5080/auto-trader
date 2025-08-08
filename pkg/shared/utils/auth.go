package utils

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// UserClaims 애플리케이션 공통 클레임
type UserClaims struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}

// GenerateTokens 액세스/리프레시 토큰 생성
func GenerateTokens(userID, secret string, accessTTL, refreshTTL time.Duration, refreshJTI string) (accessToken string, refreshToken string, err error) {
	// Access Token
	accessClaims := &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        refreshJTI,
			Issuer:    "auto-trader",
			Subject:   "access-token",
			Audience:  jwt.ClaimStrings{"auto-trader"},
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = at.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	// Refresh Token (with JTI)
	refreshClaims := &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        refreshJTI,
			Issuer:    "auto-trader",
			Subject:   "refresh-token",
			Audience:  jwt.ClaimStrings{"auto-trader"},
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = rt.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// ParseToken 토큰 파싱/검증
func ParseToken(tokenString, secret string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("유효하지 않은 토큰입니다")
	}
	return claims, nil
}

// GetUserID Fiber 컨텍스트에서 유저ID 가져오기
func GetUserID(c *fiber.Ctx) string {
	if v := c.Locals("userID"); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return ""
}
