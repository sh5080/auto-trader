package router

import (
	"auto-trader/pkg/domain/strategy"
	"auto-trader/pkg/shared/config"
	"auto-trader/pkg/shared/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupStrategyRoutes 전략 관련 라우트 설정
func SetupStrategyRoutes(v1 fiber.Router, controller *strategy.Controller, cfg *config.Config) {
	strategies := v1.Group("/strategies")
	protected := strategies.Group("/", middleware.AuthMiddleware(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL))

	// 전략 목록 조회
	protected.Get("/", controller.GetAllStrategies)

	// 전략 상세 조회
	protected.Get("/:id", controller.GetStrategy)

	// 전략 상태 조회
	protected.Get("/:id/status", controller.GetStrategyStatus)

	// 전략 성과 조회
	protected.Get("/:id/performance", controller.GetStrategyPerformance)

	// 전략 시작
	protected.Post("/:id/start", controller.StartStrategy)

	// 전략 중지
	protected.Post("/:id/stop", controller.StopStrategy)

	// 전략 재시작
	protected.Post("/:id/restart", controller.RestartStrategy)

	// 전략 생성
	protected.Post("/", controller.CreateStrategy)

	// 전략 수정
	protected.Put("/:id", controller.UpdateStrategy)

	// 전략 삭제
	protected.Delete("/:id", controller.DeleteStrategy)
}
