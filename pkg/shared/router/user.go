package router

import (
	"auto-trader/pkg/domain/user"
	"auto-trader/pkg/shared/config"
	"auto-trader/pkg/shared/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes 유저 관련 라우트 설정
func SetupUserRoutes(v1 fiber.Router, controller *user.Controller, cfg *config.Config) {
	users := v1.Group("/users")
	// 회원가입 (공개)
	users.Post("/", controller.CreateUser)
	protected := users.Group("/", middleware.AuthMiddleware(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL))
	// 유저 조회 (보호)
	protected.Get("/:id", controller.GetUser)
}
