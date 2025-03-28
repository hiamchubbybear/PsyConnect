package route

import (
	"consultationservice/handler"
	"github.com/gin-gonic/gin"
)

func RouterInit() *gin.Engine {
	router := gin.Default()
	therapistRoutes := router.Group("/consultation")
	{
		therapistRoutes.GET("/:profile_id", handlers.GetTherapistHandler)
		therapistRoutes.POST("/", handlers.PostTherapistHandler)

	}
	return router
}
