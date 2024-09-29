package queue

import (
	"github.com/streadway/amqp"
)

// IQueueClient defines the interface for a queue client
type IQueueClient interface {
	Connect(url string) error
	SendMessage(queueName string, message []byte) error
	ReceiveMessage(queueName string) (<-chan amqp.Delivery, error)
}

// QueueClient implements the IQueueClient interface for RabbitMQ
type QueueClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewQueueClient creates a new instance of QueueClient
func NewQueueClient() *QueueClient {
	return &QueueClient{}
}

// Connect establishes a connection to the RabbitMQ server
func (qc *QueueClient) Connect(url string) error {
	var err error
	qc.connection, err = amqp.Dial(url)
	if err != nil {
		return err
	}

	qc.channel, err = qc.connection.Channel()
	if err != nil {
		return err
	}

	return nil
}

// SendMessage sends a message to the specified queue
func (qc *QueueClient) SendMessage(queueName string, message []byte) error {
	_, err := qc.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = qc.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	return err
}

// ReceiveMessage receives messages from the specified queue
func (qc *QueueClient) ReceiveMessage(queueName string) (<-chan amqp.Delivery, error) {
	_, err := qc.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	msgs, err := qc.channel.Consume(
		queueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

// Close closes the connection and channel
func (qc *QueueClient) Close() {
	if qc.channel != nil {
		qc.channel.Close()
	}
	if qc.connection != nil {
		qc.connection.Close()
	}
}
