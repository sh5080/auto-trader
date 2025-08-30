package types

import "errors"

// 커스텀 에러 타입
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)
