package user

import (
	"auto-trader/pkg/domain/user/dto"
	"auto-trader/pkg/shared/types"
	"auto-trader/pkg/shared/utils"

	"github.com/gofiber/fiber/v2"
)

type Controller struct{ service Service }

func NewController(s Service) *Controller { return &Controller{service: s} }

// CreateUser
// @Summary 사용자 회원가입
// @Description 새 사용자를 생성합니다
// @Tags users
// @Accept json
// @Produce json
// @Param body body dto.CreateUserRequest true "회원가입 요청"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users [post]
func (ctl *Controller) CreateUser(c *fiber.Ctx) error {
	var dto dto.CreateUserRequest

	if err := c.BodyParser(&dto); err != nil {
		return utils.ValidationErrorResponse(c, "잘못된 요청 본문입니다")
	}
	if err := utils.ValidateStruct(dto); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	u, err := ctl.service.CreateUser(dto)
	if err != nil {
		return utils.CommonErrorResponse(c, err, "사용자 생성 실패")
	}
	return utils.SuccessResponse(c, u)
}

// GetUser
// @Summary 사용자 조회
// @Description ID로 사용자를 조회합니다
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "사용자 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users/{id} [get]
func (ctl *Controller) GetUser(c *fiber.Ctx) error {
	var path types.Id
	path.ID = c.Params("id")
	if err := utils.ValidateStruct(path); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}
	u, err := ctl.service.GetByID(path.ID)
	if err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}
	return utils.SuccessResponse(c, u)
}
