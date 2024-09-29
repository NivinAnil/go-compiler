package factory

import (
	"go-compiler/request-service/internal/domain/factory"
	"go-compiler/request-service/internal/ports/controllers"
)

type PortFactory struct {
	RequestController controllers.RequestController
}

func NewPortFactory() *PortFactory {
	domains := factory.NewDomainFactory()
	return &PortFactory{
		RequestController: *controllers.NewRequestController(domains.ExecutionService),
	}
}
