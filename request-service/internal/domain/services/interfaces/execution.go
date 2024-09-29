package interfaces

import (
	"context"
	"go-compiler/request-service/internal/domain/dto/request"
)

type IExecutionService interface {
	ProcessRequest(ctx context.Context, payload request.NewExecutionRequest) error
}
