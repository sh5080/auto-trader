package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupCORS CORS 미들웨어 설정
func SetupCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false,
	})
}

// SetupCORSWithConfig 커스텀 CORS 설정
func SetupCORSWithConfig(origins, methods, headers string, credentials bool) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     methods,
		AllowHeaders:     headers,
		AllowCredentials: credentials,
	})
}

// SetupProductionCORS 프로덕션 환경용 CORS 설정
func SetupProductionCORS(allowedOrigins []string) fiber.Handler {
	var origins string
	if len(allowedOrigins) > 0 {
		for i, origin := range allowedOrigins {
			if i > 0 {
				origins += ","
			}
			origins += origin
		}
	} else {
		origins = "*" // 기본값
	}

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		MaxAge:           86400, // 24시간
	})
}
