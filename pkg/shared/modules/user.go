package modules

import (
	"auto-trader/ent"
	"auto-trader/pkg/domain/user"
	"auto-trader/pkg/shared/config"
)

// UserModule 사용자 모듈
type UserModule struct {
	Repository user.Repository
	Service    user.Service
	Controller *user.Controller
	cfg        *config.Config
}

// NewUserModule 사용자 모듈 초기화
func NewUserModule(entClient *ent.Client, cfg *config.Config) *UserModule {
	// Repository -> Service -> Controller 순서로 초기화
	repo := user.NewEntRepository(entClient)
	service := user.NewService(repo)
	controller := user.NewController(service)

	return &UserModule{
		Repository: repo,
		Service:    service,
		Controller: controller,
		cfg:        cfg,
	}
}
