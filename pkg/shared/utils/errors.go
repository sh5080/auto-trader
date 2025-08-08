package utils

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError 검증 오류 구조체
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error 검증 오류 메시지 반환
func (e ValidationError) Error() string {
	return fmt.Sprintf("필드 '%s': %s", e.Field, e.Message)
}

// ValidationErrors 여러 검증 오류들을 담는 구조체
type ValidationsError struct {
	Errors []ValidationError `json:"errors"`
}

// Error 검증 오류들 메시지 반환
func (e ValidationsError) Error() string {
	if len(e.Errors) == 0 {
		return "검증 오류가 없습니다"
	}

	var messages []string
	for _, err := range e.Errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// AddError 검증 오류 추가
func (e *ValidationsError) AddError(field, message string) {
	e.Errors = append(e.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors 오류가 있는지 확인
func (e ValidationsError) HasErrors() bool {
	return len(e.Errors) > 0
}

// WrappedError 래핑된 오류 구조체 (검증 오류 보존용)
type WrappedError struct {
	OriginalError error
	Message       string
}

// Error 래핑된 오류 메시지 반환
func (e WrappedError) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.OriginalError)
}

// Unwrap 원본 오류 반환
func (e WrappedError) Unwrap() error {
	return e.OriginalError
}

// ==================== AppError (의도 기반 에러) ====================

// ErrorKind 에러 의도(HTTP 의미와 매핑 가능한 종류)
type ErrorKind string

const (
	ErrorBadRequest   ErrorKind = "bad_request"
	ErrorUnauthorized ErrorKind = "unauthorized"
	ErrorForbidden    ErrorKind = "forbidden"
	ErrorNotFound     ErrorKind = "not_found"
	ErrorConflict     ErrorKind = "conflict"
	ErrorInternal     ErrorKind = "internal"
)

// AppError 애플리케이션 의도 에러
type AppError struct {
	Kind    ErrorKind
	Field   string
	Message string
	Err     error
}

// Error 에러 문자열
func (e AppError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

// Unwrap 원인 에러 제공
func (e AppError) Unwrap() error { return e.Err }

// 생성자들
func BadRequest(message string) error {
	return AppError{Kind: ErrorBadRequest, Field: "", Message: message, Err: nil}
}
func Unauthorized(message string) error {
	return AppError{Kind: ErrorUnauthorized, Field: "", Message: message, Err: nil}
}
func Forbidden(message string) error {
	return AppError{Kind: ErrorForbidden, Field: "", Message: message, Err: nil}
}
func NotFound(field, message string) error {
	return AppError{Kind: ErrorNotFound, Field: field, Message: message, Err: nil}
}
func Conflict(field, message string) error {
	return AppError{Kind: ErrorConflict, Field: field, Message: message, Err: nil}
}
func Internal(message string, cause error) error {
	return AppError{Kind: ErrorInternal, Field: "", Message: message, Err: cause}
}

// 판별자들
func IsBadRequest(err error) bool   { return isKind(err, ErrorBadRequest) }
func IsUnauthorized(err error) bool { return isKind(err, ErrorUnauthorized) }
func IsForbidden(err error) bool    { return isKind(err, ErrorForbidden) }
func IsNotFound(err error) bool     { return isKind(err, ErrorNotFound) }
func IsConflict(err error) bool     { return isKind(err, ErrorConflict) }
func IsInternal(err error) bool     { return isKind(err, ErrorInternal) }

func isKind(err error, k ErrorKind) bool {
	if err == nil {
		return false
	}
	var ae AppError
	if errors.As(err, &ae) {
		return ae.Kind == k
	}
	if u, ok := err.(interface{ Unwrap() error }); ok {
		return isKind(u.Unwrap(), k)
	}
	return false
}

// IsValidationError 오류가 검증 오류인지 확인 (래핑된 오류도 확인)
func IsValidationError(err error) bool {
	if err == nil {
		return false
	}
	var ve ValidationError
	if errors.As(err, &ve) {
		return true
	}
	var vel ValidationsError
	if errors.As(err, &vel) {
		return true
	}
	var wrapped WrappedError
	if errors.As(err, &wrapped) {
		return IsValidationError(wrapped.OriginalError)
	}
	if u, ok := err.(interface{ Unwrap() error }); ok {
		return IsValidationError(u.Unwrap())
	}
	return false
}

// WrapValidationError 검증 오류를 래핑 (타입 보존)
func WrapValidationError(err error, message string) error {
	return WrappedError{
		OriginalError: err,
		Message:       message,
	}
}

// UnwrapValidationError 래핑된 오류에서 검증 오류 추출
func UnwrapValidationError(err error) error {
	if err == nil {
		return nil
	}

	// 직접적인 검증 오류
	var ve ValidationError
	if errors.As(err, &ve) {
		return ve
	}
	var ves ValidationsError
	if errors.As(err, &ves) {
		return ves
	}

	// 래핑된 오류에서 추출
	var wrappedErr WrappedError
	if errors.As(err, &wrappedErr) {
		return UnwrapValidationError(wrappedErr.OriginalError)
	}

	// fmt.Errorf로 래핑된 오류에서 추출
	if wrappedErr, ok := err.(interface{ Unwrap() error }); ok {
		return UnwrapValidationError(wrappedErr.Unwrap())
	}

	return err
}
