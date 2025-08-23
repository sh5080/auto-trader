package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Portfolio 포트폴리오 응답 데이터
type Portfolio struct {
	ID          string          `json:"id"`
	UserID      string          `json:"user_id"`
	TotalValue  decimal.Decimal `json:"total_value"`
	TotalProfit decimal.Decimal `json:"total_profit"`
	ProfitRate  decimal.Decimal `json:"profit_rate"`
	LastUpdated time.Time       `json:"last_updated"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// PortfolioList 포트폴리오 목록 응답 데이터
type PortfolioList struct {
	Portfolios []*Portfolio `json:"portfolios"`
	Total      int          `json:"total"`
	Limit      int          `json:"limit"`
	Offset     int          `json:"offset"`
}

// Position 포지션 응답 데이터
type Position struct {
	Symbol        string          `json:"symbol"`
	Quantity      decimal.Decimal `json:"quantity"`
	AvgPrice      decimal.Decimal `json:"avg_price"`
	MarketValue   decimal.Decimal `json:"market_value"`
	UnrealizedPnL decimal.Decimal `json:"unrealized_pnl"`
	Return        decimal.Decimal `json:"return"`
	TotalValue    decimal.Decimal `json:"total_value"`
	TotalProfit   decimal.Decimal `json:"total_profit"`
	LastUpdated   time.Time       `json:"last_updated"`
}

// PortfolioSummary 포트폴리오 요약 응답 데이터
type PortfolioSummary struct {
	TotalValue         decimal.Decimal `json:"total_value"`
	TotalCost          decimal.Decimal `json:"total_cost"`
	TotalUnrealizedPnL decimal.Decimal `json:"total_unrealized_pnl"`
	TotalReturn        decimal.Decimal `json:"total_return"`
	PositionCount      int             `json:"position_count"`
	LastUpdated        time.Time       `json:"last_updated"`
	TotalPositions     int             `json:"total_positions"`
}

// CurrentPrice 현재가 응답 데이터
type CurrentPrice struct {
	Symbol    string          `json:"symbol"`
	Price     decimal.Decimal `json:"price"`
	Change    decimal.Decimal `json:"change"`
	ChangePct decimal.Decimal `json:"change_pct"`
	Volume    int64           `json:"volume"`
	Timestamp time.Time       `json:"timestamp"`
}

// StockPrice 주식 가격 응답 데이터
type StockPrice struct {
	Symbol    string          `json:"symbol"`
	Price     decimal.Decimal `json:"price"`
	Timestamp time.Time       `json:"timestamp"`
}

// TradeHistory 거래 내역 응답 데이터
type TradeHistory struct {
	ID        uuid.UUID       `json:"id"`
	Symbol    string          `json:"symbol"`
	Type      string          `json:"type"` // "buy" or "sell"
	Quantity  decimal.Decimal `json:"quantity"`
	Price     decimal.Decimal `json:"price"`
	Amount    decimal.Decimal `json:"amount"`
	Timestamp time.Time       `json:"timestamp"`
}

// TradeHistoryList 거래 내역 목록 응답 데이터
type TradeHistoryList struct {
	Trades []*TradeHistory `json:"trades"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

// CompanyInfo 회사 정보 응답 데이터
type CompanyInfo struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Industry    string `json:"industry"`
	MarketCap   string `json:"market_cap"`
	Employees   int    `json:"employees"`
	Founded     int    `json:"founded"`
	Description string `json:"description"`
}

// ChartData 차트 데이터 응답 데이터
type ChartData struct {
	Symbol    string          `json:"symbol"`
	Timestamp time.Time       `json:"timestamp"`
	Open      decimal.Decimal `json:"open"`
	High      decimal.Decimal `json:"high"`
	Low       decimal.Decimal `json:"low"`
	Close     decimal.Decimal `json:"close"`
	Volume    int64           `json:"volume"`
}

// CreatePortfolioResponse 포트폴리오 생성 응답 데이터
type CreatePortfolioResponse struct {
	Portfolio *Portfolio `json:"portfolio"`
	Message   string     `json:"message"`
}

// UpdatePortfolioResponse 포트폴리오 수정 응답 데이터
type UpdatePortfolioResponse struct {
	Portfolio *Portfolio `json:"portfolio"`
	Message   string     `json:"message"`
}

// DeletePortfolioResponse 포트폴리오 삭제 응답 데이터
type DeletePortfolioResponse struct {
	Message string `json:"message"`
}

// ErrorResponse 에러 응답 데이터
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
