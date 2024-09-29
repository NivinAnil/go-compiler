package factory

import "go-compiler/notification-service/internal/port/controller/health"

type PortFactory struct {
	HealthController health.HealthController
}

func GetPorts() *PortFactory {
	return &PortFactory{
		HealthController: *health.NewHealthController(),
	}
}
