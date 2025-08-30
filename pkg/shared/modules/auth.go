package modules

import (
	"auto-trader/ent"
	"auto-trader/pkg/domain/auth"
	"auto-trader/pkg/domain/user"
	"auto-trader/pkg/shared/config"
)

// AuthModule 인증 모듈
type AuthModule struct {
	Service    auth.Service
	Controller *auth.Controller
	cfg        *config.Config
}

// NewAuthModule 인증 모듈 초기화
func NewAuthModule(entClient *ent.Client, cfg *config.Config) *AuthModule {
	// User 도메인 초기화
	userRepo := user.NewEntRepository(entClient)
	userService := user.NewService(userRepo)
	
	// Auth 도메인 초기화 (cfg에서 JWT 설정 추출)
	authService := auth.NewService(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL, userService)
	authController := auth.NewController(authService)

	return &AuthModule{
		Service:    authService,
		Controller: authController,
		cfg:        cfg,
	}
}
