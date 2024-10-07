package impl

import (
	"context"
	"encoding/base64"
	"go-compiler/common/pkg/utils"
	"go-compiler/common/pkg/utils/logger"
	"go-compiler/execution-service/internal/adapter/clients/kubernetes"
	"go-compiler/execution-service/internal/adapter/clients/queue"
	"go-compiler/execution-service/internal/domain/dto/request"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"os/exec"
	"time"
)

type ExecutionRequestService struct {
	QueueClient      queue.IQueueClient
	KubernetesClient kubernetes.IKubernetesClient
	cache            utils.ICacheClient
}

func NewExecutionRequestService(qc queue.IQueueClient, kc kubernetes.IKubernetesClient, c utils.ICacheClient) *ExecutionRequestService {
	return &ExecutionRequestService{
		QueueClient:      qc,
		KubernetesClient: kc,
		cache:            c,
	}
}

func (s *ExecutionRequestService) HandleExecution(ctx context.Context, payload request.NewExecutionRequest) error {
	log := logger.GetLogger(ctx)
	methodName := "HandleExecution"
	log.Info("Entering", "methodName", methodName)

	switch payload.LanguageId {
	case 1:
		err := s.ProcessPythonRequest(ctx, payload)
		if err != nil {
			log.Error("Error processing Python request", "error", err)
			return err
		}

	case 2:
		err := s.ProcessJavaScriptRequest(ctx, payload)
		if err != nil {
			log.Error("Error processing JavaScript request", "error", err)
			return err
		}

	case 3:
		err := s.ProcessBashRequest(ctx, payload)
		if err != nil {
			log.Error("Error processing Bash request", "error", err)
			return err
		}
	}

	return nil
}
func (e *ExecutionRequestService) ProcessPythonRequest(ctx context.Context, payload request.NewExecutionRequest) error {
	log := logger.GetLogger(ctx)
	methodName := "ProcessPythonRequest"
	log.Info("Entering", "methodName", methodName)

	// Decode the base64-encoded Python code
	decodedCode, err := base64.StdEncoding.DecodeString(payload.Code)
	if err != nil {
		log.Error("Error decoding base64 string", "error", err)
		return err
	}

	// Create the command to execute the decoded Python code
	cmd := exec.Command("python3", "-c", string(decodedCode))

	// Get the stdin pipe for sending input to the Python script
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Error("Error creating stdin pipe", "error", err)
		return err
	}

	// Run the command asynchronously and provide the input from payload.Stdin
	go func() {
		defer stdin.Close()
		if payload.StdIn != "" {
			// Write the provided input from payload.StdIn to the stdin of the Python script
			io.WriteString(stdin, payload.StdIn)
		}
	}()

	// Capture the combined stdout and stderr output
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Error executing code", "error", err, "output", string(output))
		return err
	}

	expiration := 1 * time.Hour
	// Push the result to the cache
	err = e.cache.Set(payload.RequestId, string(output), expiration)
	if err != nil {
		log.Error("Error setting cache", "error", err)
		return err
	}

	return nil
}
func (e *ExecutionRequestService) ProcessJavaScriptRequest(ctx context.Context, payload request.NewExecutionRequest) error {
	log := logger.GetLogger(ctx)
	methodName := "ProcessJavaScriptRequest"
	log.Info("Entering", "methodName", methodName)

	decodedCode, err := base64.StdEncoding.DecodeString(payload.Code)
	if err != nil {
		log.Error("Error decoding base64 string", "error", err)
		return err
	}

	// Create the command to execute the JavaScript code
	cmd := exec.Command("node", "-e", string(decodedCode))

	// Get the stdin pipe for sending input to the JavaScript script
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Error("Error creating stdin pipe", "error", err)
		return err
	}

	// Run the command asynchronously and provide the input from payload.StdIn
	go func() {
		defer stdin.Close()
		if payload.StdIn != "" {
			// Write the provided input from payload.StdIn to the stdin of the JavaScript script
			io.WriteString(stdin, payload.StdIn)
		}
	}()

	// Capture the combined stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		//log the error but save the error as output in cache
		log.Error("Error executing code", "error", err, "output", string(output))
		expiration := 1 * time.Hour
		// Push the result to the cache
		err = e.cache.Set(payload.RequestId, string(output), expiration)
		if err != nil {
			log.Error("Error setting cache", "error", err)
			return err
		}

	}
	expiration := 1 * time.Hour
	// Push the result to the cache
	err = e.cache.Set(payload.RequestId, string(output), expiration)
	if err != nil {
		log.Error("Error setting cache", "error", err)
		return err
	}

	return nil
}

func (e *ExecutionRequestService) ProcessBashRequest(ctx context.Context, payload request.NewExecutionRequest) error {
	log := logger.GetLogger(ctx)
	methodName := "ProcessBashRequest"
	log.Info("Entering", "methodName", methodName)
	decodedCode, err := base64.StdEncoding.DecodeString(payload.Code)
	if err != nil {
		log.Error("Error decoding base64 string", "error", err)
		return err
	}
	// Create the command to execute the Bash code
	cmd := exec.Command("bash", "-c", string(decodedCode))

	// Get the stdin pipe for sending input to the Bash script
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Error("Error creating stdin pipe", "error", err)
		return err
	}

	// Run the command asynchronously and provide the input from payload.StdIn
	go func() {
		defer stdin.Close()
		if payload.StdIn != "" {
			// Write the provided input from payload.StdIn to the stdin of the Bash script
			io.WriteString(stdin, payload.StdIn)
		}
	}()

	// Capture the combined stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Error executing code", "error", err, "output", string(output))
		expiration := 1 * time.Hour
		// Push the result to the cache
		err = e.cache.Set(payload.RequestId, string(output), expiration)
		if err != nil {
			log.Error("Error setting cache", "error", err)
			return err
		}
	}
	expiration := 1 * time.Hour
	// Push the result to the cache
	err = e.cache.Set(payload.RequestId, string(output), expiration)
	if err != nil {
		log.Error("Error setting cache", "error", err)
		return err
	}

	return nil
}
func (e *ExecutionRequestService) GetExecution(ctx context.Context, requestId string) (interface{}, error) {
	log := logger.GetLogger(ctx)
	methodName := "GetExecution"
	log.Info("Entering", "methodName", methodName)

	// Get the result from the cache
	result, err := e.cache.Get(requestId)
	if err != nil {
		log.Error("Error getting cache", "error", err)
		return nil, err
	}
	return bson.M{
		"output": result,
	}, nil
}
