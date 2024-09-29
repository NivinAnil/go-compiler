package factory

import (
	"go-compiler/execution-service/internal/domain/factory"
	"go-compiler/execution-service/internal/ports/controllers"
	"go-compiler/execution-service/internal/ports/handlers"
)

type PortFactory struct {
	RequestController controllers.RequestController
	ExecutionHandler  handlers.ExecutionHandler
}

func NewPortFactory() *PortFactory {
	domains := factory.NewDomainFactory()
	return &PortFactory{
		RequestController: *controllers.NewRequestController(domains.ExecutionService),
		ExecutionHandler:  *handlers.NewExecutionHandler(domains.ExecutionService),
	}
}
