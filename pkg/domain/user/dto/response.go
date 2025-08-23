package dto

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse 사용자 응답 데이터
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	IsValid   bool      `json:"is_valid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserListResponse 사용자 목록 응답 데이터
type UserListResponse struct {
	Users  []*UserResponse `json:"users"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

// UserWithRelationsResponse 관계 정보를 포함한 사용자 응답 데이터
type UserWithRelationsResponse struct {
	User       *UserResponse        `json:"user"`
	Strategies []*StrategyResponse  `json:"strategies,omitempty"`
	Portfolios []*PortfolioResponse `json:"portfolios,omitempty"`
}

// StrategyResponse 전략 응답 데이터 (간단한 정보)
type StrategyResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Symbol    string    `json:"symbol"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

// PortfolioResponse 포트폴리오 응답 데이터 (간단한 정보)
type PortfolioResponse struct {
	ID        uuid.UUID `json:"id"`
	Symbol    string    `json:"symbol"`
	Quantity  float64   `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

// UserStatsResponse 사용자 통계 응답 데이터
type UserStatsResponse struct {
	TotalUsers    int `json:"total_users"`
	ActiveUsers   int `json:"active_users"`
	InactiveUsers int `json:"inactive_users"`
}

// CreateUserResponse 사용자 생성 응답 데이터
type CreateUserResponse struct {
	User    *UserResponse `json:"user"`
	Message string        `json:"message"`
}

// UpdateUserResponse 사용자 수정 응답 데이터
type UpdateUserResponse struct {
	User    *UserResponse `json:"user"`
	Message string        `json:"message"`
}

// DeleteUserResponse 사용자 삭제 응답 데이터
type DeleteUserResponse struct {
	Message string `json:"message"`
}

// ErrorResponse 에러 응답 데이터
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
