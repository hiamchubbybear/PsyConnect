package route

import (
	"consultationservice/handler"
	"github.com/gin-gonic/gin"
)

func RouterInit() *gin.Engine {
	router := gin.Default()
	therapistRoutes := router.Group("/consultation")
	{
		therapistRoutes.GET("/", handlers.GetTherapist)
		therapistRoutes.POST("/", handlers.PostTherapist)
	}

	return router
}
