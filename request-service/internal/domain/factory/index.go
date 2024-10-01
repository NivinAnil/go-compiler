package factory

import (
	"go-compiler/request-service/internal/adapter/factory"
	"go-compiler/request-service/internal/domain/services/impl"
	"go-compiler/request-service/internal/domain/services/interfaces"
)

type DomainFactory struct {
	ExecutionService interfaces.IExecutionService
}

func NewDomainFactory() *DomainFactory {
	adapters := factory.NewAdapterFactory()
	return &DomainFactory{
		ExecutionService: impl.NewExecutionRequestService(adapters.QueueClient, "amqp://guest:guest@rabbitmq-service:5672/"),
	}
}
