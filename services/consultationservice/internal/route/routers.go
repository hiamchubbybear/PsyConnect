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

	// DDD Handlers
	clientHandler *clientHTTP.Handler,
	therapistHandler *therapistHTTP.Handler,
	matchingHandler *matchHTTP.Handler,
	sessionHandler *httpHandler.Handler,
	swipeHandler *swipeHTTP.Handler,

	// Newsfeed handlers
	postHandler *postHTTP.Handler, // DDD Handler
	reactionHandler *reactionHTTP.Handler, // DDD Handler
	commentHandler *commentHTTP.Handler, // DDD Handler
	socialHandler *socialHTTP.Handler, // DDD Handler
) {
	urI := fmt.Sprintf("%v:%v", env.Addr, env.Port)
	router := gin.Default()

	// Healthcheck endpoint
	router.GET("/actuator/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// Add logging middleware
	router.Use(middleware.LoggingMiddleware(kafkaLogger))

	defer func() {
		if err := recover(); err != nil {
			kafkaLogger.Fatal("Router panic", map[string]interface{}{
				"error": fmt.Sprintf("%v", err),
			})
			// log.Fatal(err)
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
		clientGroup.GET("/", clientHandler.GetClient)     // Deprecated
		clientGroup.POST("/", clientHandler.CreateClient) // Deprecated
		clientGroup.PUT("/", clientHandler.UpdateClient)  // Deprecated

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
		userSessionGroup.GET("/all", sessionHandler.GetSessionsByProfile) // Deprecated
		userSessionGroup.POST("/", sessionHandler.CreateSession)          // Deprecated
		userSessionGroup.DELETE("/", sessionHandler.DeleteSession)        // Deprecated
		userSessionGroup.GET("/:id", sessionHandler.GetSession)           // Deprecated
	}

	uncategoryGroup := router.Group("/consultation")
	uncategoryGroup.Use(middleware.RoleRequire(""))
	{
		uncategoryGroup.GET("/therapist/:id", therapistHandler.GetTherapistByID) // Deprecated
		uncategoryGroup.GET("/client/:id", clientHandler.GetClientByID)          // Deprecated
	}

	api := router.Group("/v1/consultation")

	// Therapists
	therapist := api.Group("/therapists")
	therapist.Use(middleware.RoleRequire("therapist"))
	{
		therapist.GET("/me", therapistHandler.GetTherapist)
		therapist.POST("/me", therapistHandler.CreateTherapist)
		therapist.PUT("/me", therapistHandler.UpdateTherapist)
		therapist.PATCH("/me/availability", therapistHandler.UpdateAvailability)

		// admin / public
		therapist.GET("/:id", therapistHandler.GetTherapistByID)
	}

	// Clients
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

		// admin / public
		client.GET("/:id", clientHandler.GetClientByID)
	}

	// Sessions
	session := api.Group("/sessions")
	{
		session.GET("/me", sessionHandler.GetSessionsByProfile)
		session.POST("/me", sessionHandler.CreateSession)
		session.DELETE("/:id", sessionHandler.DeleteSession)
		session.GET("/:id", sessionHandler.GetSession) // admin / public
		session.POST("/:id/call/start", sessionHandler.StartCall)
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
		uncategoryGroupV1.GET("/therapist/:id", therapistHandler.GetTherapistByID)
		uncategoryGroupV1.GET("/me/recommend/top", swipeHandler.PopTop5V1)
		uncategoryGroupV1.GET("/client/:id", clientHandler.GetClientByID)
	}

	// Posts
	posts := api.Group("/posts")
	posts.Use(middleware.RoleRequire("")) // All authenticated users
	{
		posts.POST("/", postHandler.CreatePost)
		posts.GET("/", postHandler.GetFeed)                  // Personalized feed
		posts.GET("/trending", postHandler.GetTrendingPosts) // Trending posts
		posts.GET("/search", postHandler.SearchPosts)        // Search query
		posts.GET("/user/:userId", postHandler.GetUserPosts) // User's posts
		posts.GET("/:id", postHandler.GetPostByID)
		posts.PUT("/:id", postHandler.UpdatePost)
		posts.DELETE("/:id", postHandler.DeletePost)
		// Reactions
		posts.POST("/:id/react", reactionHandler.AddReaction)
		posts.DELETE("/:id/react", reactionHandler.RemoveReaction)
		posts.GET("/:id/reactions", reactionHandler.GetPostReactions)

		// Comments
		posts.POST("/:id/comments", commentHandler.CreateComment)
		posts.GET("/:id/comments", commentHandler.GetComments)

		// Share
		posts.POST("/:id/share", socialHandler.SharePost)

		// Bookmarks
		posts.POST("/:id/bookmark", socialHandler.AddBookmark)
		posts.DELETE("/:id/bookmark", socialHandler.RemoveBookmark)
	}

	// Comments (standalone routes for edit/delete/replies)
	comments := api.Group("/comments")
	comments.Use(middleware.RoleRequire(""))
	{
		comments.PUT("/:id", commentHandler.UpdateComment)
		comments.DELETE("/:id", commentHandler.DeleteComment)
		comments.GET("/:id/replies", commentHandler.GetReplies)
	}

	// Social Features
	users := api.Group("/users")
	users.Use(middleware.RoleRequire(""))
	{
		// Follow/Unfollow
		users.POST("/:id/follow", socialHandler.FollowUser)
		users.DELETE("/:id/follow", socialHandler.UnfollowUser)
		users.GET("/:id/followers", socialHandler.GetFollowers)
		users.GET("/:id/following", socialHandler.GetFollowing)
	}

	// Bookmarks
	bookmarks := api.Group("/bookmarks")
	bookmarks.Use(middleware.RoleRequire(""))
	{
		bookmarks.GET("/", socialHandler.GetBookmarks)
	}

	// Tags & Categories
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
