package router

import (
	"auto-trader/pkg/domain/portfolio"
	"auto-trader/pkg/shared/config"
	"auto-trader/pkg/shared/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupPortfolioRoutes 포트폴리오 관련 라우트 설정
func SetupPortfolioRoutes(v1 fiber.Router, controller *portfolio.Controller, cfg *config.Config) {
	portfolio := v1.Group("/portfolio")
	protected := portfolio.Group("/", middleware.AuthMiddleware(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL))

	// 포트폴리오 조회
	protected.Get("/", controller.GetPortfolio)

	// 포트폴리오 요약 조회
	protected.Get("/summary", controller.GetPortfolioSummary)

	// 포트폴리오 새로고침
	protected.Post("/refresh", controller.RefreshPortfolio)

	// 보유 주식 관련
	positions := protected.Group("/positions")
	positions.Get("/", controller.GetPositions)       // 보유 주식 목록 조회
	positions.Get("/:symbol", controller.GetPosition) // 특정 보유 주식 조회

	// 가격 관련
	prices := protected.Group("/prices")
	prices.Get("/", controller.GetCurrentPrices)       // 여러 종목 현재가 조회
	prices.Get("/:symbol", controller.GetCurrentPrice) // 특정 종목 현재가 조회

	// 거래 내역
	trades := protected.Group("/trades")
	trades.Get("/", controller.GetTradeHistory) // 거래 내역 조회
}
