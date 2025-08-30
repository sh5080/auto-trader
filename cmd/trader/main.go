package main

import (
	"log"

	_ "auto-trader/docs" // swagger docs
	"auto-trader/pkg/shared/config"
	"auto-trader/pkg/shared/database"
	"auto-trader/pkg/shared/middleware"
	"auto-trader/pkg/shared/modules"
	"auto-trader/pkg/shared/router"

	"github.com/sirupsen/logrus"
)

// @title Auto Trader API
// @version 1.0
// @description 자동 주식 거래 시스템 API
// @termsOfService http://swagger.io/terms/

// @host localhost:8087
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// 애플리케이션 초기화
	app := initializeApp()

	// 서버 시작
	startServer(app)
}

func initializeApp() *router.Router {
	logrus.Info("🚀 Auto Trader 초기화 시작")

	// 설정 로드
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("❌ 설정 로드 실패: %v", err)
	}

	// 의존성 초기화
	dependencies := initializeDependencies(cfg)

	// 라우터 생성 및 설정
	mainRouter := router.New(dependencies.RiskManager)

	// Swagger 라우트 추가
	mainRouter.SetupSwagger()
	// 라우트 설정
	mainRouter.SetupRoutes(
		dependencies.Modules.Strategy.Controller,
		dependencies.Modules.Portfolio.Controller,
		dependencies.Modules.Auth.Controller,
		dependencies.Modules.User.Controller,
		cfg,
	)

	// 백그라운드 작업 시작
	startBackgroundTasks(dependencies)

	logrus.Info("✅ Auto Trader 초기화 완료")
	return mainRouter
}

// Dependencies 애플리케이션 의존성들
type Dependencies struct {
	Database    database.DB
	RiskManager *middleware.Manager
	Modules     *modules.Modules
}

func initializeDependencies(cfg *config.Config) *Dependencies {
	logrus.Info("🔧 의존성 초기화 중...")

	// 1. 데이터베이스 연결 초기화
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		logrus.Fatalf("❌ 데이터베이스 연결 실패: %v", err)
	}
	logrus.Info("✅ 데이터베이스 연결 성공")

	// 2. 리스크 관리자 초기화
	riskManager := middleware.NewManager(cfg)
	logrus.Info("✅ 리스크 관리자 초기화")

	// 3. ent 클라이언트 가져오기
	entClient := db.GetEntClient()
	logrus.Info("✅ ent 클라이언트 준비 완료")

	// 4. 모듈 초기화
	allModules := modules.InitializeModules(entClient, riskManager, cfg)

	return &Dependencies{
		Database:    db,
		RiskManager: riskManager,
		Modules:     allModules,
	}
}

func startBackgroundTasks(deps *Dependencies) {
	logrus.Info("🔄 백그라운드 서비스 시작 중...")

	// 전략 서비스 시작 (비동기)
	go func() {
		if err := deps.Modules.Strategy.Service.Start(); err != nil {
			logrus.Errorf("❌ 전략 서비스 시작 실패: %v", err)
		} else {
			logrus.Info("✅ 전략 서비스 시작 완료")
		}
	}()

	logrus.Info("🎯 백그라운드 서비스 시작 완료")
}

func startServer(mainRouter *router.Router) {
	// 설정에서 포트 가져오기
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("❌ 설정 로드 실패: %v", err)
	}
	port := ":" + cfg.Server.Port

	logrus.Info("🌟 ================================")
	logrus.Info("🚀 Auto Trader 서버 시작")
	logrus.Infof("📡 포트: %s", port)
	logrus.Infof("🌐 서버: http://localhost%s", port)
	logrus.Infof("❤️  헬스체크: http://localhost%s/health", port)
	logrus.Infof("📊 API: http://localhost%s/api/v1", port)
	logrus.Infof("🎯 전략: http://localhost%s/api/v1/strategies", port)
	logrus.Infof("💼 포트폴리오: http://localhost%s/api/v1/portfolio", port)
	logrus.Infof("📚 Swagger: http://localhost%s/docs/", port)
	logrus.Infof("📖 Docs: http://localhost%s/docs", port)
	logrus.Info("🌟 ================================")

	if err := mainRouter.GetApp().Listen(port); err != nil {
		log.Fatalf("❌ 서버 시작 실패: %v", err)
	}
}
