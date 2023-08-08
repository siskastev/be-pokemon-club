package middleware

import (
	"be-pokemon-club/internal/helpers/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Get the logger instance from the logger package
		log := logger.Logger()
		c.Locals("logger", log) // Store the logger in Fiber context locals

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		// Log request information using the logger instance
		log.WithFields(logrus.Fields{
			"method":  c.Method(),
			"route":   c.Path(),
			"status":  c.Response().StatusCode(),
			"latency": latency,
		}).Info("Request processed")

		return nil
	}
}
