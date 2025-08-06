package main

import (
	"log"

	"auto-trader/pkg/domain/strategy"
	"auto-trader/pkg/shared/config"
	"auto-trader/pkg/shared/database"
	"auto-trader/pkg/shared/middleware"
	"auto-trader/pkg/shared/router"

	"github.com/sirupsen/logrus"
)

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

	// 라우트 설정
	mainRouter.SetupRoutes(dependencies.StrategyController)

	// 백그라운드 작업 시작
	startBackgroundTasks(dependencies)

	logrus.Info("✅ Auto Trader 초기화 완료")
	return mainRouter
}

// 정리된 Dependencies - 실제 사용하는 것만
type Dependencies struct {
	Database           database.DB
	RiskManager        *middleware.Manager
	StrategyService    strategy.Service
	StrategyController *strategy.Controller
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

	// 3. 전략 리포지토리 초기화
	strategyRepo := strategy.NewDBRepository(db.GetDB())
	logrus.Info("✅ 전략 리포지토리 초기화)")

	// 4. 전략 서비스 초기화
	// 현재는 data/order 도메인이 없으므로 nil
	var dataCollector strategy.Collector = nil
	var executor strategy.Executor = nil

	strategyService := strategy.NewService(
		strategyRepo,
		dataCollector, // TODO: data 도메인 완성 후 연결
		executor,      // TODO: order 도메인 완성 후 연결
		riskManager,
		cfg,
	)
	logrus.Info("✅ 전략 서비스 초기화")

	// 5. 전략 컨트롤러 초기화
	strategyController := strategy.NewController(strategyService)
	logrus.Info("✅ 전략 컨트롤러 초기화")

	logrus.Info("🎉 모든 의존성 초기화 완료")

	return &Dependencies{
		Database:           db,
		RiskManager:        riskManager,
		StrategyService:    strategyService,
		StrategyController: strategyController,
	}
}

func startBackgroundTasks(deps *Dependencies) {
	logrus.Info("🔄 백그라운드 서비스 시작 중...")

	// 전략 서비스 시작 (비동기)
	go func() {
		if err := deps.StrategyService.Start(); err != nil {
			logrus.Errorf("❌ 전략 서비스 시작 실패: %v", err)
		} else {
			logrus.Info("✅ 전략 서비스 시작 완료")
		}
	}()

	logrus.Info("🎯 백그라운드 서비스 시작 완료")
}

func startServer(mainRouter *router.Router) {
	port := ":8787"

	logrus.Info("🌟 ================================")
	logrus.Info("🚀 Auto Trader 서버 시작")
	logrus.Infof("📡 포트: %s", port)
	logrus.Infof("🌐 서버: http://localhost%s", port)
	logrus.Infof("❤️  헬스체크: http://localhost%s/health", port)
	logrus.Infof("📊 API: http://localhost%s/api/v1", port)
	logrus.Infof("🎯 전략: http://localhost%s/api/v1/strategies", port)
	logrus.Info("🌟 ================================")

	if err := mainRouter.GetApp().Listen(port); err != nil {
		log.Fatalf("❌ 서버 시작 실패: %v", err)
	}
}
