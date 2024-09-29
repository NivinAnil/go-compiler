package router

import (
	"go-compiler/common/pkg/utils"
	"go-compiler/request-service/internal/ports/factory"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())
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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}
