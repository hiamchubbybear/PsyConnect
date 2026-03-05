package middleware

import (
	"chatservice/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)


func LoggingMiddleware(log *logger.KafkaLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		
		c.Next()

		
		duration := time.Since(start).Milliseconds()

		
		userID, exists := c.Get("userId")
		userIDStr := ""
		if exists && userID != nil {
			userIDStr, _ = userID.(string)
		}

		
		fields := map[string]interface{}{
			"ip":        c.ClientIP(),
			"userAgent": c.Request.UserAgent(),
		}

		if userIDStr != "" {
			fields["userId"] = userIDStr
		}

		
		log.LogHTTPRequest(
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
			fields,
		)
	}
}
