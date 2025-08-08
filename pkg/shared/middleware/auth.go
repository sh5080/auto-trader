package middleware

import (
	"errors"
	"time"

	authstore "auto-trader/pkg/shared/auth"
	"auto-trader/pkg/shared/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AuthMiddleware 액세스 토큰 검증 및 만료 시 리프레시 토큰으로 자동 재발급(RTR)
// 클라이언트는 AT는 Authorization 헤더, RT는 X-Refresh-Token 헤더로 전송
func AuthMiddleware(secret string, accessTTL, refreshTTL time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return utils.UnauthorizedResponse(c, "Authorization 헤더가 필요합니다")
		}
		token := auth
		if len(auth) > 7 && (auth[:7] == "Bearer " || auth[:7] == "bearer ") {
			token = auth[7:]
		}

		claims, err := utils.ParseToken(token, secret)
		if err == nil {
			c.Locals("userID", claims.UserID)
			return c.Next()
		}

		// 만료라면 자동 리프레시 시도
		if errors.Is(err, jwt.ErrTokenExpired) {
			refreshHeader := c.Get("X-Refresh-Token")
			if refreshHeader == "" {
				return utils.UnauthorizedResponse(c, "리프레시 토큰이 필요합니다")
			}
			// refresh 토큰은 Bearer 포맷일 수도 있음
			refreshToken := refreshHeader
			if len(refreshHeader) > 7 && (refreshHeader[:7] == "Bearer " || refreshHeader[:7] == "bearer ") {
				refreshToken = refreshHeader[7:]
			}

			rClaims, rErr := utils.ParseToken(refreshToken, secret)
			if rErr != nil {
				return utils.UnauthorizedResponse(c, "리프레시 토큰이 유효하지 않습니다")
			}
			userID := rClaims.UserID
			jti := rClaims.ID
			current, ok := authstore.GetRefreshJTI(userID)
			if !ok || current != jti {
				return utils.UnauthorizedResponse(c, "리프레시 토큰이 회전되어 사용 불가합니다")
			}

			// 회전 및 재발급
			newJTI := uuid.NewString()
			access, refresh, genErr := utils.GenerateTokens(userID, secret, accessTTL, refreshTTL, newJTI)
			if genErr != nil {
				return utils.InternalServerErrorResponse(c, "토큰 재발급 실패", genErr)
			}
			authstore.SetRefreshJTI(userID, newJTI)

			// 새 토큰을 응답 헤더에 제공 (클라이언트가 저장하도록 위임)
			c.Set("X-New-Access-Token", access)
			c.Set("X-New-Refresh-Token", refresh)
			c.Locals("userID", userID)
			return c.Next()
		}

		// 그 외 에러는 인증 실패
		return utils.UnauthorizedResponse(c, "액세스 토큰 검증 실패")
	}
}
