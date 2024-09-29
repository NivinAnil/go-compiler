package factory

import (
	"go-compiler/common/pkg/utils"
	"go-compiler/execution-service/internal/adapter/clients/kubernetes"
	"go-compiler/execution-service/internal/adapter/clients/queue"
	"go-compiler/execution-service/internal/adapter/clients/queue/impl"
)

type AdapterFactory struct {
	QueueClient queue.IQueueClient
	Kubernetes  kubernetes.IKubernetesClient
}

func NewAdapterFactory() *AdapterFactory {
	log := utils.GetLogger()
	newQueue, _ := impl.NewQueueClient("executions", "amqp://guest:guest@localhost:5672/")
	kuberenetClient, kuberenetError := kubernetes.NewKubernetesClient("default")
	if kuberenetError != nil {
		log.Fatalf("Failed to create kubernetes client: %v", kuberenetError)
	}
	return &AdapterFactory{
		QueueClient: newQueue,
		Kubernetes:  kuberenetClient,
	}
}
