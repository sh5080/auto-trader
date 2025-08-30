package modules

import (
	"auto-trader/ent"
	"auto-trader/pkg/domain/portfolio"
	"auto-trader/pkg/shared/config"
)

// PortfolioModule 포트폴리오 모듈
type PortfolioModule struct {
	Repository portfolio.Repository
	Service    portfolio.Service
	Controller *portfolio.Controller
	cfg        *config.Config
}

// NewPortfolioModule 포트폴리오 모듈 초기화
func NewPortfolioModule(entClient *ent.Client, cfg *config.Config) *PortfolioModule {
	// Repository -> Service -> Controller 순서로 초기화
	repo := portfolio.NewEntRepository(entClient)
	service := portfolio.NewService(repo)
	controller := portfolio.NewController(service)

	return &PortfolioModule{
		Repository: repo,
		Service:    service,
		Controller: controller,
		cfg:        cfg,
	}
}
