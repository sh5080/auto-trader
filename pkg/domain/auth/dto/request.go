package dto

// Body DTOs (JSON 요청 본문)

// LoginBody 로그인 요청 데이터
type LoginBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterBody 회원가입 요청 데이터
type RegisterBody struct {
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Nickname string `json:"nickname" validate:"required,min=1,max=50"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

// RefreshTokenBody 토큰 갱신 요청 데이터
type RefreshTokenBody struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LogoutBody 로그아웃 요청 데이터
type LogoutBody struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ChangePasswordBody 비밀번호 변경 요청 데이터
type ChangePasswordBody struct {
	CurrentPassword string `json:"current_password" validate:"required,min=6"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=255"`
}

// ForgotPasswordBody 비밀번호 찾기 요청 데이터
type ForgotPasswordBody struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordBody 비밀번호 재설정 요청 데이터
type ResetPasswordBody struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=255"`
}

// Query DTOs (URL 쿼리 파라미터)

// GetAuthStatusQuery 인증 상태 조회 쿼리 파라미터
type GetAuthStatusQuery struct {
	// 현재는 파라미터가 없지만 향후 확장 가능
}

// Path DTOs (URL 경로 파라미터)

// TokenPath 토큰 경로 파라미터
type TokenPath struct {
	Token string `param:"token" validate:"required"`
}
