package controllers

import (
	"go-compiler/common/pkg/utils/logger"
	"go-compiler/request-service/internal/domain/dto/request"
	"go-compiler/request-service/internal/domain/services/interfaces"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestController struct {
	ExecutionService interfaces.IExecutionService
}

func NewRequestController(es interfaces.IExecutionService) *RequestController {
	return &RequestController{
		ExecutionService: es,
	}
}

func (rc *RequestController) GetRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := logger.GetLogger(ctx)
		methodName := "GetRequest"
		log.Info("Entering", "methodName", methodName)
		start := time.Now()
		log.Info("Request received at: %v", start)

		var Payload request.NewExecutionRequest
		e := ctx.BindJSON(&Payload)
		if e != nil {
			log.Error("Error in binding request", "error", e.Error())
			ctx.JSON(400, gin.H{"error": e.Error()})
			return
		}

		domainError := rc.ExecutionService.ProcessRequest(ctx, Payload)
		if domainError != nil {
			log.Error("Error in processing request", "error", domainError.Error())
			ctx.JSON(500, gin.H{"error": domainError.Error()})
			return
		}

		log.Info("Request submitted to RabbitMQ at: %v", time.Since(start))

		ctx.JSON(200, gin.H{"message": "Request processed successfully"})
	}
}
