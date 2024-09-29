package factory

import (
	"go-compiler/notification-service/internal/adapter/clients/queue"
	"go-compiler/notification-service/internal/adapter/clients/queue/impl"
	"log"
)

type AdapterFactory struct {
	QueueClient queue.IQueueClient
}

func GetAdapters() *AdapterFactory {
	queueClient, err := impl.NewQueueClient("executions_result", "amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to create queue client: %v", err)
	}
	return &AdapterFactory{
		QueueClient: queueClient,
	}
}
