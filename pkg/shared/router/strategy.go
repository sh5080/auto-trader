package router

import (
	"auto-trader/pkg/domain/strategy"

	"github.com/gofiber/fiber/v2"
)

// SetupStrategyRoutes 전략 관련 라우트 설정
func SetupStrategyRoutes(v1 fiber.Router, controller *strategy.Controller) {
	strategies := v1.Group("/strategies")

	// 전략 목록 조회
	strategies.Get("/", controller.GetAllStrategies)

	// 전략 상세 조회
	strategies.Get("/:id", controller.GetStrategy)

	// 전략 상태 조회
	strategies.Get("/:id/status", controller.GetStrategyStatus)

	// 전략 성과 조회
	strategies.Get("/:id/performance", controller.GetStrategyPerformance)

	// 전략 시작
	strategies.Post("/:id/start", controller.StartStrategy)

	// 전략 중지
	strategies.Post("/:id/stop", controller.StopStrategy)

	// 전략 재시작
	strategies.Post("/:id/restart", controller.RestartStrategy)

	// 전략 생성
	strategies.Post("/", controller.CreateStrategy)

	// 전략 수정
	strategies.Put("/:id", controller.UpdateStrategy)

	// 전략 삭제
	strategies.Delete("/:id", controller.DeleteStrategy)
}
