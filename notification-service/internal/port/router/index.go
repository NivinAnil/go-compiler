package router

import (
	"encoding/json"
	"go-compiler/common/pkg/utils"
	"go-compiler/notification-service/internal/port/factory"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

// WebSocket upgrader to upgrade HTTP requests to WebSocket protocol
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin (change based on your CORS policy)
	},
}

type MessageBody struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	ConnectionID string `json:"connection_id,omitempty"`
}

type ExecutionBody struct {
	ConnectionID string `json:"connection_id"`
	Output       string `json:"output"`
}

// Global map to store active WebSocket connections by `connection_id`
var connections = make(map[string]*websocket.Conn)

var connMutex sync.Mutex // Mutex to handle concurrent access to `connections`

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	ports := factory.GetPorts()
	//health

	router.GET("/health", ports.HealthController.Status())

	// WebSocket route to listen for real-time execution results
	router.GET("/ws", func(c *gin.Context) {
		handleWebSocket(c.Writer, c.Request)
	})

	router.Use(CORSMiddleware())

	go ListenToQueue("executions", "amqp://guest:guest@rabbitmq-service:5672/")

	return router

}

func ListenToQueue(queueName, rabbitMQURL string) {
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

	for msg := range msgs {
		var result ExecutionBody
		if err := json.Unmarshal(msg.Body, &result); err != nil {
			log.Printf("Error decoding RabbitMQ message: %v", err)
			continue
		}

		// Forward message to WebSocket client based on connection_id
		SendMessageToClient(result.ConnectionID, result)
	}
}

// SendMessageToClient sends a message to the WebSocket client based on connection_id
func SendMessageToClient(connectionID string, result ExecutionBody) {
	connMutex.Lock()

	log := utils.GetLogger()
	conn, exists := connections[connectionID]
	connMutex.Unlock()

	if !exists {
		log.Info("No active WebSocket connection for connection_id: %s", connectionID)
		return
	}

	err := conn.WriteJSON(result)
	if err != nil {
		log.Info("Error sending message to WebSocket client %s: %v", connectionID, err)
		conn.Close()
		// Remove the connection from the map
		connMutex.Lock()
		delete(connections, connectionID)
		connMutex.Unlock()
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get connection_id from query parameters
	connectionID := r.URL.Query().Get("connection_id")
	if connectionID == "" {
		log.Printf("Missing connection_id")
		http.Error(w, "connection_id is required", http.StatusBadRequest)
		return
	}

	// Upgrade initial GET request to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Register the WebSocket connection in the global map
	connMutex.Lock()
	connections[connectionID] = conn
	connMutex.Unlock()

	defer func() {
		// Clean up when the connection is closed
		connMutex.Lock()
		delete(connections, connectionID)
		connMutex.Unlock()
	}()

	// Send a message to the client that just connected
	connectedMessage := MessageBody{
		Status:  "connected",
		Message: "You are successfully connected",
	}
	err = conn.WriteJSON(connectedMessage)
	if err != nil {
		log.Printf("Error sending connected message: %v", err)
		return
	}

	// Wait for messages from the client or close the connection if there's an error
	for {
		var msg MessageBody
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading JSON from WebSocket: %v", err)
			break
		}
		log.Printf("Received message: %v", msg)
	}
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}
