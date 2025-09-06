package route

import (
	"consultationservice/bootstrap"
	handlers "consultationservice/internal/handler"
	"consultationservice/internal/middleware"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func RouterInit(
	env *bootstrap.Env,
	clientHandler *handlers.ClientHandler,
	therapistHandler *handlers.TherapistHandler,
	matchingHandler *handlers.MatchHandler,
	sessionHandler *handlers.SessionHandler,
	swipeHandler *handlers.SwipeHandler,
) {
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
		therapistGroup.GET("/", therapistHandler.GetTherapistHandler)                        // Deprecated
		therapistGroup.POST("/", therapistHandler.PostTherapistHandler)                      // Deprecated
		therapistGroup.PUT("/", therapistHandler.PutTherapistHandler)                        // Deprecated
		therapistGroup.PUT("/status/:status", therapistHandler.ChangeTherapistProfileStatus) // Deprecated
	}

	clientGroup := router.Group("/consultation/client")
	clientGroup.Use(middleware.RoleRequire("client"))
	{
		clientGroup.GET("/", clientHandler.GetClientHandler)   // Deprecated
		clientGroup.POST("/", clientHandler.PostClientHandler) // Deprecated
		clientGroup.PUT("/", clientHandler.PutClientHandler)   // Deprecated

		clientGroup.POST("/recommend", swipeHandler.TriggerUpdate) // Deprecated
		clientGroup.GET("/recommend/top", swipeHandler.PopTop5)    // Deprecated
		clientGroup.POST("/match", matchingHandler.MatchRequest)   // Deprecated
		clientGroup.POST("/swipe", swipeHandler.SwipeTherapist)    // Deprecated
	}

	publicGroup := router.Group("/consultation/therapist/match")
	{
		publicGroup.GET("/", matchingHandler.GetAllMatchTherapist) // Deprecated
		publicGroup.POST("/", matchingHandler.MatchRequest)        // Deprecated
	}

	adminSessionGroup := router.Group("/consultation/admin/session")
	adminSessionGroup.Use(middleware.RoleRequire("admin"))
	{
		adminSessionGroup.GET("/", sessionHandler.GetAllSessions) // Deprecated
	}

	userSessionGroup := router.Group("/consultation/session")
	userSessionGroup.Use(middleware.RoleRequire(""))
	{
		userSessionGroup.GET("/all", sessionHandler.GetAllSessionByID)           // Deprecated
		userSessionGroup.POST("/", sessionHandler.CreateNewSessionHandler)       // Deprecated
		userSessionGroup.DELETE("/", sessionHandler.DeleteCurrentSessionHandler) // Deprecated
		userSessionGroup.GET("/:id", sessionHandler.GetSessionByID)              // Deprecated
	}

	uncategoryGroup := router.Group("/consultation")
	uncategoryGroup.Use(middleware.RoleRequire(""))
	{
		uncategoryGroup.GET("/therapist/:id", therapistHandler.GetTherapistByIdHandlerV1) // Deprecated
		uncategoryGroup.GET("/client/:id", clientHandler.GetClientByIdHandler)            // Deprecated
	}

	api := router.Group("/v1/consultation")

	// Therapists
	therapist := api.Group("/therapists")
	therapist.Use(middleware.RoleRequire("therapist"))
	{
		therapist.GET("/me", therapistHandler.GetTherapistHandlerV1)
		therapist.POST("/me", therapistHandler.PostTherapistHandlerV1)
		therapist.PUT("/me", therapistHandler.PutTherapistHandlerV1)
		therapist.PUT("/me/status/:status", therapistHandler.ChangeTherapistProfileStatus)

		// admin / public
		therapist.GET("/:id", therapistHandler.GetTherapistByIdHandlerV1)
	}

	// Clients
	client := api.Group("/clients")
	client.Use(middleware.RoleRequire("client"))
	{
		// thao tác với chính mình
		client.GET("/me", clientHandler.GetClientHandler)
		client.POST("/me", clientHandler.PostClientHandler)
		client.PUT("/me", clientHandler.PutClientHandler)

		client.POST("/me/recommend", swipeHandler.TriggerUpdate)
		client.GET("/me/recommend/top", swipeHandler.PopTop5)
		client.POST("/me/match", matchingHandler.MatchRequest)
		client.POST("/me/swipe", swipeHandler.SwipeTherapist)

		// admin / public
		client.GET("/:id", clientHandler.GetClientByIdHandler)
	}

	// Sessions
	session := api.Group("/sessions")
	{
		session.GET("/me", sessionHandler.GetAllSessionByID)
		session.POST("/me", sessionHandler.CreateNewSessionHandler)
		session.DELETE("/:id", sessionHandler.DeleteCurrentSessionHandler)
		session.GET("/:id", sessionHandler.GetSessionByID) // admin / public
	}

	// Admin sessions
	adminSession := api.Group("/admin/sessions")
	adminSession.Use(middleware.RoleRequire("admin"))
	{
		adminSession.GET("/", sessionHandler.GetAllSessions)
	}

	// Public match
	match := api.Group("/matches")
	{
		match.GET("/", matchingHandler.GetAllMatchTherapist)
		match.POST("/", matchingHandler.MatchRequest)
	}
	// Uncategorize Routes
	uncategoryGroupV1 := router.Group("/v1/consultation")
	uncategoryGroup.Use(middleware.RoleRequire(""))
	{

		uncategoryGroupV1.GET("/therapist/:id", therapistHandler.GetTherapistByIdHandlerV1)
		uncategoryGroupV1.GET("/me/recommend/top", swipeHandler.PopTop5V1)
		uncategoryGroupV1.GET("/client/:id", clientHandler.GetClientByIdHandler)
	}

	router.Run(urI)
}
