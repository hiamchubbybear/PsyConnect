package route

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"consultationservice/bootstrap"
	handlers "consultationservice/internal/handler"
	"consultationservice/internal/middleware"
	"consultationservice/pkg/logger"
)

func RouterInit(
	env *bootstrap.Env,
	kafkaLogger *logger.KafkaLogger,
	clientHandler *handlers.ClientHandler,
	therapistHandler *handlers.TherapistHandler,
	matchingHandler *handlers.MatchHandler,
	sessionHandler *handlers.SessionHandler,
	swipeHandler *handlers.SwipeHandler,
	// Newsfeed handlers
	postHandler *handlers.PostHandler,
	reactionHandler *handlers.ReactionHandler,
	commentHandler *handlers.CommentHandler,
	socialHandler *handlers.SocialHandler,
) {
	urI := fmt.Sprintf("%v:%v", env.Addr, env.Port)
	router := gin.Default()

	// Add logging middleware
	router.Use(middleware.LoggingMiddleware(kafkaLogger))

	defer func() {
		if err := recover(); err != nil {
			kafkaLogger.Fatal("Router panic", map[string]interface{}{
				"error": fmt.Sprintf("%v", err),
			})
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
		client.GET("/me", clientHandler.GetClientHandler)
		client.POST("/me", clientHandler.PostClientHandler)
		client.PUT("/me", clientHandler.PutClientHandler)

		client.POST("/me/recommend", swipeHandler.TriggerUpdateV1)
		client.GET("/me/recommend/top", swipeHandler.PopTop5V1)
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


	// Posts
	posts := api.Group("/posts")
	posts.Use(middleware.RoleRequire("")) // All authenticated users
	{
		posts.POST("/", postHandler.CreatePost)
		posts.GET("/", postHandler.GetFeed)                          // Personalized feed
		posts.GET("/trending", postHandler.GetTrendingPosts)         // Trending posts
		posts.GET("/search", postHandler.SearchPosts)                // Search query
		posts.GET("/user/:userId", postHandler.GetUserPosts)         // User's posts
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
