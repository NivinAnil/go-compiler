package factory

import (
	"go-compiler/execution-service/internal/adapter/clients/queue"
	"go-compiler/execution-service/internal/adapter/clients/queue/impl"
)

type AdapterFactory struct {
	QueueClient queue.IQueueClient
}

func NewAdapterFactory() *AdapterFactory {
	newQueue, _ := impl.NewQueueClient("executions", "amqp://guest:guest@localhost:5672/")
	return &AdapterFactory{
		QueueClient: newQueue,
	}
}
