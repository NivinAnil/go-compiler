package pkg

import (
	"go-compiler/notification-service/internal/port/router"

	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	engine := router.NewRouter()
	return engine
}
