package main

import (
	"encoding/json"
	"fmt"
	"go-compiler/execution-service/internal/domain/dto/request"
	"go-compiler/execution-service/internal/ports/factory"
	"go-compiler/request-service/pkg/router"
	"log"
	"net/http"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	appRouter := router.GetRouter()

	// Listen to the queue in a separate goroutine
	go ListenToQueue("submissions", "amqp://guest:guest@localhost:5672/")

	port := ":8081"

	httpServer := &http.Server{
		Addr:    port,
		Handler: appRouter,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}

	fmt.Println("Server is running on port", port)
}

func ListenToQueue(queueName, rabbitMQURL string) {

	// Create ports for execution handling
	ports := factory.NewPortFactory()

	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Register a consumer
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Listen for messages from the queue
	for msg := range msgs {
		start := time.Now() // Start timing when a message is received
		log.Printf("Message received at: %v", start)

		// Unmarshal the incoming RabbitMQ message into the execution request
		var result request.NewExecutionRequest
		if err := json.Unmarshal(msg.Body, &result); err != nil {
			log.Printf("Error decoding RabbitMQ message: %v", err)
			continue
		}

		log.Printf("Decoded message: %+v", result)

		// Convert the execution request to string for handling
		resultStr, err := json.Marshal(result)
		if err != nil {
			log.Printf("Error converting result to string: %v", err)
			continue
		}

		// Call the execution service to handle the execution request
		handleStart := time.Now() // Start timing for handling the execution
		handlingError := ports.ExecutionHandler.Handle(string(resultStr))
		if handlingError != nil {
			log.Printf("Error handling message: %v", handlingError)
		}
		log.Printf("Execution handling completed in: %v", time.Since(handleStart))

		// Log the total time taken for processing the message
		log.Printf("Total time taken for message processing: %v", time.Since(start))
	}
}
