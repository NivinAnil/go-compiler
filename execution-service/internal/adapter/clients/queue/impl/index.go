package impl

import (
	"github.com/streadway/amqp"
	"log"
)

type QueueClient struct {
	QueueName  string `json:"queue_name"`
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewQueueClient(queueName string, rabbitMQURL string) (*QueueClient, error) {
	// Dial the RabbitMQ server
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		conn.Close()
		return nil, err
	}

	// Declare a queue to ensure it exists
	_, err = ch.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
		conn.Close()
		ch.Close()
		return nil, err
	}

	// Return the QueueClient with connection and channel established
	return &QueueClient{
		QueueName:  queueName,
		Connection: conn,
		Channel:    ch,
	}, nil
}

// PublishMessage sends a message to the queue.
func (qc *QueueClient) PublishMessage(messageBody string) error {
	err := qc.Channel.Publish(
		"",           // exchange
		qc.QueueName, // routing key (queue name)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(messageBody),
		})
	if err != nil {
		log.Printf("Failed to publish a message: %v", err)
		return err
	}
	log.Printf("Message published to queue %s: %s", qc.QueueName, messageBody)
	return nil
}

func (qc *QueueClient) ConsumeMessages(handleMessage func(string)) error {
	msgs, err := qc.Channel.Consume(
		qc.QueueName, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Printf("Failed to register a consumer: %v", err)
		return err
	}

	go func() {
		for d := range msgs {
			handleMessage(string(d.Body))
		}
	}()

	log.Printf("Waiting for messages from queue %s", qc.QueueName)
	return nil
}

func (qc *QueueClient) Close() error {
	err := qc.Channel.Close()
	if err != nil {
		log.Printf("Failed to close channel: %v", err)
		return err
	}

	err = qc.Connection.Close()
	if err != nil {
		log.Printf("Failed to close connection: %v", err)
		return err
	}

	log.Printf("Successfully closed connection to RabbitMQ")
	return nil
}
