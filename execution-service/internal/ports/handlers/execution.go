package handlers

import (
	"context"
	"encoding/json"
	"go-compiler/common/pkg/utils"
	"go-compiler/execution-service/internal/domain/dto/request"
	"go-compiler/execution-service/internal/domain/services/interfaces"
	"time"
)

type ExecutionHandler struct {
	ExecutionService interfaces.IExecutionService
}

func NewExecutionHandler(es interfaces.IExecutionService) *ExecutionHandler {
	return &ExecutionHandler{
		ExecutionService: es,
	}
}

func (handler *ExecutionHandler) Handle(payload string) error {
	// Create a new context
	ctx := context.Background()

	// Get the logger
	log := utils.GetLogger(ctx)
	methodName := "Handle"
	start := time.Now()
	log.Info("Entering", "methodName", methodName, "start_time", start)

	// Start measuring the unmarshalling time
	unmarshalStart := time.Now()
	var RequestBody *request.NewExecutionRequest
	err := json.Unmarshal([]byte(payload), &RequestBody)
	if err != nil {
		log.Error("Unmarshalling Error", "error", err)
		return nil
	}
	log.Info("Unmarshalling completed", "time_taken", time.Since(unmarshalStart))

	// Start measuring the time taken by HandleExecution
	handleExecutionStart := time.Now()
	resp := handler.ExecutionService.HandleExecution(ctx, *RequestBody)
	if resp != nil {
		log.Error("Error in processing request", "error", resp.Error())
		return nil
	}
	log.Info("ExecutionService.HandleExecution completed", "time_taken", time.Since(handleExecutionStart))

	// Log the total time taken by the Handle function
	log.Info("Exiting", "methodName", methodName, "total_time_taken", time.Since(start))
	return nil
}
