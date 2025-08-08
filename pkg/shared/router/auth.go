package router

import (
	"auto-trader/pkg/domain/auth"

	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes 인증 관련 라우트 설정
func SetupAuthRoutes(v1 fiber.Router, controller *auth.Controller) {
	auth := v1.Group("/auth")

	auth.Post("/login", controller.Login)
	auth.Post("/refresh", controller.Refresh)
}
