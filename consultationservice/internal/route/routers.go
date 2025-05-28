package route

import (
	"consultationservice/bootstrap"
	handlers "consultationservice/internal/handler"
	"consultationservice/internal/middleware"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func RouterInit(env *bootstrap.Env, clientHandler *handlers.ClientHandler, therapistHandler *handlers.TherapistHandler, matchingHandler *handlers.MatchHandler, sessionHandler *handlers.SessionHandler) {
	urI := fmt.Sprintf("%v:%v", env.Addr, env.Port)
	router := gin.Default()

	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	therapistGroup := router.Group("/consultation/therapist")
	therapistGroup.Use(middleware.RoleRequire("therapist"))
	{
		therapistGroup.GET("/", therapistHandler.GetTherapistHandler)
		therapistGroup.POST("/", therapistHandler.PostTherapistHandler)
		therapistGroup.PUT("/", therapistHandler.PutTherapistHandler)
		therapistGroup.POST("/match/response", matchingHandler.ResponseMatchingRequest)
		therapistGroup.PUT("/status/:status", therapistHandler.ChangeTherapistProfileStatus)
	}

	clientGroup := router.Group("/consultation/client")
	clientGroup.Use(middleware.RoleRequire("client"))
	{
		clientGroup.GET("/", clientHandler.GetClientHandler)
		clientGroup.POST("/", clientHandler.PostClientHandler)
		clientGroup.PUT("/", clientHandler.PutClientHandler)
	}

	publicGroup := router.Group("/consultation/therapist/match")
	{
		publicGroup.GET("/", matchingHandler.GetAllMatchTherapist)
		publicGroup.POST("/", matchingHandler.MatchRequest)
	}

	adminSessionGroup := router.Group("/consultation/admin/session")
	adminSessionGroup.Use(middleware.RoleRequire("admin"))
	{
		adminSessionGroup.GET("/", sessionHandler.GetAllSessions)
	}
	userPublicGroup := router.Group("/consultation/session")
	{
		userPublicGroup.POST("/", sessionHandler.CreateNewSessionHandler)
		userPublicGroup.DELETE("/", sessionHandler.DeleteCurrentSessionHandler)
		userPublicGroup.GET("/", sessionHandler.GetSessionByID)
	}

	router.Run(urI)
}
