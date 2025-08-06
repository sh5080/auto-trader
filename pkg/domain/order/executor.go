package execution

import (
	"context"
	"fmt"
	"sync"
	"time"

	"auto-trader/pkg/api"
	"auto-trader/pkg/shared/config"
	"auto-trader/pkg/shared/middleware"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type Executor struct {
	apiClient   *api.Client
	riskManager *middleware.Manager
	config      *config.Config
	orders      map[string]*api.OrderResponse
	mutex       sync.RWMutex
}

func NewExecutor(apiClient *api.Client, riskManager *middleware.Manager, cfg *config.Config) *Executor {
	return &Executor{
		apiClient:   apiClient,
		riskManager: riskManager,
		config:      cfg,
		orders:      make(map[string]*api.OrderResponse),
	}
}

func (e *Executor) ExecuteOrder(ctx context.Context, symbol string, side string, quantity decimal.Decimal, price decimal.Decimal, orderType string) (*api.OrderResponse, error) {
	// 리스크 체크
	riskCheck := e.riskManager.CheckOrderRisk(symbol, side, quantity, price)
	if !riskCheck.Allowed {
		return nil, fmt.Errorf("리스크 체크 실패: %s", riskCheck.Reason)
	}

	// 주문 요청 생성
	orderReq := &api.OrderRequest{
		Symbol:   symbol,
		Side:     side,
		Quantity: quantity,
		Type:     orderType,
	}

	// 시장가 주문이 아닌 경우 가격 설정
	if orderType != "market" {
		orderReq.Price = price
	}

	// 주문 실행
	var orderResp *api.OrderResponse
	var err error

	for attempt := 0; attempt < e.config.Trading.RetryAttempts; attempt++ {
		orderResp, err = e.apiClient.PlaceOrder(ctx, orderReq)
		if err == nil {
			break
		}

		logrus.Warnf("주문 재시도 %d/%d: %v", attempt+1, e.config.Trading.RetryAttempts, err)
		if attempt < e.config.Trading.RetryAttempts-1 {
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("주문 실행 실패: %w", err)
	}

	// 주문 기록 저장
	e.mutex.Lock()
	e.orders[orderResp.ID] = orderResp
	e.mutex.Unlock()

	// 포지션 업데이트
	if orderResp.Status == "filled" {
		e.riskManager.UpdatePosition(symbol, side, quantity, orderResp.Price)
	}

	logrus.Infof("주문 실행 완료: %s %s %s @ %s (ID: %s)",
		side, quantity.String(), symbol, orderResp.Price.String(), orderResp.ID)

	return orderResp, nil
}

func (e *Executor) CancelOrder(orderID string) error {
	err := e.apiClient.CancelOrder(orderID)
	if err != nil {
		return fmt.Errorf("주문 취소 실패: %w", err)
	}

	// 주문 상태 업데이트
	e.mutex.Lock()
	if order, exists := e.orders[orderID]; exists {
		order.Status = "cancelled"
	}
	e.mutex.Unlock()

	return nil
}

func (e *Executor) GetOrderStatus(orderID string) (*api.OrderResponse, error) {
	orderResp, err := e.apiClient.GetOrderStatus(orderID)
	if err != nil {
		return nil, fmt.Errorf("주문 상태 조회 실패: %w", err)
	}

	// 로컬 캐시 업데이트
	e.mutex.Lock()
	e.orders[orderID] = orderResp
	e.mutex.Unlock()

	return orderResp, nil
}

func (e *Executor) GetOrders() map[string]*api.OrderResponse {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	result := make(map[string]*api.OrderResponse)
	for k, v := range e.orders {
		result[k] = v
	}
	return result
}

func (e *Executor) ExecuteStopLoss(symbol string, currentPrice decimal.Decimal) error {
	positions := e.riskManager.GetPositions()
	if position, exists := positions[symbol]; exists {
		// 스탑로스 조건 확인
		if e.riskManager.CheckStopLoss(symbol, currentPrice) {
			// 반대 방향 주문으로 청산
			side := "sell"
			if position.Side == "short" {
				side = "buy"
			}

			ctx, cancel := context.WithTimeout(context.Background(), e.config.Trading.OrderTimeout)
			defer cancel()

			_, err := e.ExecuteOrder(ctx, symbol, side, position.Quantity, currentPrice, "market")
			if err != nil {
				return fmt.Errorf("스탑로스 실행 실패: %w", err)
			}

			logrus.Warnf("스탑로스 실행: %s %s %s @ %s",
				side, position.Quantity.String(), symbol, currentPrice.String())
		}
	}

	return nil
}

func (e *Executor) ExecuteTakeProfit(symbol string, currentPrice decimal.Decimal, takeProfitPercentage float64) error {
	positions := e.riskManager.GetPositions()
	if position, exists := positions[symbol]; exists {
		takeProfitThreshold := position.AvgPrice.Mul(decimal.NewFromFloat(1 + takeProfitPercentage))

		if position.Side == "long" && currentPrice.GreaterThanOrEqual(takeProfitThreshold) {
			// 롱 포지션 익절
			ctx, cancel := context.WithTimeout(context.Background(), e.config.Trading.OrderTimeout)
			defer cancel()

			_, err := e.ExecuteOrder(ctx, symbol, "sell", position.Quantity, currentPrice, "market")
			if err != nil {
				return fmt.Errorf("익절 실행 실패: %w", err)
			}

			logrus.Infof("익절 실행: sell %s %s @ %s",
				position.Quantity.String(), symbol, currentPrice.String())
		}
	}

	return nil
}
