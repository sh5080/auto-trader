package dto

// Body DTOs (JSON 요청 본문)

// CreateUserBody 사용자 생성 요청 데이터
type CreateUserBody struct {
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Nickname string `json:"nickname" validate:"required,min=1,max=50"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

// UpdateUserBody 사용자 수정 요청 데이터
type UpdateUserBody struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Nickname *string `json:"nickname,omitempty" validate:"omitempty,min=1,max=50"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email,max=255"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6,max=255"`
	IsValid  *bool   `json:"is_valid,omitempty"`
}

// Query DTOs (URL 쿼리 파라미터)

// GetUserListQuery 사용자 목록 조회 쿼리 파라미터
type GetUserListQuery struct {
	Limit  int `query:"limit" validate:"min=1,max=100"`
	Offset int `query:"offset" validate:"min=0"`
}

// GetUserStatsQuery 사용자 통계 조회 쿼리 파라미터
type GetUserStatsQuery struct {
	// 현재는 파라미터가 없지만 향후 확장 가능
}

// Path DTOs (URL 경로 파라미터)

// UserPath 사용자 ID 경로 파라미터
type UserPath struct {
	ID string `param:"id" validate:"required"`
}

// DeleteUserPath 사용자 삭제 경로 파라미터
type DeleteUserPath struct {
	ID string `param:"id" validate:"required"`
}
