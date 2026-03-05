package route

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"consultationservice/bootstrap"
	clientHTTP "consultationservice/internal/client/transport/http"
	httpHandler "consultationservice/internal/consultation/transport/http"
	matchHTTP "consultationservice/internal/matching/transport/http"
	"consultationservice/internal/middleware"
	commentHTTP "consultationservice/internal/newsfeed/comment/transport/http"
	groupHTTP "consultationservice/internal/newsfeed/group/transport/http"
	postHTTP "consultationservice/internal/newsfeed/post/transport/http"
	reactionHTTP "consultationservice/internal/newsfeed/reaction/transport/http"
	socialHTTP "consultationservice/internal/newsfeed/social/transport/http"
	swipeHTTP "consultationservice/internal/swipe/transport/http"
	therapistHTTP "consultationservice/internal/therapist/transport/http"
	"consultationservice/pkg/logger"
)

func RouterInit(
	env *bootstrap.Env,
	kafkaLogger *logger.KafkaLogger,

	clientHandler *clientHTTP.Handler,
	therapistHandler *therapistHTTP.Handler,
	matchingHandler *matchHTTP.Handler,
	sessionHandler *httpHandler.Handler,
	swipeHandler *swipeHTTP.Handler,

	postHandler *postHTTP.Handler,
	reactionHandler *reactionHTTP.Handler,
	commentHandler *commentHTTP.Handler,
	socialHandler *socialHTTP.Handler,
	groupHandler *groupHTTP.Handler,
) {
	urI := fmt.Sprintf("%v:%v", env.Addr, env.Port)
	router := gin.Default()

	router.GET("/actuator/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	router.Use(middleware.LoggingMiddleware(kafkaLogger))

	defer func() {
		if err := recover(); err != nil {
			kafkaLogger.Fatal("Router panic", map[string]interface{}{
				"error": fmt.Sprintf("%v", err),
			})

		}
	}()

	therapistGroup := router.Group("/consultation/therapist")
	therapistGroup.Use(middleware.RoleRequire("therapist"))
	{
		therapistGroup.GET("/", therapistHandler.GetTherapist)
		therapistGroup.POST("/", therapistHandler.CreateTherapist)
		therapistGroup.PUT("/", therapistHandler.UpdateTherapist)
		therapistGroup.PATCH("/status", therapistHandler.UpdateAvailability)
	}

	clientGroup := router.Group("/consultation/client")
	clientGroup.Use(middleware.RoleRequire("client"))
	{
		clientGroup.GET("/", clientHandler.GetClient)
		clientGroup.POST("/", clientHandler.CreateClient)
		clientGroup.PUT("/", clientHandler.UpdateClient)

		clientGroup.POST("/recommend", swipeHandler.TriggerUpdate)
		clientGroup.GET("/recommend/top", swipeHandler.PopTop5)
		clientGroup.POST("/match", matchingHandler.MatchRequest)
		clientGroup.POST("/swipe", swipeHandler.SwipeTherapist)
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

	userSessionGroup := router.Group("/consultation/session")
	userSessionGroup.Use(middleware.RoleRequire(""))
	{
		userSessionGroup.GET("/all", sessionHandler.GetSessionsByProfile)
		userSessionGroup.POST("/", sessionHandler.CreateSession)
		userSessionGroup.DELETE("/", sessionHandler.DeleteSession)
		userSessionGroup.GET("/:id", sessionHandler.GetSession)
	}

	uncategoryGroup := router.Group("/consultation")
	uncategoryGroup.Use(middleware.RoleRequire(""))
	{
		uncategoryGroup.GET("/therapist/:id", therapistHandler.GetTherapistByID)
		uncategoryGroup.GET("/client/:id", clientHandler.GetClientByID)
	}

	api := router.Group("/v1/consultation")

	therapist := api.Group("/therapists")
	therapist.Use(middleware.RoleRequire("therapist"))
	{
		therapist.GET("/me", therapistHandler.GetTherapist)
		therapist.POST("/me", therapistHandler.CreateTherapist)
		therapist.PUT("/me", therapistHandler.UpdateTherapist)
		therapist.PATCH("/me/availability", therapistHandler.UpdateAvailability)

		therapist.GET("/:id", therapistHandler.GetTherapistByID)
	}

	client := api.Group("/clients")
	client.Use(middleware.RoleRequire("client"))
	{
		client.GET("/me", clientHandler.GetClient)
		client.POST("/me", clientHandler.CreateClient)
		client.PUT("/me", clientHandler.UpdateClient)

		client.POST("/me/recommend", swipeHandler.TriggerUpdateV1)
		client.GET("/me/recommend/top", swipeHandler.PopTop5V1)
		client.POST("/me/match", matchingHandler.MatchRequest)
		client.POST("/me/swipe", swipeHandler.SwipeTherapist)

		client.GET("/:id", clientHandler.GetClientByID)
	}

	session := api.Group("/sessions")
	{
		session.GET("/me", sessionHandler.GetSessionsByProfile)
		session.GET("/me/calendar", sessionHandler.GetCalendar)
		session.GET("/me/overview", sessionHandler.GetOverview)
		session.POST("/me", sessionHandler.CreateSession)
		session.DELETE("/:id", sessionHandler.DeleteSession)
		session.GET("/:id", sessionHandler.GetSession)
		session.POST("/:id/call/start", sessionHandler.StartCall)

		// Payment & Refund Routes
		session.GET("/:id/payment-url", sessionHandler.GetPaymentURL)          // Mock Payment URL
		session.POST("/webhook/payment", sessionHandler.ProcessPaymentWebhook) // Payment Webhook
		session.POST("/:id/refund", sessionHandler.RefundSession)              // Process Refund
	}

	adminSession := api.Group("/admin/sessions")
	adminSession.Use(middleware.RoleRequire("admin"))
	{
		adminSession.GET("/", sessionHandler.GetAllSessions)
	}

	match := api.Group("/matches")
	{
		match.GET("/", matchingHandler.GetAllMatchTherapist)
		match.POST("/", matchingHandler.MatchRequest)
	}

	uncategoryGroupV1 := router.Group("/v1/consultation")
	uncategoryGroupV1.Use(middleware.RoleRequire(""))
	{
		uncategoryGroupV1.GET("/clients", clientHandler.GetAllClients)
		uncategoryGroupV1.GET("/therapist/:id", therapistHandler.GetTherapistByID)
		uncategoryGroupV1.GET("/search/therapists", therapistHandler.SearchTherapists)
		uncategoryGroupV1.GET("/me/recommend/top", swipeHandler.PopTop5V1)
		uncategoryGroupV1.GET("/client/:id", clientHandler.GetClientByID)
	}

	posts := api.Group("/posts")
	posts.Use(middleware.RoleRequire(""))
	{
		posts.POST("/", postHandler.CreatePost)
		posts.GET("/", postHandler.GetFeed)
		posts.GET("/trending", postHandler.GetTrendingPosts)
		posts.GET("/search", postHandler.SearchPosts)
		posts.GET("/tags/popular", postHandler.GetPopularTags)
		posts.GET("/user/:userId", postHandler.GetUserPosts)
		posts.GET("/:id", postHandler.GetPostByID)
		posts.PUT("/:id", postHandler.UpdatePost)
		posts.DELETE("/:id", postHandler.DeletePost)
		posts.POST("/:id/view", postHandler.IncrementViewCount)

		posts.POST("/:id/react", reactionHandler.AddReaction)
		posts.DELETE("/:id/react", reactionHandler.RemoveReaction)
		posts.GET("/:id/reactions", reactionHandler.GetPostReactions)

		posts.POST("/:id/comments", commentHandler.CreateComment)
		posts.GET("/:id/comments", commentHandler.GetComments)

		posts.POST("/:id/share", socialHandler.SharePost)

		posts.POST("/:id/bookmark", socialHandler.AddBookmark)
		posts.DELETE("/:id/bookmark", socialHandler.RemoveBookmark)
	}

	comments := api.Group("/comments")
	comments.Use(middleware.RoleRequire(""))
	{
		comments.PUT("/:id", commentHandler.UpdateComment)
		comments.DELETE("/:id", commentHandler.DeleteComment)
		comments.GET("/:id/replies", commentHandler.GetReplies)
	}

	users := api.Group("/users")
	users.Use(middleware.RoleRequire(""))
	{

		users.POST("/:id/follow", socialHandler.FollowUser)
		users.DELETE("/:id/follow", socialHandler.UnfollowUser)
		users.GET("/:id/followers", socialHandler.GetFollowers)
		users.GET("/:id/following", socialHandler.GetFollowing)
	}

	bookmarks := api.Group("/bookmarks")
	bookmarks.Use(middleware.RoleRequire(""))
	{
		bookmarks.GET("/", socialHandler.GetBookmarks)
	}

	groups := api.Group("/groups")
	groups.Use(middleware.RoleRequire(""))
	{
		groups.POST("/", groupHandler.CreateGroup)
		groups.GET("/", groupHandler.GetGroups)
		groups.GET("/:id", groupHandler.GetGroupByID)
		groups.POST("/:id/join", groupHandler.JoinGroup)
	}

	tags := api.Group("/tags")
	{
		tags.GET("/:tag/posts", postHandler.GetPostsByTag)
	}

	categories := api.Group("/categories")
	{
		categories.GET("/:category/posts", postHandler.GetPostsByCategory)
	}

	router.Run(urI)
}
