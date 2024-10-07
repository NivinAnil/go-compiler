package factory

import (
	"go-compiler/common/pkg/utils"
	"go-compiler/execution-service/internal/adapter/factory"
	"go-compiler/execution-service/internal/domain/services/impl"
	"go-compiler/execution-service/internal/domain/services/interfaces"
)

type DomainFactory struct {
	ExecutionService interfaces.IExecutionService
}

func NewDomainFactory() *DomainFactory {
	adapters := factory.NewAdapterFactory()
	cache := utils.NewCacheClient("localhost:6379", "", 0)
	return &DomainFactory{
		ExecutionService: impl.NewExecutionRequestService(adapters.QueueClient, cache),
	}
}
