package middleware

import (
	"consultationservice/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(log *logger.KafkaLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Milliseconds()

		// Get user ID from context (if available)
		userID, exists := c.Get("userId")
		userIDStr := ""
		if exists && userID != nil {
			userIDStr, _ = userID.(string)
		}

		// Prepare fields
		fields := map[string]interface{}{
			"ip":        c.ClientIP(),
			"userAgent": c.Request.UserAgent(),
		}

		if userIDStr != "" {
			fields["userId"] = userIDStr
		}

		// Log request
		log.LogHTTPRequest(
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
			fields,
		)
	}
}
