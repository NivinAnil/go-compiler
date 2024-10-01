package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/streadway/amqp"
)

type NewExecutionRequest struct {
	Code         string `json:"code"`
	StdIn        string `json:"stdin"`
	ConnectionId string `json:"connection_id"`
	LanguageId   int    `json:"language_id"`
}

type ExecutionResult struct {
	ConnectionId string `json:"connection_id"`
	Output       string `json:"output"`
}

func main() {
	// Fetch RabbitMQ URL from the environment variable
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		log.Fatal("RABBITMQ_URL environment variable not set")
	}

	// Log the RabbitMQ URL to verify it during troubleshooting
	log.Printf("Connecting to RabbitMQ at: %s", rabbitMQURL)

	// Connect to RabbitMQ using the URL from the environment variable
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Start the worker to listen to the submissions queue
	go listenForTasks(ch)

	// Keep the worker running
	select {}
}

func listenForTasks(ch *amqp.Channel) {
	// Declare the queue to ensure it exists
	queue, err := ch.QueueDeclare(
		"submissions", // name of the queue
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	log.Printf("Declared queue: %s", queue.Name)

	// Now consume from the queue
	msgs, err := ch.Consume(
		"submissions", // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for msg := range msgs {
		start := time.Now()
		var req NewExecutionRequest
		if err := json.Unmarshal(msg.Body, &req); err != nil {
			log.Printf("Error decoding RabbitMQ message: %v", err)
			continue
		}

		// Process the execution request
		output := processExecution(req)

		// Send the result back to the executions queue
		sendResult(ch, req.ConnectionId, output)

		log.Printf("Execution completed in: %v", time.Since(start))
	}
}

func processExecution(req NewExecutionRequest) string {
	// Decode the base64-encoded code
	decodedCode, err := base64.StdEncoding.DecodeString(req.Code)
	if err != nil {
		return fmt.Sprintf("Error decoding base64 code: %v", err)
	}

	// Initialize command based on the LanguageId
	var cmd *exec.Cmd
	switch req.LanguageId {
	case 1: // Python
		cmd = exec.Command("python3", "-c", string(decodedCode))
	case 2: // JavaScript (Node.js)
		cmd = exec.Command("node", "-e", string(decodedCode))
	case 3: // Go
		tempFile := "/tmp/temp.go"
		err := os.WriteFile(tempFile, decodedCode, 0644)
		if err != nil {
			return fmt.Sprintf("Error writing Go code to file: %v", err)
		}
		cmd = exec.Command("go", "run", tempFile)
	case 4: // Java
		tempFile := "/tmp/Main.java"
		err := os.WriteFile(tempFile, decodedCode, 0644)
		if err != nil {
			return fmt.Sprintf("Error writing Java code to file: %v", err)
		}
		// Compile Java code
		compileCmd := exec.Command("javac", tempFile)
		_, err = compileCmd.CombinedOutput()
		if err != nil {
			return fmt.Sprintf("Error compiling Java code: %v", err)
		}
		cmd = exec.Command("java", "-cp", "/tmp", "Main")
	case 9: // Bash
		// Write the Bash code to a temporary file
		tempFile := "/tmp/temp.sh"
		err := os.WriteFile(tempFile, decodedCode, 0755) // 0755 permissions to execute
		if err != nil {
			return fmt.Sprintf("Error writing Bash code to file: %v", err)
		}
		// Use bash to execute the script
		cmd = exec.Command("bash", tempFile)
	default:
		return "Unsupported language"
	}

	// Set up stdin for the command if provided
	if req.StdIn != "" {
		cmd.Stdin = bytes.NewBufferString(req.StdIn)
	}

	// Execute the code
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error executing code: %v\nOutput: %s", err, string(output))
	}

	return string(output)
}
func sendResult(ch *amqp.Channel, connectionId, output string) {
	result := ExecutionResult{
		ConnectionId: connectionId,
		Output:       output,
	}

	body, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error marshalling execution result: %v", err)
		return
	}

	err = ch.Publish(
		"",           // exchange
		"executions", // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish execution result: %v", err)
	}
}
