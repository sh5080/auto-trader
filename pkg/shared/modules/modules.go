package modules

import (
	"auto-trader/ent"
	"auto-trader/pkg/shared/config"
	"auto-trader/pkg/shared/middleware"

	"github.com/sirupsen/logrus"
)

type Modules struct {
	User      *UserModule
	Auth      *AuthModule
	Strategy  *StrategyModule
	Portfolio *PortfolioModule
}

// InitializeModules 모듈 초기화
func InitializeModules(entClient *ent.Client, riskManager *middleware.Manager, cfg *config.Config) *Modules {
	logrus.Info("🔧 모듈별 의존성 초기화 중...")

	// 1. User 모듈 초기화
	userModule := NewUserModule(entClient, cfg)
	logrus.Info("✅ User 모듈 초기화 완료")

	// 2. Auth 모듈 초기화
	authModule := NewAuthModule(entClient, cfg)
	logrus.Info("✅ Auth 모듈 초기화 완료")

	// 3. Strategy 모듈 초기화
	strategyModule := NewStrategyModule(entClient, riskManager, cfg)
	logrus.Info("✅ Strategy 모듈 초기화 완료")

	// 4. Portfolio 모듈 초기화
	portfolioModule := NewPortfolioModule(entClient, cfg)
	logrus.Info("✅ Portfolio 모듈 초기화 완료")

	return &Modules{
		User:      userModule,
		Auth:      authModule,
		Strategy:  strategyModule,
		Portfolio: portfolioModule,
	}
}
