package router

import (
	"go-compiler/request-service/internal/ports/router"

	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	return router.NewRouter()
}
