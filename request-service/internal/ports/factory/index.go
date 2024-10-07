package factory

import (
	"go-compiler/request-service/internal/domain/factory"
	"go-compiler/request-service/internal/ports/controllers"
	"go-compiler/request-service/internal/ports/grpc"
)

type PortFactory struct {
	RequestController controllers.RequestController
	GRPCServer        *grpc.GRPCServer
}

func NewPortFactory() *PortFactory {
	domains := factory.NewDomainFactory()
	return &PortFactory{
		RequestController: *controllers.NewRequestController(domains.ExecutionService),
		GRPCServer:        grpc.NewGRPCServer(domains.ExecutionService),
	}
}
