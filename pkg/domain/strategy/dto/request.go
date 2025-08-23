package dto

import (
	"github.com/google/uuid"
)

// Body DTOs (JSON 요청 본문)

// CreateStrategyBody 전략 생성 요청 데이터
type CreateStrategyBody struct {
	Name        string    `json:"name" validate:"required,min=1,max=100"`
	Symbol      string    `json:"symbol" validate:"required,min=1,max=20"`
	Description *string   `json:"description,omitempty" validate:"omitempty,max=500"`
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	Active      bool      `json:"active"`
}

// UpdateStrategyBody 전략 수정 요청 데이터
type UpdateStrategyBody struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Symbol      *string `json:"symbol,omitempty" validate:"omitempty,min=1,max=20"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Active      *bool   `json:"active,omitempty"`
}

// Query DTOs (URL 쿼리 파라미터)

// GetStrategyListQuery 전략 목록 조회 쿼리 파라미터
type GetStrategyListQuery struct {
	UserID uuid.UUID `query:"user_id,omitempty"`
	Active *bool     `query:"active,omitempty"`
	Limit  int       `query:"limit" validate:"min=1,max=100"`
	Offset int       `query:"offset" validate:"min=0"`
}

// GetStrategyPerformanceQuery 전략 성과 조회 쿼리 파라미터
type GetStrategyPerformanceQuery struct {
	StrategyID string `query:"strategy_id" validate:"required"`
	StartDate  string `query:"start_date,omitempty"`
	EndDate    string `query:"end_date,omitempty"`
}

// Path DTOs (URL 경로 파라미터)

// StrategyPath 전략 ID 경로 파라미터
type StrategyPath struct {
	ID string `param:"id" validate:"required"`
}

// DeleteStrategyPath 전략 삭제 경로 파라미터
type DeleteStrategyPath struct {
	ID string `param:"id" validate:"required"`
}
