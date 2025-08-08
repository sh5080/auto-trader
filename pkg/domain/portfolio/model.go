package portfolio

import (
	"time"

	"github.com/shopspring/decimal"
)

// Portfolio 포트폴리오 정보
type Portfolio struct {
	ID          string          `json:"id" db:"id"`
	UserID      string          `json:"user_id" db:"user_id"`
	TotalValue  decimal.Decimal `json:"total_value" db:"total_value"`
	TotalProfit decimal.Decimal `json:"total_profit" db:"total_profit"`
	ProfitRate  decimal.Decimal `json:"profit_rate" db:"profit_rate"`
	LastUpdated time.Time       `json:"last_updated" db:"last_updated"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

// Position 보유 주식 정보
type Position struct {
	ID              string          `json:"id" db:"id"`
	UserID          string          `json:"user_id" db:"user_id"`
	Symbol          string          `json:"symbol" db:"symbol"`
	CompanyName     string          `json:"company_name" db:"company_name"`
	Quantity        decimal.Decimal `json:"quantity" db:"quantity"`
	AveragePrice    decimal.Decimal `json:"average_price" db:"average_price"`
	CurrentPrice    decimal.Decimal `json:"current_price" db:"current_price"`
	TotalValue      decimal.Decimal `json:"total_value" db:"total_value"`
	TotalProfit     decimal.Decimal `json:"total_profit" db:"total_profit"`
	ProfitRate      decimal.Decimal `json:"profit_rate" db:"profit_rate"`
	DailyProfit     decimal.Decimal `json:"daily_profit" db:"daily_profit"`
	DailyProfitRate decimal.Decimal `json:"daily_profit_rate" db:"daily_profit_rate"`
	LastUpdated     time.Time       `json:"last_updated" db:"last_updated"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

// StockPrice 주식 가격 정보
type StockPrice struct {
	Symbol        string          `json:"symbol" db:"symbol"`
	Price         decimal.Decimal `json:"price" db:"price"`
	Change        decimal.Decimal `json:"change" db:"change"`
	ChangeRate    decimal.Decimal `json:"change_rate" db:"change_rate"`
	Volume        int64           `json:"volume" db:"volume"`
	MarketCap     decimal.Decimal `json:"market_cap" db:"market_cap"`
	High          decimal.Decimal `json:"high" db:"high"`
	Low           decimal.Decimal `json:"low" db:"low"`
	Open          decimal.Decimal `json:"open" db:"open"`
	PreviousClose decimal.Decimal `json:"previous_close" db:"previous_close"`
	Timestamp     time.Time       `json:"timestamp" db:"timestamp"`
}

// TradeHistory 거래 내역
type TradeHistory struct {
	ID        string          `json:"id" db:"id"`
	UserID    string          `json:"user_id" db:"user_id"`
	Symbol    string          `json:"symbol" db:"symbol"`
	Type      string          `json:"type" db:"type"` // BUY, SELL
	Quantity  decimal.Decimal `json:"quantity" db:"quantity"`
	Price     decimal.Decimal `json:"price" db:"price"`
	Total     decimal.Decimal `json:"total" db:"total"`
	Fee       decimal.Decimal `json:"fee" db:"fee"`
	Timestamp time.Time       `json:"timestamp" db:"timestamp"`
}

// CompanyInfo 회사 정보
type CompanyInfo struct {
	Symbol      string          `json:"symbol" db:"symbol"`
	Name        string          `json:"name" db:"name"`
	Sector      string          `json:"sector" db:"sector"`
	Industry    string          `json:"industry" db:"industry"`
	Country     string          `json:"country" db:"country"`
	MarketCap   decimal.Decimal `json:"market_cap" db:"market_cap"`
	Employees   int             `json:"employees" db:"employees"`
	Website     string          `json:"website" db:"website"`
	Description string          `json:"description" db:"description"`
}

// ChartData 차트 데이터
type ChartData struct {
	Symbol    string          `json:"symbol" db:"symbol"`
	Timestamp time.Time       `json:"timestamp" db:"timestamp"`
	Open      decimal.Decimal `json:"open" db:"open"`
	High      decimal.Decimal `json:"high" db:"high"`
	Low       decimal.Decimal `json:"low" db:"low"`
	Close     decimal.Decimal `json:"close" db:"close"`
	Volume    int64           `json:"volume" db:"volume"`
}

// PortfolioSummary 포트폴리오 요약
type PortfolioSummary struct {
	TotalPositions  int             `json:"total_positions"`
	TotalValue      decimal.Decimal `json:"total_value"`
	TotalProfit     decimal.Decimal `json:"total_profit"`
	TotalProfitRate decimal.Decimal `json:"total_profit_rate"`
	DailyProfit     decimal.Decimal `json:"daily_profit"`
	DailyProfitRate decimal.Decimal `json:"daily_profit_rate"`
	TopGainers      []Position      `json:"top_gainers"`
	TopLosers       []Position      `json:"top_losers"`
	Positions       []Position      `json:"positions"`
	LastUpdated     time.Time       `json:"last_updated"`
	DataFreshness   string          `json:"data_freshness"` // REALTIME, CACHED, STALE
}

// Request/Response DTOs

// GetPortfolioRequest 포트폴리오 조회 요청
type GetPortfolioRequest struct {
	UserID       string `json:"user_id" validate:"required"`
	ForceRefresh bool   `json:"force_refresh"` // 강제 새로고침 여부
}

// GetPositionRequest 보유 주식 조회 요청
type GetPositionRequest struct {
	UserID       string `json:"user_id" validate:"required"`
	Symbol       string `json:"symbol,omitempty"` // 특정 주식만 조회할 때
	ForceRefresh bool   `json:"force_refresh"`
}

// GetStockPriceRequest 주식 가격 조회 요청
type GetStockPriceRequest struct {
	Symbols      []string `json:"symbols" validate:"required"`
	ForceRefresh bool     `json:"force_refresh"`
}

// GetTradeHistoryRequest 거래 내역 조회 요청
type GetTradeHistoryRequest struct {
	UserID    string    `json:"user_id" validate:"required"`
	Symbol    string    `json:"symbol,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	Offset    int       `json:"offset,omitempty"`
}

// CacheConfig 캐시 설정
type CacheConfig struct {
	PortfolioCacheTTL time.Duration `json:"portfolio_cache_ttl"` // 포트폴리오 캐시 TTL
	PriceCacheTTL     time.Duration `json:"price_cache_ttl"`     // 가격 캐시 TTL
	PositionCacheTTL  time.Duration `json:"position_cache_ttl"`  // 포지션 캐시 TTL
}
