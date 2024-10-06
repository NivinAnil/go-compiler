package impl

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-compiler/common/pkg/utils/logger"
	"go-compiler/execution-service/internal/adapter/clients/kubernetes"
	"go-compiler/execution-service/internal/adapter/clients/queue"
	"go-compiler/execution-service/internal/domain/dto/request"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type ExecutionRequestService struct {
	QueueClient      queue.IQueueClient
	KubernetesClient kubernetes.IKubernetesClient
}

func NewExecutionRequestService(qc queue.IQueueClient, kc kubernetes.IKubernetesClient) *ExecutionRequestService {
	return &ExecutionRequestService{
		QueueClient:      qc,
		KubernetesClient: kc,
	}
}

func (s *ExecutionRequestService) HandleExecution(ctx context.Context, payload request.NewExecutionRequest) error {
	log := logger.GetLogger(ctx)
	methodName := "HandleExecution"
	start := time.Now()
	log.Info("Entering", "methodName", methodName, "start_time", start)

	// Call the Python processing function
	resp := s.ProcessPython(payload.Code, payload.StdIn, payload.ConnectionId)
	if resp == "" {
		log.Error("Error processing request")
		return fmt.Errorf("Error processing request")
	}

	// Log the total time taken for the HandleExecution function
	log.Info("Exiting", "methodName", methodName, "total_time_taken", time.Since(start))
	return nil
}

func (e *ExecutionRequestService) ProcessPython(Code string, Stdin string, ConnectionId string) string {
	log := logger.GetLogger()
	processStart := time.Now()
	log.Info("Inside ProcessPython", "start_time", processStart)

	// Step 1: Decode the base64-encoded Python code
	decodeStart := time.Now()
	decodedCodeBytes, err := base64.StdEncoding.DecodeString(Code)
	if err != nil {
		log.Error("Error decoding base64 code:", err)
		return ""
	}
	log.Info("Decoded base64 code", "time_taken", time.Since(decodeStart))

	// Step 2: Decode the base64-encoded stdin input (if provided)
	decodeStdinStart := time.Now()
	var stdin string
	if Stdin != "" {
		decodedStdinBytes, err := base64.StdEncoding.DecodeString(Stdin)
		if err != nil {
			log.Error("Error decoding base64 stdin:", err)
			return ""
		}
		stdin = string(decodedStdinBytes)
	} else {
		stdin = "" // No stdin provided
	}
	log.Info("Decoded base64 stdin", "time_taken", time.Since(decodeStdinStart))

	// Step 3: Define the container image and command to run the code
	cmd := []string{"bash", "-c", fmt.Sprintf("echo '%s' | python -c '%s'", stdin, string(decodedCodeBytes))}
	env := []string{"CODE=" + string(decodedCodeBytes)}

	// Step 4: Log the job creation start time
	jobStart := time.Now()
	log.Info("Starting Kubernetes Job", "start_time", jobStart)

	// Step 5: Create the Kubernetes Job to execute the Python code
	jobName, createErr := e.KubernetesClient.CreateJob("ammyy9908/go-python", cmd, env)
	if createErr != nil {
		log.Error("Error creating Kubernetes Job", createErr)
		return ""
	}
	log.Info("Kubernetes Job created", "job_name", jobName, "time_taken", time.Since(jobStart))

	// Step 6: Wait for the Job to complete
	waitStart := time.Now()
	log.Info("Waiting for the Job to finish execution")
	waitErr := e.KubernetesClient.WaitForJobCompletion(jobName)
	if waitErr != nil {
		log.Error("Error waiting for Kubernetes Job completion", waitErr)
		return ""
	}
	log.Info("Kubernetes Job execution completed", "time_taken", time.Since(waitStart))

	// Step 7: Retrieve the logs from the Job's Pod
	logRetrievalStart := time.Now()
	log.Info("Retrieving logs from the Job")
	logs, logsErr := e.KubernetesClient.GetJobLogs(jobName)
	if logsErr != nil {
		log.Error("Error retrieving logs", logsErr)
		return ""
	}
	log.Info("Logs retrieved successfully", "logs", logs, "time_taken", time.Since(logRetrievalStart))

	// Step 8: Publish the execution result to the queue
	publishStart := time.Now()
	newExecution := bson.M{
		"connection_id": ConnectionId,
		"output":        logs,
	}

	Payload, err := json.Marshal(newExecution)
	if err != nil {
		log.Error("Error marshalling execution result", err)
		return ""
	}
	publishError := e.QueueClient.PublishMessage(string(Payload))
	if publishError != nil {
		log.Error("Error publishing execution result to queue", publishError)
		return ""
	}
	log.Info("Execution result published to queue", "time_taken", time.Since(publishStart))

	// Step 9: Clean up the Kubernetes Job after execution
	cleanupStart := time.Now()
	log.Info("Removing the Kubernetes Job")
	removeErr := e.KubernetesClient.DeleteJob(jobName)
	if removeErr != nil {
		log.Error("Error deleting Kubernetes Job", removeErr)
		return ""
	}
	log.Info("Kubernetes Job removed successfully", "time_taken", time.Since(cleanupStart))

	// Step 10: Return the logs (output from the Python code)
	log.Info("ProcessPython completed", "total_time_taken", time.Since(processStart))
	return logs
}
