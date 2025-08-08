package utils

import (
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
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// Error 검증 오류들 메시지 반환
func (e ValidationErrors) Error() string {
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
func (e *ValidationErrors) AddError(field, message string) {
	e.Errors = append(e.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors 오류가 있는지 확인
func (e ValidationErrors) HasErrors() bool {
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

// IsValidationError 오류가 검증 오류인지 확인 (래핑된 오류도 확인)
func IsValidationError(err error) bool {
	if err == nil {
		return false
	}

	// 직접적인 검증 오류 타입 확인
	if _, ok := err.(ValidationError); ok {
		return true
	}
	if _, ok := err.(ValidationErrors); ok {
		return true
	}

	// 래핑된 오류에서 검증 오류 확인
	if wrappedErr, ok := err.(WrappedError); ok {
		return IsValidationError(wrappedErr.OriginalError)
	}

	// fmt.Errorf로 래핑된 오류에서 검증 오류 확인
	if wrappedErr, ok := err.(interface{ Unwrap() error }); ok {
		return IsValidationError(wrappedErr.Unwrap())
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
	if _, ok := err.(ValidationError); ok {
		return err
	}
	if _, ok := err.(ValidationErrors); ok {
		return err
	}

	// 래핑된 오류에서 추출
	if wrappedErr, ok := err.(WrappedError); ok {
		return UnwrapValidationError(wrappedErr.OriginalError)
	}

	// fmt.Errorf로 래핑된 오류에서 추출
	if wrappedErr, ok := err.(interface{ Unwrap() error }); ok {
		return UnwrapValidationError(wrappedErr.Unwrap())
	}

	return err
}
