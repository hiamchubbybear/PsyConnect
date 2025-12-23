package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(api *API) *gin.Engine {
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200", "https://chessy.dev", "http://localhost:8100"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Mail routes
	mail := router.Group("/mail")
	{
		mail.POST("/activate", api.SendActivateEmail)
		mail.POST("/account-update", api.SendAccountUpdateEmail)
	}

	// Notification routes
	notification := router.Group("/notification")
	{
		notification.POST("/token", api.SaveFCMToken)
		notification.GET("/:userId", api.GetNotifications)
		notification.PUT("/:id/read", api.MarkAsRead)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "notification-service",
		})
	})

	router.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ready",
			"service": "notification-service",
		})
	})

	return router
}
