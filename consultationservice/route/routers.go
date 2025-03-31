package route

import (
	"consultationservice/bootstrap"
	"consultationservice/handler"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RouterInit(env *bootstrap.Env) {
	urI := fmt.Sprintf("%v:%v", env.Addr, env.Port)
	router := gin.Default()
	var (
		therapistHandler handlers.TherapistHandler
		clientHandler    handlers.ClientHandler
	)
	therapistRoutes := router.Group("/consultation")
	{
		{
			therapistRoutes.GET("/therapist", therapistHandler.GetTherapistHandler)
			therapistRoutes.POST("/therapist", therapistHandler.PostTherapistHandler)
			therapistRoutes.PUT("/therapist", therapistHandler.PutTherapistHandler)
			therapistRoutes.PUT("/therapist:status", therapistHandler.ChangeTherapistProfileStatus)
		}
		{
			therapistRoutes.GET("/client", clientHandler.GetClientHandler)
			therapistRoutes.POST("/client", clientHandler.PostClientHandler)
			therapistRoutes.PUT("/client", clientHandler.PutClientHandler)
		}
	}
	router.Run(urI)
}
