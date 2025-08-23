package auth

import (
	"auto-trader/pkg/domain/auth/dto"
	authstore "auto-trader/pkg/shared/auth"
	"auto-trader/pkg/shared/utils"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service Service
}

func NewController(s Service) *Controller { return &Controller{service: s} }

// Login
// @Summary 로그인
// @Description 이메일과 비밀번호로 로그인합니다
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "로그인 요청"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/login [post]
func (ctl *Controller) Login(c *fiber.Ctx) error {
	loginDto, err := utils.ParseAndValidate[dto.LoginBody](c)
	if err != nil {
		return err
	}

	res, err := ctl.service.Login(loginDto)
	if err != nil {
		return utils.UnauthorizedResponse(c, "이메일 또는 비밀번호가 올바르지 않습니다")
	}
	// RTR 저장
	authstore.SetRefreshJTI(res.UserID, res.RefreshJTI)
	return utils.SuccessResponse(c, fiber.Map{
		"accessToken":  res.AccessToken,
		"refreshToken": res.RefreshToken,
	})
}

// Refresh
// @Summary 토큰 갱신
// @Description Authorization 헤더의 Refresh 토큰으로 Access 토큰을 갱신합니다
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {refreshToken}"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /auth/refresh [post]
func (ctl *Controller) Refresh(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return utils.UnauthorizedResponse(c, "Authorization 헤더가 필요합니다")
	}
	res, err := ctl.service.Refresh(token)
	if err != nil {
		return utils.UnauthorizedResponse(c, err.Error())
	}
	// RTR 저장
	authstore.SetRefreshJTI(res.UserID, res.RefreshJTI)
	return utils.SuccessResponse(c, fiber.Map{
		"accessToken":  res.AccessToken,
		"refreshToken": res.RefreshToken,
	})
}
