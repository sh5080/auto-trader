package portfolio

import (
	"strings"

	"auto-trader/pkg/domain/portfolio/dto"
	"auto-trader/pkg/shared/utils"

	"github.com/gofiber/fiber/v2"
)

// Controller 포트폴리오 컨트롤러
type Controller struct {
	service Service
}

// NewController 새로운 포트폴리오 컨트롤러 생성
func NewController(service Service) *Controller {
	return &Controller{
		service: service,
	}
}

// GetPortfolio 포트폴리오 조회
// @Summary 포트폴리오 조회
// @Description 사용자의 포트폴리오 정보를 조회합니다
// @Tags portfolio
// @Accept json
// @Produce json
// @Param forceRefresh query bool false "강제 새로고침 여부"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /portfolio [get]
func (ctrl *Controller) GetPortfolio(c *fiber.Ctx) error {
	var q dto.GetPortfolioQuery
	userID := utils.GetUserID(c)
	q.ForceRefresh = c.QueryBool("forceRefresh", false)
	if err := utils.ValidateStruct(q); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	portfolio, err := ctrl.service.GetPortfolio(userID, q)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "포트폴리오 조회 실패", err)
	}

	return utils.SuccessResponse(c, portfolio)
}

// GetPortfolioSummary 포트폴리오 요약 조회
// @Summary 포트폴리오 요약 조회
// @Description 사용자의 포트폴리오 요약 정보를 조회합니다
// @Tags portfolio
// @Accept json
// @Produce json
// @Param forceRefresh query bool false "강제 새로고침 여부"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /portfolio/summary [get]
func (ctrl *Controller) GetPortfolioSummary(c *fiber.Ctx) error {
	var q dto.GetPortfolioSummaryQuery
	userID := utils.GetUserID(c)
	q.ForceRefresh = c.QueryBool("forceRefresh", false)
	if err := utils.ValidateStruct(q); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	summary, err := ctrl.service.GetPortfolioSummary(userID, q)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "포트폴리오 요약 조회 실패", err)
	}

	return utils.SuccessResponse(c, summary)
}

// GetPositions 보유 주식 목록 조회
// @Summary 보유 주식 목록 조회
// @Description 사용자의 보유 주식 목록을 조회합니다
// @Tags portfolio
// @Accept json
// @Produce json
// @Param symbol query string false "특정 종목 심볼"
// @Param forceRefresh query bool false "강제 새로고침 여부"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /portfolio/positions [get]
func (ctrl *Controller) GetPositions(c *fiber.Ctx) error {
	var q dto.GetPositionsQuery
	userID := utils.GetUserID(c)
	q.Symbol = c.Query("symbol")
	q.ForceRefresh = c.QueryBool("forceRefresh", false)
	if err := utils.ValidateStruct(q); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	positions, err := ctrl.service.GetPositions(userID, q)
	if err != nil {
		if utils.IsValidationError(err) {
			validationErr := utils.UnwrapValidationError(err)
			return utils.ValidationErrorResponse(c, validationErr.Error())
		}
		return utils.InternalServerErrorResponse(c, "보유 주식 조회 실패", err)
	}

	return utils.SuccessResponse(c, positions)
}

// GetPosition 특정 보유 주식 조회
// @Summary 특정 보유 주식 조회
// @Description 사용자의 특정 보유 주식을 조회합니다
// @Tags portfolio
// @Accept json
// @Produce json
// @Param symbol path string true "종목 심볼"
// @Param forceRefresh query bool false "강제 새로고침 여부"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /portfolio/positions/{symbol} [get]
func (ctrl *Controller) GetPosition(c *fiber.Ctx) error {
	var path dto.SymbolPath
	path.Symbol = c.Params("symbol")
	if err := utils.ValidateStruct(path); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	userID := utils.GetUserID(c)

	position, err := ctrl.service.GetPosition(userID, dto.GetPositionsQuery{
		Symbol: path.Symbol,
	})
	if err != nil {
		return utils.NotFoundResponse(c, "보유 주식을 찾을 수 없습니다")
	}

	return utils.SuccessResponse(c, position)
}

// GetCurrentPrice 현재가 조회
// @Summary 현재가 조회
// @Description 특정 종목의 현재가를 조회합니다
// @Tags portfolio
// @Accept json
// @Produce json
// @Param symbol path string true "종목 심볼"
// @Param forceRefresh query bool false "강제 새로고침 여부"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /portfolio/prices/{symbol} [get]
func (ctrl *Controller) GetCurrentPrice(c *fiber.Ctx) error {
	var q dto.GetCurrentPricesQuery
	q.Symbols = c.Params("symbol")
	if err := utils.ValidateStruct(q); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	price, err := ctrl.service.GetCurrentPrice(q)
	if err != nil {
		if utils.IsValidationError(err) {
			validationErr := utils.UnwrapValidationError(err)
			return utils.ValidationErrorResponse(c, validationErr.Error())
		}
		return utils.InternalServerErrorResponse(c, "현재가 조회 실패", err)
	}

	return utils.SuccessResponse(c, price)
}

// GetCurrentPrices 여러 종목 현재가 조회
// @Summary 여러 종목 현재가 조회
// @Description 여러 종목의 현재가를 조회합니다
// @Tags portfolio
// @Accept json
// @Produce json
// @Param symbols query string true "종목 심볼들 (쉼표로 구분)"
// @Param forceRefresh query bool false "강제 새로고침 여부"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /portfolio/prices [get]
func (ctrl *Controller) GetCurrentPrices(c *fiber.Ctx) error {
	var q dto.GetCurrentPricesQuery
	q.Symbols = c.Query("symbols")
	q.ForceRefresh = c.QueryBool("forceRefresh", false)
	if err := utils.ValidateStruct(q); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	symbols := strings.Split(q.Symbols, ",")
	if len(symbols) == 0 {
		return utils.ValidationErrorResponse(c, "유효한 종목 심볼이 없습니다")
	}

	prices, err := ctrl.service.GetCurrentPrices(q)
	if err != nil {
		if utils.IsValidationError(err) {
			validationErr := utils.UnwrapValidationError(err)
			return utils.ValidationErrorResponse(c, validationErr.Error())
		}
		return utils.InternalServerErrorResponse(c, "현재가 조회 실패", err)
	}

	return utils.SuccessResponse(c, prices)
}

// GetTradeHistory 거래 내역 조회
// @Summary 거래 내역 조회
// @Description 사용자의 거래 내역을 조회합니다
// @Tags portfolio
// @Accept json
// @Produce json
// @Param symbol query string false "특정 종목 심볼"
// @Param start_date query string false "시작 날짜 (YYYY-MM-DD)"
// @Param end_date query string false "종료 날짜 (YYYY-MM-DD)"
// @Param limit query int false "조회 개수 제한"
// @Param offset query int false "오프셋"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /portfolio/trades [get]
func (ctrl *Controller) GetTradeHistory(c *fiber.Ctx) error {
	var q dto.GetTradeHistoryQuery
	userID := utils.GetUserID(c)
	q.Symbol = c.Query("symbol")
	q.Limit = c.QueryInt("limit", 100)
	q.Offset = c.QueryInt("offset", 0)
	if err := utils.ValidateStruct(q); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	trades, err := ctrl.service.GetTradeHistory(userID, q)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "거래 내역 조회 실패", err)
	}

	return utils.SuccessResponse(c, trades)
}

// RefreshPortfolio 포트폴리오 새로고침
// @Summary 포트폴리오 새로고침
// @Description 포트폴리오 데이터를 강제로 새로고침합니다
// @Tags portfolio
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /portfolio/refresh [post]
func (ctrl *Controller) RefreshPortfolio(c *fiber.Ctx) error {
	userID := utils.GetUserID(c)
	if err := ctrl.service.RefreshPortfolio(userID); err != nil {
		return utils.InternalServerErrorResponse(c, "포트폴리오 새로고침 실패", err)
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "포트폴리오가 성공적으로 새로고침되었습니다",
		"userID":  userID,
	})
}
