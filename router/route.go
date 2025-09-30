package router

import (
	"awesomeProject/controller"
	"awesomeProject/logger"
	"awesomeProject/middlewares"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(
		// logger.GinLogger(),
		logger.GinRecovery(true),
		middlewares.RateLimitMiddleware(2*time.Second, 100))

	v1 := r.Group("/api/v1")
	v1.POST("/login", controller.LoginHandler)
	v1.POST("/signup", controller.SignUpHandler)
	v1.GET("/refresh_token", controller.RefreshTokenHandler)

	v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
		//
		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.PostDetailHandler)
		//v1.GET("/post", controller.PostListHandler)

		v1.GET("/posts2", controller.PostList2Handler)
		//v1.GET("/post2", controller.PostList2Handler)
		//
		v1.POST("/vote", controller.VoteHandler)
		v1.POST("/follow", controller.FollowHandler)
		v1.POST("/unfollow", controller.UnfollowHandler)
		v1.POST("/feed/publish", controller.PublishMessageHandler)
		v1.GET("/feed/get", controller.GetFeedHandler)
		//
		//v1.POST("/comment", controller.CommentHandler)
		//v1.GET("/comment", controller.CommentListHandler)
		v1.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	return r
}
