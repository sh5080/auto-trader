package strategy

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Controller 전략 HTTP 컨트롤러
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
func (ctrl *Controller) GetAllStrategies(c *fiber.Ctx) error {
	strategies, err := ctrl.service.GetAllStrategies()
	if err != nil {
		logrus.Errorf("전략 목록 조회 실패: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "전략 목록을 조회할 수 없습니다",
		})
	}

	activeCount := 0
	for _, strategy := range strategies {
		if strategy.Active {
			activeCount++
		}
	}

	return c.JSON(StrategyListResponse{
		Success:    true,
		Strategies: strategies,
		Count:      len(strategies),
		Active:     activeCount,
	})
}

// GetStrategy 특정 전략 조회
func (ctrl *Controller) GetStrategy(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "전략 ID가 필요합니다",
		})
	}

	strategy, err := ctrl.service.GetStrategy(id)
	if err != nil {
		logrus.Errorf("전략 조회 실패 (ID: %s): %v", id, err)
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "전략을 찾을 수 없습니다",
		})
	}

	return c.JSON(StrategyResponse{
		Success:  true,
		Strategy: strategy,
		Message:  "전략 조회 성공",
	})
}

// GetStrategyStatus 전략 상태 조회
func (ctrl *Controller) GetStrategyStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "전략 ID가 필요합니다",
		})
	}

	status, err := ctrl.service.GetStrategyStatus(id)
	if err != nil {
		logrus.Errorf("전략 상태 조회 실패 (ID: %s): %v", id, err)
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "전략 상태를 조회할 수 없습니다",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"status":  status,
	})
}

// GetStrategyPerformance 전략 성과 조회
func (ctrl *Controller) GetStrategyPerformance(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "전략 ID가 필요합니다",
		})
	}

	performance, err := ctrl.service.GetStrategyPerformance(id)
	if err != nil {
		logrus.Errorf("전략 성과 조회 실패 (ID: %s): %v", id, err)
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"error":   "전략 성과를 조회할 수 없습니다",
		})
	}

	return c.JSON(fiber.Map{
		"success":     true,
		"performance": performance,
	})
}

// StartStrategy 전략 시작
func (ctrl *Controller) StartStrategy(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "전략 ID가 필요합니다",
		})
	}

	err := ctrl.service.StartStrategy(id)
	if err != nil {
		logrus.Errorf("전략 시작 실패 (ID: %s): %v", id, err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "전략을 시작할 수 없습니다",
		})
	}

	logrus.Infof("전략 시작됨 (ID: %s)", id)
	return c.JSON(StrategyResponse{
		Success: true,
		Message: "전략이 성공적으로 시작되었습니다",
	})
}

// StopStrategy 전략 중지
func (ctrl *Controller) StopStrategy(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "전략 ID가 필요합니다",
		})
	}

	err := ctrl.service.StopStrategy(id)
	if err != nil {
		logrus.Errorf("전략 중지 실패 (ID: %s): %v", id, err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "전략을 중지할 수 없습니다",
		})
	}

	logrus.Infof("전략 중지됨 (ID: %s)", id)
	return c.JSON(StrategyResponse{
		Success: true,
		Message: "전략이 성공적으로 중지되었습니다",
	})
}

// RestartStrategy 전략 재시작
func (ctrl *Controller) RestartStrategy(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "전략 ID가 필요합니다",
		})
	}

	err := ctrl.service.RestartStrategy(id)
	if err != nil {
		logrus.Errorf("전략 재시작 실패 (ID: %s): %v", id, err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "전략을 재시작할 수 없습니다",
		})
	}

	logrus.Infof("전략 재시작됨 (ID: %s)", id)
	return c.JSON(StrategyResponse{
		Success: true,
		Message: "전략이 성공적으로 재시작되었습니다",
	})
}

// CreateStrategy 전략 생성
func (ctrl *Controller) CreateStrategy(c *fiber.Ctx) error {
	var req CreateStrategyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "잘못된 요청 형식입니다",
		})
	}

	strategy, err := ctrl.service.CreateStrategy(&req)
	if err != nil {
		logrus.Errorf("전략 생성 실패: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "전략을 생성할 수 없습니다",
		})
	}

	logrus.Infof("새 전략 생성됨 (ID: %s, Name: %s)", strategy.ID, strategy.Name)
	return c.Status(201).JSON(StrategyResponse{
		Success:  true,
		Strategy: strategy,
		Message:  "전략이 성공적으로 생성되었습니다",
	})
}

// UpdateStrategy 전략 수정
func (ctrl *Controller) UpdateStrategy(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "전략 ID가 필요합니다",
		})
	}

	var req UpdateStrategyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "잘못된 요청 형식입니다",
		})
	}

	strategy, err := ctrl.service.UpdateStrategy(id, &req)
	if err != nil {
		logrus.Errorf("전략 수정 실패 (ID: %s): %v", id, err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "전략을 수정할 수 없습니다",
		})
	}

	logrus.Infof("전략 수정됨 (ID: %s)", id)
	return c.JSON(StrategyResponse{
		Success:  true,
		Strategy: strategy,
		Message:  "전략이 성공적으로 수정되었습니다",
	})
}

// DeleteStrategy 전략 삭제
func (ctrl *Controller) DeleteStrategy(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "전략 ID가 필요합니다",
		})
	}

	err := ctrl.service.DeleteStrategy(id)
	if err != nil {
		logrus.Errorf("전략 삭제 실패 (ID: %s): %v", id, err)
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   "전략을 삭제할 수 없습니다",
		})
	}

	logrus.Infof("전략 삭제됨 (ID: %s)", id)
	return c.JSON(StrategyResponse{
		Success: true,
		Message: "전략이 성공적으로 삭제되었습니다",
	})
}
