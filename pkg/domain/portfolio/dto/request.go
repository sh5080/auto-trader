package dto

import (
	"github.com/google/uuid"
)

// Body DTOs (JSON 요청 본문)

// CreatePortfolioBody 포트폴리오 생성 요청 데이터
type CreatePortfolioBody struct {
	Symbol   string    `json:"symbol" validate:"required,min=1,max=20"`
	Quantity float64   `json:"quantity" validate:"required,gt=0"`
	UserID   uuid.UUID `json:"user_id" validate:"required"`
}

// UpdatePortfolioBody 포트폴리오 수정 요청 데이터
type UpdatePortfolioBody struct {
	Symbol   *string  `json:"symbol,omitempty" validate:"omitempty,min=1,max=20"`
	Quantity *float64 `json:"quantity,omitempty" validate:"omitempty,gt=0"`
}

// Query DTOs (URL 쿼리 파라미터)

// GetPortfolioListQuery 포트폴리오 목록 조회 쿼리 파라미터
type GetPortfolioListQuery struct {
	UserID uuid.UUID `query:"user_id,omitempty"`
	Symbol string    `query:"symbol,omitempty"`
	Limit  int       `query:"limit" validate:"min=1,max=100"`
	Offset int       `query:"offset" validate:"min=0"`
}

// GetPortfolioSummaryQuery 포트폴리오 요약 조회 쿼리 파라미터
type GetPortfolioSummaryQuery struct {
	UserID uuid.UUID `query:"user_id,omitempty"`
}

// GetPositionsQuery 보유 주식 목록 조회 쿼리 파라미터
type GetPositionsQuery struct {
	UserID uuid.UUID `query:"user_id,omitempty"`
	Symbol string    `query:"symbol,omitempty"`
}

// GetTradeHistoryQuery 거래 내역 조회 쿼리 파라미터
type GetTradeHistoryQuery struct {
	UserID    uuid.UUID `query:"user_id,omitempty"`
	Symbol    string    `query:"symbol,omitempty"`
	StartDate string    `query:"start_date,omitempty"`
	EndDate   string    `query:"end_date,omitempty"`
	Limit     int       `query:"limit" validate:"min=1,max=100"`
	Offset    int       `query:"offset" validate:"min=0"`
}

// GetCurrentPricesQuery 여러 종목 현재가 조회 쿼리 파라미터
type GetCurrentPricesQuery struct {
	Symbols      string `query:"symbols" validate:"required,min=1"`
	ForceRefresh bool   `query:"forceRefresh"`
}

// Path DTOs (URL 경로 파라미터)

// PortfolioPath 포트폴리오 ID 경로 파라미터
type PortfolioPath struct {
	ID string `param:"id" validate:"required"`
}

// DeletePortfolioPath 포트폴리오 삭제 경로 파라미터
type DeletePortfolioPath struct {
	ID string `param:"id" validate:"required"`
}

// SymbolPath 특정 종목 심볼 경로 파라미터
type SymbolPath struct {
	Symbol string `param:"symbol" validate:"required,min=1"`
}

// GetPortfolioQuery 포트폴리오 조회 쿼리 파라미터
type GetPortfolioQuery struct {
	ForceRefresh bool `query:"forceRefresh"`
}
