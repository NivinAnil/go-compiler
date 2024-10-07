package interfaces

import (
	"context"
	"go-compiler/execution-service/internal/domain/dto/request"
)

type IExecutionService interface {
	HandleExecution(ctx context.Context, payload request.NewExecutionRequest) error
	GetExecution(ctx context.Context, requestId string) (interface{}, error)
}
