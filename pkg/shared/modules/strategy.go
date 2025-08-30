package modules

import (
	"auto-trader/ent"
	"auto-trader/pkg/domain/strategy"
	"auto-trader/pkg/shared/config"
	"auto-trader/pkg/shared/middleware"
)

// StrategyModule 전략 모듈
type StrategyModule struct {
	Repository strategy.Repository
	Service    strategy.Service
	Controller *strategy.Controller
	cfg        *config.Config
}

// NewStrategyModule 전략 모듈 초기화
func NewStrategyModule(entClient *ent.Client, riskManager *middleware.Manager, cfg *config.Config) *StrategyModule {
	// Repository 초기화
	repo := strategy.NewEntRepository(entClient)

	// 현재는 portfolio/order 도메인이 없으므로 nil
	var dataCollector strategy.Collector = nil
	var executor strategy.Executor = nil

	// Service 초기화
	service := strategy.NewService(
		repo,
		dataCollector, // TODO: portfolio 도메인 완성 후 연결
		executor,      // TODO: order 도메인 완성 후 연결
		riskManager,
		cfg,
	)

	// Controller 초기화
	controller := strategy.NewController(service)

	return &StrategyModule{
		Repository: repo,
		Service:    service,
		Controller: controller,
		cfg:        cfg,
	}
}
