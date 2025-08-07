package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// ErrorResponse 에러 응답 구조체
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// SetupAdvancedErrorHandler 고급 에러 핸들링 미들웨어
func SetupAdvancedErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"path":   c.Path(),
				"method": c.Method(),
				"ip":     c.IP(),
				"error":  err.Error(),
			}).Error("요청 처리 중 오류 발생")

			// Fiber 기본 에러 처리
			if e, ok := err.(*fiber.Error); ok {
				return c.Status(e.Code).JSON(ErrorResponse{
					Success: false,
					Error:   e.Message,
					Code:    fmt.Sprintf("HTTP_%d", e.Code),
				})
			}

			// 일반 오류
			return c.Status(500).JSON(ErrorResponse{
				Success: false,
				Error:   "서버 내부 오류가 발생했습니다",
				Code:    "INTERNAL_ERROR",
			})
		}
		return nil
	}
}

// SetupPanicRecovery 패닉 복구 미들웨어
func SetupPanicRecovery() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// 스택 트레이스 가져오기
				stack := debug.Stack()

				logrus.WithFields(logrus.Fields{
					"panic":      r,
					"path":       c.Path(),
					"method":     c.Method(),
					"ip":         c.IP(),
					"user_agent": c.Get("User-Agent"),
					"stack":      string(stack),
				}).Error("패닉 발생 - 서버 복구됨")

				errorResponse := ErrorResponse{
					Success: false,
					Error:   "심각한 서버 오류가 발생했습니다",
					Code:    "PANIC_RECOVERED",
				}

				c.Status(500).JSON(errorResponse)
			}
		}()

		return c.Next()
	}
}

// SetupNotFoundHandler 404 핸들러
func SetupNotFoundHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		logrus.WithFields(logrus.Fields{
			"path":   c.Path(),
			"method": c.Method(),
			"ip":     c.IP(),
		}).Warn("404 - 요청한 리소스를 찾을 수 없습니다")

		return c.Status(404).JSON(ErrorResponse{
			Success: false,
			Error:   "요청한 리소스를 찾을 수 없습니다",
			Code:    "NOT_FOUND",
		})
	}
}
