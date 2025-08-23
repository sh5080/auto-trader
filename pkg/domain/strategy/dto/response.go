package dto

import (
	"time"

	"github.com/google/uuid"
)

// StrategyResponse 전략 응답 데이터
type StrategyResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Symbol      string    `json:"symbol"`
	Description *string   `json:"description,omitempty"`
	UserID      uuid.UUID `json:"user_id"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StrategyListResponse 전략 목록 응답 데이터
type StrategyListResponse struct {
	Strategies []*StrategyResponse `json:"strategies"`
	Total      int                 `json:"total"`
	Limit      int                 `json:"limit"`
	Offset     int                 `json:"offset"`
}

// StrategyDetails 전략 상세 정보
type StrategyDetails struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Symbol      string    `json:"symbol"`
	Description *string   `json:"description,omitempty"`
	UserID      uuid.UUID `json:"user_id"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StrategyStatus 전략 상태 정보
type StrategyStatus struct {
	ID     uuid.UUID `json:"id"`
	Active bool      `json:"active"`
}

// StrategyPerformance 전략 성과 데이터
type StrategyPerformance struct {
	StrategyID       string  `json:"strategy_id"`
	TotalReturn      float64 `json:"total_return"`
	AnnualReturn     float64 `json:"annual_return"`
	SharpeRatio      float64 `json:"sharpe_ratio"`
	MaxDrawdown      float64 `json:"max_drawdown"`
	WinRate          float64 `json:"win_rate"`
	TotalTrades      int     `json:"total_trades"`
	ProfitableTrades int     `json:"profitable_trades"`
}

// CreateStrategyResponse 전략 생성 응답 데이터
type CreateStrategyResponse struct {
	Strategy *StrategyResponse `json:"strategy"`
	Message  string            `json:"message"`
}

// UpdateStrategyResponse 전략 수정 응답 데이터
type UpdateStrategyResponse struct {
	Strategy *StrategyResponse `json:"strategy"`
	Message  string            `json:"message"`
}

// DeleteStrategyResponse 전략 삭제 응답 데이터
type DeleteStrategyResponse struct {
	Message string `json:"message"`
}

// ErrorResponse 에러 응답 데이터
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
