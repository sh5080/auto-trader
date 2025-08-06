package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// ErrorResponse 표준 에러 응답 구조체
type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// SetupErrorHandler 기본 에러 핸들링 미들웨어
func SetupErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return c.Status(500).JSON(ErrorResponse{
				Success: false,
				Error:   err.Error(),
			})
		}
		return nil
	}
}

// SetupAdvancedErrorHandler 고급 에러 핸들링 미들웨어
func SetupAdvancedErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			// Fiber 에러 처리
			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				return c.Status(fiberErr.Code).JSON(ErrorResponse{
					Success: false,
					Error:   fiberErr.Message,
					Code:    "FIBER_ERROR",
				})
			}

			// 기본 에러 처리
			logrus.WithError(err).WithFields(logrus.Fields{
				"path":   c.Path(),
				"method": c.Method(),
				"ip":     c.IP(),
			}).Error("HTTP 요청 처리 중 예상치 못한 오류 발생")

			return c.Status(500).JSON(ErrorResponse{
				Success: false,
				Error:   "내부 서버 오류가 발생했습니다",
				Code:    "INTERNAL_SERVER_ERROR",
			})
		}
		return nil
	}
}

// SetupCustomErrorHandler 커스텀 에러 핸들링 미들웨어
func SetupCustomErrorHandler(
	showDetails bool,
	logger *logrus.Logger,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			var fiberErr *fiber.Error
			var statusCode int = 500
			var errorCode string = "INTERNAL_SERVER_ERROR"
			var errorMessage string = "내부 서버 오류가 발생했습니다"

			if errors.As(err, &fiberErr) {
				statusCode = fiberErr.Code
				errorMessage = fiberErr.Message
				errorCode = "FIBER_ERROR"
			} else if showDetails {
				errorMessage = err.Error()
			}

			// 로깅
			logEntry := logger.WithError(err).WithFields(logrus.Fields{
				"path":        c.Path(),
				"method":      c.Method(),
				"ip":          c.IP(),
				"user_agent":  c.Get("User-Agent"),
				"request_id":  c.Get("X-Request-ID"),
				"status_code": statusCode,
			})

			if statusCode >= 500 {
				logEntry.Error("서버 오류 발생")
			} else {
				logEntry.Warn("클라이언트 오류 발생")
			}

			response := ErrorResponse{
				Success: false,
				Error:   errorMessage,
				Code:    errorCode,
			}

			if showDetails {
				response.Details = map[string]interface{}{
					"path":      c.Path(),
					"method":    c.Method(),
					"timestamp": c.Context().Time(),
				}
			}

			return c.Status(statusCode).JSON(response)
		}
		return nil
	}
}

// SetupNotFoundHandler 404 Not Found 핸들러
func SetupNotFoundHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(404).JSON(ErrorResponse{
			Success: false,
			Error:   "요청하신 리소스를 찾을 수 없습니다",
			Code:    "NOT_FOUND",
			Details: map[string]interface{}{
				"path":   c.Path(),
				"method": c.Method(),
			},
		})
	}
}

// SetupPanicRecovery 패닉 복구 미들웨어
func SetupPanicRecovery() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithFields(logrus.Fields{
					"panic":  r,
					"path":   c.Path(),
					"method": c.Method(),
					"ip":     c.IP(),
				}).Error("패닉 발생 - 서버 복구됨")

				c.Status(500).JSON(ErrorResponse{
					Success: false,
					Error:   "심각한 서버 오류가 발생했습니다",
					Code:    "PANIC_RECOVERED",
				})
			}
		}()

		return c.Next()
	}
}
