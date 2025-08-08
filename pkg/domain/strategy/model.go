package strategy

import (
	"time"
)

// Strategy 전략 인터페이스
type Strategy interface {
	ID() string
	Name() string
	Symbols() []string
	Execute() error
	Start() error
	Stop() error
}

// StrategyDetails 전략 상세 정보
type StrategyDetails struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Active      bool                 `json:"active"`
	Symbols     []string             `json:"symbols"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	Status      string               `json:"status"` // "running", "stopped", "error"
	Performance *StrategyPerformance `json:"performance,omitempty"`
}

// StrategyPerformance 전략 성과 정보
type StrategyPerformance struct {
	StrategyID    string    `json:"strategy_id"`
	TotalReturn   float64   `json:"total_return"`
	WinRate       float64   `json:"win_rate"`
	ProfitLoss    float64   `json:"profit_loss"`
	TradeCount    int64     `json:"trade_count"`
	LastTradeTime time.Time `json:"last_trade_time"`
	MaxDrawdown   float64   `json:"max_drawdown"`
	SharpeRatio   float64   `json:"sharpe_ratio"`
}

// StrategyStatus 전략 실행 상태
type StrategyStatus struct {
	ID             string    `json:"id"`
	Status         string    `json:"status"` // "active", "inactive", "error"
	LastExecution  time.Time `json:"last_execution"`
	ErrorMessage   string    `json:"error_message,omitempty"`
	ExecutionCount int64     `json:"execution_count"`
	Uptime         int64     `json:"uptime"` // 초 단위
}

// StrategyConfig 전략 설정
type StrategyConfig struct {
	ID         string                 `json:"id"`
	Parameters map[string]interface{} `json:"parameters"`
	Enabled    bool                   `json:"enabled"`
	Schedule   string                 `json:"schedule"` // cron 형식
	RiskLimits RiskLimits             `json:"risk_limits"`
}

// RiskLimits 전략별 리스크 제한
type RiskLimits struct {
	MaxPositionSize   float64 `json:"max_position_size"`
	MaxDailyLoss      float64 `json:"max_daily_loss"`
	StopLossPercent   float64 `json:"stop_loss_percent"`
	TakeProfitPercent float64 `json:"take_profit_percent"`
}

// StrategyRequest 전략 요청 DTO
type StrategyRequest struct {
	Action string `json:"action" validate:"required,oneof=start stop restart"`
}

// StrategyResponse 전략 응답 DTO
type StrategyResponse struct {
	Success  bool             `json:"success"`
	Message  string           `json:"message"`
	Strategy *StrategyDetails `json:"strategy,omitempty"`
	Error    string           `json:"error,omitempty"`
}

// StrategyListResponse 전략 목록 응답 DTO
type StrategyListResponse struct {
	Success    bool               `json:"success"`
	Strategies []*StrategyDetails `json:"strategies"`
	Count      int                `json:"count"`
	Active     int                `json:"active_count"`
}

// CreateStrategyRequest 전략 생성 요청
type CreateStrategyRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Type        string                 `json:"type" validate:"required"`
	Description string                 `json:"description"`
	Symbols     []string               `json:"symbols" validate:"required"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// UpdateStrategyRequest 전략 수정 요청
type UpdateStrategyRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Symbols     []string               `json:"symbols,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Enabled     *bool                  `json:"enabled,omitempty"`
}
