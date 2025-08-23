package dto

import (
	"time"

	"github.com/google/uuid"
)

// LoginResponse 로그인 응답 데이터
type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	TokenType    string   `json:"token_type"`
	ExpiresIn    int      `json:"expires_in"`
	User         UserInfo `json:"user"`
	Message      string   `json:"message"`
}

// RegisterResponse 회원가입 응답 데이터
type RegisterResponse struct {
	User    UserInfo `json:"user"`
	Message string   `json:"message"`
}

// UserInfo 사용자 정보 (인증 응답용)
type UserInfo struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	IsValid   bool      `json:"is_valid"`
	CreatedAt time.Time `json:"created_at"`
}

// TokenResponse 토큰 응답 데이터
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// RefreshTokenResponse 토큰 갱신 응답 데이터
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Message     string `json:"message"`
}

// LogoutResponse 로그아웃 응답 데이터
type LogoutResponse struct {
	Message string `json:"message"`
}

// ChangePasswordResponse 비밀번호 변경 응답 데이터
type ChangePasswordResponse struct {
	Message string `json:"message"`
}

// ForgotPasswordResponse 비밀번호 찾기 응답 데이터
type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

// ResetPasswordResponse 비밀번호 재설정 응답 데이터
type ResetPasswordResponse struct {
	Message string `json:"message"`
}

// ErrorResponse 에러 응답 데이터
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
