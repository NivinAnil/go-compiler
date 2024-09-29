package router

import (
	"go-compiler/common/pkg/utils"
	"go-compiler/execution-service/internal/ports/factory"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	rateLimiter := RateLimiterMiddleware(10, time.Minute)
	router.Use(rateLimiter)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	portFactory := factory.NewPortFactory()

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/submission", portFactory.RequestController.GetRequest())
		}
	}

	return router
}

func RateLimiterMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	limiter := utils.NewRateLimiter(maxRequests, window)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}
		c.Next()
	}
}
