package controllers

import (
	"go-compiler/common/pkg/utils/logger"
	"go-compiler/execution-service/internal/domain/services/interfaces"

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

func (rc *RequestController) GetExecution() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log := logger.GetLogger(ctx)
		methodName := "GetExecution"
		log.Info("Entering", "methodName", methodName)
		requestId := ctx.Param("request_id")

		resp, domainError := rc.ExecutionService.GetExecution(ctx, requestId)
		if domainError != nil {
			log.Error("Error in processing request", "error", domainError.Error())
			ctx.JSON(500, gin.H{"error": domainError.Error()})
			return
		}

		ctx.JSON(200, resp)
	}
}
