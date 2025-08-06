package strategy

import (
	"context"
	"fmt"
	"strings"

	"auto-trader/pkg/shared/config"
	"auto-trader/pkg/shared/middleware"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// DynamicStrategy 완전 DB 기반 동적 전략 엔진
type DynamicStrategy struct {
	// 의존성들
	dataCollector Collector
	executor      Executor
	riskManager   *middleware.Manager
	appConfig     *config.Config

	// 전략 설정 (DB에서 로드)
	strategyConfig *StrategyConfig

	// 런타임 상태
	stopChan chan struct{}
}

// NewDynamicStrategy 동적 전략 생성
func NewDynamicStrategy(
	dataCollector Collector,
	executor Executor,
	riskManager *middleware.Manager,
	appConfig *config.Config,
	strategyConfig *StrategyConfig,
) Strategy {
	return &DynamicStrategy{
		dataCollector:  dataCollector,
		executor:       executor,
		riskManager:    riskManager,
		appConfig:      appConfig,
		strategyConfig: strategyConfig,
		stopChan:       make(chan struct{}),
	}
}

func (s *DynamicStrategy) ID() string {
	return s.strategyConfig.ID
}

func (s *DynamicStrategy) Name() string {
	if name, ok := s.strategyConfig.Parameters["name"].(string); ok {
		return name
	}
	return "동적 전략"
}

func (s *DynamicStrategy) Symbols() []string {
	if symbols, ok := s.strategyConfig.Parameters["symbols"].([]interface{}); ok {
		result := make([]string, len(symbols))
		for i, symbol := range symbols {
			result[i] = symbol.(string)
		}
		return result
	}
	return []string{}
}

func (s *DynamicStrategy) Execute() error {
	if !s.strategyConfig.Enabled {
		return nil
	}

	// DB에서 로드한 전략 로직을 동적으로 실행
	return s.executeStrategyLogic()
}

func (s *DynamicStrategy) Start() error {
	logrus.Infof("🚀 동적 전략 시작: %s", s.Name())
	return nil
}

func (s *DynamicStrategy) Stop() error {
	close(s.stopChan)
	logrus.Infof("🛑 동적 전략 중지: %s", s.Name())
	return nil
}

// executeStrategyLogic DB에 저장된 전략 로직을 동적으로 실행
func (s *DynamicStrategy) executeStrategyLogic() error {
	symbols := s.Symbols()

	for _, symbol := range symbols {
		if err := s.executeForSymbol(symbol); err != nil {
			logrus.Errorf("전략 실행 오류 (%s): %v", symbol, err)
		}
	}
	return nil
}

// executeForSymbol 특정 심볼에 대해 전략 실행
func (s *DynamicStrategy) executeForSymbol(symbol string) error {
	// 현재가 조회
	priceData, err := s.dataCollector.GetCurrentPrice(symbol)
	if err != nil {
		return fmt.Errorf("현재가 조회 실패: %w", err)
	}

	// 전략 조건들을 DB에서 로드하여 평가
	conditions := s.getConditions()

	for _, condition := range conditions {
		if s.evaluateCondition(condition, symbol, priceData) {
			if err := s.executeAction(condition.Action, symbol, priceData); err != nil {
				return fmt.Errorf("액션 실행 실패: %w", err)
			}
		}
	}

	return nil
}

// Condition 전략 조건 구조체
type Condition struct {
	Type     string                 `json:"type"`
	Operator string                 `json:"operator"`
	Value    interface{}            `json:"value"`
	Action   Action                 `json:"action"`
	Priority int                    `json:"priority"`
	Metadata map[string]interface{} `json:"metadata"`
}

// Action 전략 액션 구조체
type Action struct {
	Type     string                 `json:"type"`     // "BUY", "SELL", "HOLD"
	Quantity interface{}            `json:"quantity"` // 숫자 또는 "ALL", "PERCENTAGE"
	Price    interface{}            `json:"price"`    // "MARKET", "LIMIT", 숫자
	Metadata map[string]interface{} `json:"metadata"`
}

// getConditions DB에서 조건들을 로드
func (s *DynamicStrategy) getConditions() []Condition {
	conditions := []Condition{}

	if conditionsData, ok := s.strategyConfig.Parameters["conditions"].([]interface{}); ok {
		for _, condData := range conditionsData {
			if condMap, ok := condData.(map[string]interface{}); ok {
				condition := Condition{
					Type:     s.getString(condMap, "type"),
					Operator: s.getString(condMap, "operator"),
					Value:    condMap["value"],
					Action: Action{
						Type:     s.getString(condMap, "action_type"),
						Quantity: condMap["action_quantity"],
						Price:    condMap["action_price"],
					},
					Priority: s.getInt(condMap, "priority", 0),
				}
				conditions = append(conditions, condition)
			}
		}
	}

	return conditions
}

// evaluateCondition 조건 평가
func (s *DynamicStrategy) evaluateCondition(condition Condition, symbol string, priceData *PriceData) bool {
	switch condition.Type {
	case "profit_percentage":
		return s.evaluateProfitCondition(condition, symbol, priceData)
	case "daily_profit":
		return s.evaluateDailyProfitCondition(condition, symbol, priceData)
	case "price_level":
		return s.evaluatePriceCondition(condition, symbol, priceData)
	case "rsi":
		return s.evaluateRSICondition(condition, symbol, priceData)
	case "moving_average":
		return s.evaluateMACondition(condition, symbol, priceData)
	case "bollinger_bands":
		return s.evaluateBBCondition(condition, symbol, priceData)
	default:
		logrus.Warnf("지원하지 않는 조건 타입: %s", condition.Type)
		return false
	}
}

// evaluateProfitCondition 수익률 조건 평가
func (s *DynamicStrategy) evaluateProfitCondition(condition Condition, symbol string, priceData *PriceData) bool {
	// 실제로는 포지션 정보에서 총 수익률 계산 필요
	totalProfit := 0.0 // TODO: 실제 포지션에서 계산

	threshold := s.getFloat(condition.Value, 0.0)

	switch condition.Operator {
	case ">=":
		return totalProfit >= threshold
	case "<=":
		return totalProfit <= threshold
	case ">":
		return totalProfit > threshold
	case "<":
		return totalProfit < threshold
	default:
		return false
	}
}

// evaluateDailyProfitCondition 일일 수익률 조건 평가
func (s *DynamicStrategy) evaluateDailyProfitCondition(condition Condition, symbol string, priceData *PriceData) bool {
	dailyProfit, err := s.dataCollector.GetDailyProfit(symbol)
	if err != nil {
		logrus.Errorf("일일 수익률 조회 실패: %v", err)
		return false
	}

	threshold := s.getFloat(condition.Value, 0.0)
	dailyProfitFloat, _ := dailyProfit.Float64()

	switch condition.Operator {
	case ">=":
		return dailyProfitFloat >= threshold
	case "<=":
		return dailyProfitFloat <= threshold
	case ">":
		return dailyProfitFloat > threshold
	case "<":
		return dailyProfitFloat < threshold
	default:
		return false
	}
}

// evaluatePriceCondition 가격 조건 평가
func (s *DynamicStrategy) evaluatePriceCondition(condition Condition, symbol string, priceData *PriceData) bool {
	currentPrice, _ := priceData.Price.Float64()
	threshold := s.getFloat(condition.Value, 0.0)

	switch condition.Operator {
	case ">=":
		return currentPrice >= threshold
	case "<=":
		return currentPrice <= threshold
	case ">":
		return currentPrice > threshold
	case "<":
		return currentPrice < threshold
	default:
		return false
	}
}

// evaluateRSICondition RSI 조건 평가 (임시)
func (s *DynamicStrategy) evaluateRSICondition(condition Condition, symbol string, priceData *PriceData) bool {
	rsi := 50.0 // TODO: 실제 RSI 계산
	threshold := s.getFloat(condition.Value, 0.0)

	switch condition.Operator {
	case ">=":
		return rsi >= threshold
	case "<=":
		return rsi <= threshold
	default:
		return false
	}
}

// evaluateMACondition 이동평균 조건 평가 (임시)
func (s *DynamicStrategy) evaluateMACondition(condition Condition, symbol string, priceData *PriceData) bool {
	shortMA := 100.0 // TODO: 실제 단기 이동평균 계산
	longMA := 105.0  // TODO: 실제 장기 이동평균 계산

	switch condition.Operator {
	case "CROSS_ABOVE":
		return shortMA > longMA
	case "CROSS_BELOW":
		return shortMA < longMA
	default:
		return false
	}
}

// evaluateBBCondition 볼린저 밴드 조건 평가 (임시)
func (s *DynamicStrategy) evaluateBBCondition(condition Condition, symbol string, priceData *PriceData) bool {
	currentPrice, _ := priceData.Price.Float64()
	upperBand := 110.0 // TODO: 실제 상단 밴드 계산
	lowerBand := 90.0  // TODO: 실제 하단 밴드 계산

	switch condition.Operator {
	case "TOUCH_UPPER":
		return currentPrice >= upperBand
	case "TOUCH_LOWER":
		return currentPrice <= lowerBand
	default:
		return false
	}
}

// executeAction 액션 실행
func (s *DynamicStrategy) executeAction(action Action, symbol string, priceData *PriceData) error {
	switch action.Type {
	case "BUY":
		return s.executeBuyAction(action, symbol, priceData)
	case "SELL":
		return s.executeSellAction(action, symbol, priceData)
	case "HOLD":
		logrus.Infof("📊 홀드: %s", symbol)
		return nil
	default:
		return fmt.Errorf("지원하지 않는 액션 타입: %s", action.Type)
	}
}

// executeBuyAction 매수 액션 실행
func (s *DynamicStrategy) executeBuyAction(action Action, symbol string, priceData *PriceData) error {
	quantity := s.calculateQuantity(action.Quantity, priceData.Price)
	orderPrice := s.calculatePrice(action.Price, priceData.Price)

	_, err := s.executor.ExecuteOrder(context.Background(), symbol, "BUY", quantity, orderPrice, "MARKET")
	if err != nil {
		return fmt.Errorf("매수 주문 실패: %w", err)
	}

	logrus.Infof("📈 매수 실행: %s, 수량: %s, 가격: %s", symbol, quantity.String(), orderPrice.String())
	return nil
}

// executeSellAction 매도 액션 실행
func (s *DynamicStrategy) executeSellAction(action Action, symbol string, priceData *PriceData) error {
	quantity := s.calculateQuantity(action.Quantity, priceData.Price)
	orderPrice := s.calculatePrice(action.Price, priceData.Price)

	_, err := s.executor.ExecuteOrder(context.Background(), symbol, "SELL", quantity, orderPrice, "MARKET")
	if err != nil {
		return fmt.Errorf("매도 주문 실패: %w", err)
	}

	logrus.Infof("📉 매도 실행: %s, 수량: %s, 가격: %s", symbol, quantity.String(), orderPrice.String())
	return nil
}

// calculateQuantity 수량 계산
func (s *DynamicStrategy) calculateQuantity(quantity interface{}, price decimal.Decimal) decimal.Decimal {
	switch q := quantity.(type) {
	case float64:
		return decimal.NewFromFloat(q).Div(price)
	case string:
		if q == "ALL" {
			return decimal.NewFromFloat(1000.0).Div(price) // TODO: 실제 보유 수량 조회
		} else if strings.HasSuffix(q, "%") {
			percent := s.getFloat(q, 100.0)
			return decimal.NewFromFloat(percent / 100.0 * 1000.0).Div(price) // TODO: 실제 보유 수량 조회
		}
	}
	return decimal.NewFromFloat(1000.0).Div(price) // 기본값
}

// calculatePrice 주문 가격 계산
func (s *DynamicStrategy) calculatePrice(price interface{}, currentPrice decimal.Decimal) decimal.Decimal {
	switch p := price.(type) {
	case float64:
		return decimal.NewFromFloat(p)
	case string:
		if p == "MARKET" {
			return currentPrice
		}
	}
	return currentPrice // 기본값
}

// Helper methods
func (s *DynamicStrategy) getString(value interface{}, key string) string {
	if m, ok := value.(map[string]interface{}); ok {
		if v, ok := m[key].(string); ok {
			return v
		}
	}
	return ""
}

func (s *DynamicStrategy) getFloat(value interface{}, defaultValue float64) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		if f, err := decimal.NewFromString(v); err == nil {
			f64, _ := f.Float64()
			return f64
		}
	}
	return defaultValue
}

func (s *DynamicStrategy) getInt(value interface{}, key string, defaultValue int) int {
	if m, ok := value.(map[string]interface{}); ok {
		if v, ok := m[key].(int); ok {
			return v
		}
	}
	return defaultValue
}
