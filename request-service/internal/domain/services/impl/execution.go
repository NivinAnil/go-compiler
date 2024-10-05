package impl

import (
	"context"
	"encoding/json"
	"go-compiler/common/pkg/utils/logger"
	"go-compiler/request-service/internal/adapter/clients/queue"
	"go-compiler/request-service/internal/domain/dto/request"
)

type ExecutionRequestService struct {
	QueueClient queue.IQueueClient
}

func NewExecutionRequestService(qc queue.IQueueClient, url string) *ExecutionRequestService {
	er := qc.Connect(url)
	log := logger.GetLogger(context.Background())
	if er != nil {
		log.Error("Error in connecting to queue", "error", er.Error())
		return nil
	}
	return &ExecutionRequestService{
		QueueClient: qc,
	}
}

func (s *ExecutionRequestService) ProcessRequest(ctx context.Context, payload request.NewExecutionRequest) error {
	log := logger.GetLogger(ctx)
	methodName := "ProcessRequest"
	log.Info("Entering", "methodName", methodName)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Error("Error marshalling payload", "error", err.Error())
		return err
	}

	err = s.QueueClient.SendMessage("submissions", payloadBytes)
	if err != nil {
		log.Error("Error sending message to queue", "error", err.Error())
		return err
	}

	return nil
}
