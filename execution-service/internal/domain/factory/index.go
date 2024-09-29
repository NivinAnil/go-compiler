package factory

import (
	"go-compiler/execution-service/internal/adapter/factory"
	"go-compiler/execution-service/internal/domain/services/impl"
	"go-compiler/execution-service/internal/domain/services/interfaces"
)

type DomainFactory struct {
	ExecutionService interfaces.IExecutionService
}

func NewDomainFactory() *DomainFactory {
	adapters := factory.NewAdapterFactory()
	return &DomainFactory{
		ExecutionService: impl.NewExecutionRequestService(adapters.QueueClient, adapters.Kubernetes),
	}
}
