package strategy

import (
	"auto-trader/pkg/domain/strategy/dto"
	"auto-trader/pkg/shared/types"
	"auto-trader/pkg/shared/utils"

	"github.com/gofiber/fiber/v2"
)

// Controller 전략 컨트롤러
type Controller struct {
	service Service
}

// NewController 새로운 전략 컨트롤러 생성
func NewController(service Service) *Controller {
	return &Controller{
		service: service,
	}
}

// GetAllStrategies 모든 전략 조회
// @Summary 모든 전략 조회
// @Description 사용자의 모든 전략 목록을 조회합니다
// @Tags strategies
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /strategies [get]
func (ctrl *Controller) GetAllStrategies(c *fiber.Ctx) error {
	strategies, err := ctrl.service.GetAllStrategies()
	if err != nil {
		return utils.InternalServerErrorResponse(c, "전략 목록 조회 실패", err)
	}

	return utils.SuccessResponse(c, strategies)
}

// GetStrategy 특정 전략 조회
// @Summary 특정 전략 조회
// @Description ID로 특정 전략을 조회합니다
// @Tags strategies
// @Accept json
// @Produce json
// @Param id path string true "전략 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /strategies/{id} [get]
func (ctrl *Controller) GetStrategy(c *fiber.Ctx) error {
	var path types.Id
	path.ID = c.Params("id")
	if err := utils.ValidateStruct(path); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	strategy, err := ctrl.service.GetStrategy(path.ID)
	if err != nil {
		return utils.NotFoundResponse(c, "전략을 찾을 수 없습니다")
	}

	return utils.SuccessResponse(c, strategy)
}

// CreateStrategy 전략 생성
// @Summary 전략 생성
// @Description 새로운 전략을 생성합니다
// @Tags strategies
// @Accept json
// @Produce json
// @Param strategy body dto.CreateStrategyBody true "전략 정보"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /strategies [post]
func (ctrl *Controller) CreateStrategy(c *fiber.Ctx) error {
	var req dto.CreateStrategyBody
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "잘못된 요청 형식")
	}
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	strategy, err := ctrl.service.CreateStrategy(&req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "전략 생성 실패", err)
	}

	return utils.SuccessResponse(c, strategy, fiber.StatusCreated)
}

// UpdateStrategy 전략 수정
// @Summary 전략 수정
// @Description 전략을 수정합니다
// @Tags strategies
// @Accept json
// @Produce json
// @Param id path string true "전략 ID"
// @Param strategy body dto.UpdateStrategyInput true "수정할 전략 정보"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /strategies/{id} [put]
func (ctrl *Controller) UpdateStrategy(c *fiber.Ctx) error {
	var path types.Id
	path.ID = c.Params("id")
	if err := utils.ValidateStruct(path); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	var req dto.UpdateStrategyBody
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "잘못된 요청 형식")
	}
	// Update는 부분 업데이트 허용이므로 필수값 검증은 스킵하거나 필요한 필드만 검증

	strategy, err := ctrl.service.UpdateStrategy(path.ID, &req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "전략 수정 실패", err)
	}

	return utils.SuccessResponse(c, strategy)
}

// DeleteStrategy 전략 삭제
// @Summary 전략 삭제
// @Description 전략을 삭제합니다
// @Tags strategies
// @Accept json
// @Produce json
// @Param id path string true "전략 ID"
// @Success 204 "No Content"
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /strategies/{id} [delete]
func (ctrl *Controller) DeleteStrategy(c *fiber.Ctx) error {
	var path types.Id
	path.ID = c.Params("id")
	if err := utils.ValidateStruct(path); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	if err := ctrl.service.DeleteStrategy(path.ID); err != nil {
		return utils.InternalServerErrorResponse(c, "전략 삭제 실패", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetStrategyStatus 전략 상태 조회
// @Summary 전략 상태 조회
// @Description 전략의 현재 상태를 조회합니다
// @Tags strategies
// @Accept json
// @Produce json
// @Param id path string true "전략 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /strategies/{id}/status [get]
func (ctrl *Controller) GetStrategyStatus(c *fiber.Ctx) error {
	var path types.Id
	path.ID = c.Params("id")
	if err := utils.ValidateStruct(path); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	status, err := ctrl.service.GetStrategyStatus(path.ID)
	if err != nil {
		return utils.NotFoundResponse(c, "전략 상태를 찾을 수 없습니다")
	}

	return utils.SuccessResponse(c, status)
}

// StartStrategy 전략 시작
// @Summary 전략 시작
// @Description 전략을 시작합니다
// @Tags strategies
// @Accept json
// @Produce json
// @Param id path string true "전략 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /strategies/{id}/start [post]
func (ctrl *Controller) StartStrategy(c *fiber.Ctx) error {
	var path types.Id
	path.ID = c.Params("id")
	if err := utils.ValidateStruct(path); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	if err := ctrl.service.StartStrategy(path.ID); err != nil {
		return utils.InternalServerErrorResponse(c, "전략 시작 실패", err)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message":     "전략이 시작되었습니다",
		"strategy_id": path.ID,
	})
}

// StopStrategy 전략 중지
// @Summary 전략 중지
// @Description 전략을 중지합니다
// @Tags strategies
// @Accept json
// @Produce json
// @Param id path string true "전략 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /strategies/{id}/stop [post]
func (ctrl *Controller) StopStrategy(c *fiber.Ctx) error {
	var path types.Id
	path.ID = c.Params("id")
	if err := utils.ValidateStruct(path); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	if err := ctrl.service.StopStrategy(path.ID); err != nil {
		return utils.InternalServerErrorResponse(c, "전략 중지 실패", err)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message":     "전략이 중지되었습니다",
		"strategy_id": path.ID,
	})
}

// RestartStrategy 전략 재시작
// @Summary 전략 재시작
// @Description 전략을 재시작합니다
// @Tags strategies
// @Accept json
// @Produce json
// @Param id path string true "전략 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /strategies/{id}/restart [post]
func (ctrl *Controller) RestartStrategy(c *fiber.Ctx) error {
	var path types.Id
	path.ID = c.Params("id")
	if err := utils.ValidateStruct(path); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	if err := ctrl.service.RestartStrategy(path.ID); err != nil {
		return utils.InternalServerErrorResponse(c, "전략 재시작 실패", err)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message":     "전략이 재시작되었습니다",
		"strategy_id": path.ID,
	})
}

// GetStrategyPerformance 전략 성과 조회
// @Summary 전략 성과 조회
// @Description 전략의 성과 정보를 조회합니다
// @Tags strategies
// @Accept json
// @Produce json
// @Param id path string true "전략 ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /strategies/{id}/performance [get]
func (ctrl *Controller) GetStrategyPerformance(c *fiber.Ctx) error {
	var path types.Id
	path.ID = c.Params("id")
	if err := utils.ValidateStruct(path); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	performance, err := ctrl.service.GetStrategyPerformance(path.ID)
	if err != nil {
		return utils.NotFoundResponse(c, "전략 성과를 찾을 수 없습니다")
	}

	return utils.SuccessResponse(c, performance)
}
