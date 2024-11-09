package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RequestLoggingMiddleware logs each incoming request with method, URL, and response status
func RequestLoggingMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	// Process request
	err := c.Next()

	// Log request details after the response is sent
	log.Printf("[%s] %s %s - %d - %v", c.Method(), c.Path(), c.IP(), c.Response().StatusCode(), time.Since(start))

	return err
}
