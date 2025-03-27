package middleware

import (
	"consultationservice/apiresponse"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, apiresponse.ApiResponse{
				Message: "Internal Server Error",
				Status:  http.StatusInternalServerError,
				Data:    nil,
			})
		}
	}
}
