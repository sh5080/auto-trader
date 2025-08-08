package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
)

// SetupLogger 기본 로깅 미들웨어 설정
func SetupLogger() fiber.Handler {
	var cfg logger.Config
	cfg.Format = "[${time}] ${status} - ${method} ${path} (${latency})\n"
	cfg.TimeFormat = "2006-01-02 15:04:05"
	cfg.TimeZone = "Local"
	return logger.New(cfg)
}

// SetupDetailedLogger 상세 로깅 미들웨어 설정
func SetupDetailedLogger() fiber.Handler {
	var cfg logger.Config
	cfg.Format = "[${time}] ${ip} | ${status} | ${latency} | ${method} | ${path} | ${error}\n"
	cfg.TimeFormat = "2006-01-02 15:04:05"
	cfg.TimeZone = "Local"
	return logger.New(cfg)
}

// SetupJSONLogger JSON 형식 로깅 미들웨어
func SetupJSONLogger() fiber.Handler {
	var cfg logger.Config
	cfg.Format = `{"time":"${time}","ip":"${ip}","method":"${method}","path":"${path}","status":${status},"latency":"${latency}","user_agent":"${ua}","error":"${error}"}` + "\n"
	cfg.TimeFormat = time.RFC3339
	cfg.TimeZone = "UTC"
	return logger.New(cfg)
}

// SetupStructuredLogger Logrus를 활용한 구조화된 로깅
func SetupStructuredLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// 다음 핸들러 실행
		err := c.Next()

		// 로그 기록
		duration := time.Since(start)

		logEntry := logrus.WithFields(logrus.Fields{
			"method":     c.Method(),
			"path":       c.Path(),
			"status":     c.Response().StatusCode(),
			"latency":    duration.String(),
			"ip":         c.IP(),
			"user_agent": c.Get("User-Agent"),
			"request_id": c.Get("X-Request-ID"),
		})

		if err != nil {
			logEntry.WithError(err).Error("HTTP 요청 처리 중 오류 발생")
		} else {
			switch {
			case c.Response().StatusCode() >= 500:
				logEntry.Error("HTTP 요청 완료 (서버 오류)")
			case c.Response().StatusCode() >= 400:
				logEntry.Warn("HTTP 요청 완료 (클라이언트 오류)")
			default:
				logEntry.Info("HTTP 요청 완료")
			}
		}

		return err
	}
}

// SetupCustomLogger 커스텀 로깅 설정
func SetupCustomLogger(format, timeFormat, timeZone string) fiber.Handler {
	var cfg logger.Config
	cfg.Format = format
	cfg.TimeFormat = timeFormat
	cfg.TimeZone = timeZone
	return logger.New(cfg)
}
