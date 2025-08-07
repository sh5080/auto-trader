package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Response 공통 응답 구조체
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    string      `json:"code,omitempty"`
}

// SuccessResponse 성공 응답
func SuccessResponse(c *fiber.Ctx, data interface{}, statusCode ...int) error {
	code := fiber.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	return c.Status(code).JSON(Response{
		Success: true,
		Data:    data,
	})
}

// ErrorResponse 에러 응답
func ErrorResponse(c *fiber.Ctx, statusCode int, message string, code ...string) error {
	response := Response{
		Success: false,
		Error:   message,
	}

	if len(code) > 0 {
		response.Code = code[0]
	}

	return c.Status(statusCode).JSON(response)
}

// BadRequestResponse 잘못된 요청 응답
func BadRequestResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message, "BAD_REQUEST")
}

// InternalServerErrorResponse 내부 서버 오류 응답
func InternalServerErrorResponse(c *fiber.Ctx, message string, err error) error {
	if err != nil {
		logrus.Errorf("%s: %v", message, err)
	}
	return ErrorResponse(c, fiber.StatusInternalServerError, message, "INTERNAL_SERVER_ERROR")
}

// NotFoundResponse 리소스 없음 응답
func NotFoundResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, message, "NOT_FOUND")
}

// UnauthorizedResponse 인증 실패 응답
func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, message, "UNAUTHORIZED")
}

// ForbiddenResponse 권한 없음 응답
func ForbiddenResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusForbidden, message, "FORBIDDEN")
}

// ValidationErrorResponse 유효성 검사 오류 응답
func ValidationErrorResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message, "VALIDATION_ERROR")
}

// ConflictResponse 충돌 응답
func ConflictResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusConflict, message, "CONFLICT")
}
