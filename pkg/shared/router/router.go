package router

import (
	"auto-trader/pkg/domain/auth"
	"auto-trader/pkg/domain/portfolio"
	"auto-trader/pkg/domain/strategy"
	"auto-trader/pkg/domain/user"
	"auto-trader/pkg/shared/config"
	"auto-trader/pkg/shared/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// Router 메인 라우터 구조체
type Router struct {
	app         *fiber.App
	riskManager *middleware.Manager
}

// New 새로운 라우터 인스턴스 생성
func New(riskManager *middleware.Manager) *Router {
	var cfg fiber.Config
	cfg.AppName = "Auto Trader"
	cfg.ServerHeader = "Auto Trader v1.0"

	app := fiber.New(cfg)

	return &Router{
		app:         app,
		riskManager: riskManager,
	}
}

// SetupRoutes 모든 라우트 설정
func (r *Router) SetupRoutes(
	strategyController *strategy.Controller,
	portfolioController *portfolio.Controller,
	authController *auth.Controller,
	userController *user.Controller,
	cfg *config.Config,
) {
	// 글로벌 미들웨어 설정
	r.setupGlobalMiddleware()

	// Health Check
	r.app.Get("/health", r.healthCheck)

	// API v1 그룹
	v1 := r.app.Group("/api/v1")
	// 각 도메인 라우트 설정
	SetupAuthRoutes(v1, authController)
	SetupStrategyRoutes(v1, strategyController, cfg)
	SetupPortfolioRoutes(v1, portfolioController, cfg)
	SetupUserRoutes(v1, userController, cfg)

	r.app.Use(middleware.SetupNotFoundHandler())
}

// SetupSwagger Swagger UI 설정
func (r *Router) SetupSwagger() {
	r.app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	r.app.Get("/docs/*", fiberSwagger.WrapHandler)
	logrus.Info("📚 Swagger UI 설정 완료: /docs/*")
}

// GetApp Fiber 앱 인스턴스 반환
func (r *Router) GetApp() *fiber.App {
	return r.app
}

// setupGlobalMiddleware 글로벌 미들웨어 설정
func (r *Router) setupGlobalMiddleware() {
	// 패닉 복구 (가장 먼저)
	r.app.Use(middleware.SetupPanicRecovery())

	// CORS 설정
	r.app.Use(middleware.SetupCORS())

	// 로깅 미들웨어
	r.app.Use(middleware.SetupStructuredLogger())

	// 에러 핸들링 미들웨어
	r.app.Use(middleware.SetupAdvancedErrorHandler())
}

// healthCheck 헬스 체크 핸들러
func (r *Router) healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"service": "auto-trader",
		"version": "1.0.0",
		"timestamp": fiber.Map{
			"now": c.Context().Time(),
		},
	})
}
