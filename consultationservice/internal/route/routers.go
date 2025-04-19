package route

import (
	"consultationservice/bootstrap"
	"consultationservice/internal/handler"
	"consultationservice/internal/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func RouterInit(env *bootstrap.Env) {
	urI := fmt.Sprintf("%v:%v", env.Addr, env.Port)
	router := gin.Default()

	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	var (
		therapistHandler handlers.TherapistHandler
		clientHandler    handlers.ClientHandler
	)
	therapistGroup := router.Group("/consultation/therapist")
	therapistGroup.Use(middleware.RoleRequire("therapist"))
	{
		therapistGroup.GET("/", therapistHandler.GetTherapistHandler)
		therapistGroup.POST("/", therapistHandler.PostTherapistHandler)
		therapistGroup.PUT("/", therapistHandler.PutTherapistHandler)
		therapistGroup.PUT("/status/:status", therapistHandler.ChangeTherapistProfileStatus)
	}

	clientGroup := router.Group("/consultation/client")
	clientGroup.Use(middleware.RoleRequire("client"))
	{
		clientGroup.GET("/", clientHandler.GetClientHandler)
		clientGroup.POST("/", clientHandler.PostClientHandler)
		clientGroup.PUT("/", clientHandler.PutClientHandler)
	}

	router.Run(urI)
}
