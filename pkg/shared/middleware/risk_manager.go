package middleware

import (
	"sync"
	"time"

	"auto-trader/pkg/shared/config"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	config        *config.Config
	dailyLoss     decimal.Decimal
	positions     map[string]*Position
	mutex         sync.RWMutex
	lastResetDate time.Time
}

type Position struct {
	Symbol    string          `json:"symbol"`
	Quantity  decimal.Decimal `json:"quantity"`
	AvgPrice  decimal.Decimal `json:"avg_price"`
	Side      string          `json:"side"` // "long" or "short"
	Timestamp time.Time       `json:"timestamp"`
}

type RiskCheck struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		config:        cfg,
		dailyLoss:     decimal.Zero,
		positions:     make(map[string]*Position),
		lastResetDate: time.Now(),
	}
}

func (m *Manager) CheckOrderRisk(symbol string, side string, quantity decimal.Decimal, price decimal.Decimal) *RiskCheck {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 일일 손실 한도 체크
	if m.dailyLoss.GreaterThanOrEqual(decimal.NewFromFloat(m.config.Risk.MaxDailyLoss)) {
		return &RiskCheck{
			Allowed: false,
			Reason:  "일일 손실 한도 초과",
		}
	}

	// 포지션 크기 체크
	orderValue := quantity.Mul(price)
	if orderValue.GreaterThan(decimal.NewFromFloat(m.config.Risk.MaxPositionSize)) {
		return &RiskCheck{
			Allowed: false,
			Reason:  "최대 포지션 크기 초과",
		}
	}

	// 기존 포지션과의 총 크기 체크
	if existing, exists := m.positions[symbol]; exists {
		totalValue := existing.Quantity.Mul(existing.AvgPrice).Add(orderValue)
		if totalValue.GreaterThan(decimal.NewFromFloat(m.config.Risk.MaxPositionSize)) {
			return &RiskCheck{
				Allowed: false,
				Reason:  "심볼별 최대 포지션 크기 초과",
			}
		}
	}

	return &RiskCheck{Allowed: true}
}

func (m *Manager) UpdatePosition(symbol string, side string, quantity decimal.Decimal, price decimal.Decimal) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 일일 손실 리셋 체크
	if time.Since(m.lastResetDate) > 24*time.Hour {
		m.dailyLoss = decimal.Zero
		m.lastResetDate = time.Now()
		logrus.Info("일일 손실 카운터 리셋")
	}

	if existing, exists := m.positions[symbol]; exists {
		// 기존 포지션 업데이트
		if existing.Side == side {
			// 같은 방향 포지션 추가
			totalQuantity := existing.Quantity.Add(quantity)
			totalValue := existing.Quantity.Mul(existing.AvgPrice).Add(quantity.Mul(price))
			existing.AvgPrice = totalValue.Div(totalQuantity)
			existing.Quantity = totalQuantity
		} else {
			// 반대 방향 포지션 (청산)
			if existing.Quantity.GreaterThanOrEqual(quantity) {
				// 부분 청산
				existing.Quantity = existing.Quantity.Sub(quantity)
				if existing.Quantity.IsZero() {
					delete(m.positions, symbol)
				}
			} else {
				// 전체 청산 후 반대 포지션
				remainingQuantity := quantity.Sub(existing.Quantity)
				delete(m.positions, symbol)
				if !remainingQuantity.IsZero() {
					m.positions[symbol] = &Position{
						Symbol:    symbol,
						Quantity:  remainingQuantity,
						AvgPrice:  price,
						Side:      side,
						Timestamp: time.Now(),
					}
				}
			}
		}
	} else {
		// 새 포지션 생성
		m.positions[symbol] = &Position{
			Symbol:    symbol,
			Quantity:  quantity,
			AvgPrice:  price,
			Side:      side,
			Timestamp: time.Now(),
		}
	}

	logrus.Infof("포지션 업데이트: %s %s %s @ %s", side, quantity.String(), symbol, price.String())
}

func (m *Manager) UpdateDailyLoss(loss decimal.Decimal) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.dailyLoss = m.dailyLoss.Add(loss)
	logrus.Warnf("일일 손실 업데이트: %s (총: %s)", loss.String(), m.dailyLoss.String())
}

func (m *Manager) GetPositions() map[string]*Position {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[string]*Position)
	for k, v := range m.positions {
		result[k] = v
	}
	return result
}

func (m *Manager) GetDailyLoss() decimal.Decimal {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.dailyLoss
}

func (m *Manager) CheckStopLoss(symbol string, currentPrice decimal.Decimal) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if position, exists := m.positions[symbol]; exists {
		stopLossThreshold := position.AvgPrice.Mul(decimal.NewFromFloat(1 - m.config.Risk.StopLossPercentage))

		if position.Side == "long" && currentPrice.LessThan(stopLossThreshold) {
			logrus.Warnf("롱 포지션 스탑로스: %s @ %s (진입가: %s)", symbol, currentPrice.String(), position.AvgPrice.String())
			return true
		}

		if position.Side == "short" && currentPrice.GreaterThan(stopLossThreshold) {
			logrus.Warnf("숏 포지션 스탑로스: %s @ %s (진입가: %s)", symbol, currentPrice.String(), position.AvgPrice.String())
			return true
		}
	}

	return false
}
