package factory

import "go-compiler/request-service/internal/adapter/clients/queue"

type AdapterFactory struct {
	QueueClient queue.IQueueClient
}

func NewAdapterFactory() *AdapterFactory {
	return &AdapterFactory{
		QueueClient: queue.NewQueueClient(),
	}
}
