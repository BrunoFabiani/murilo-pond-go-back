package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type Telemetry struct {
	DeviceID      string    `json:"device_id"`
	Timestamp     time.Time `json:"timestamp"`
	SensorType    string    `json:"sensor_type"`
	ReadingNature string    `json:"reading_nature"`
	Value         any       `json:"value"`
}

var conn *amqp.Connection
var ch *amqp.Channel
var q amqp.Queue

func validateTelemetry(t *Telemetry) error {
	t.DeviceID = strings.TrimSpace(t.DeviceID)
	t.SensorType = strings.TrimSpace(t.SensorType)
	t.ReadingNature = strings.ToLower(strings.TrimSpace(t.ReadingNature))

	if t.DeviceID == "" {
		return errors.New("device_id is required")
	}
	if t.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}
	if t.SensorType == "" {
		return errors.New("sensor_type is required")
	}
	switch t.ReadingNature {
	case "analog":
		if _, ok := t.Value.(float64); !ok {
			return errors.New("analog value must be numeric")
		}
	case "discrete":
		switch t.Value.(type) {
		case bool, string:
		default:
			return errors.New("discrete value must be bool or string")
		}
	default:
		return errors.New("reading_nature must be analog or discrete")
	}
	return nil
}

func postTelemetry(c *gin.Context) {
	var msg Telemetry
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
		return
	}
	if err := validateTelemetry(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid telemetry fields"})
		return
	}

	body, err := json.Marshal(msg)
	failOnError(err, "Failed to marshal JSON")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s\n", body)
	c.IndentedJSON(http.StatusCreated, msg)
}

func main() {
	var err error
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err = amqp.Dial(rabbitURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err = ch.QueueDeclare(
		"group", // name
		true,    // durability
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		},
	)
	failOnError(err, "Failed to declare a queue")

	router := gin.Default()
	router.POST("/telemetry", postTelemetry)

	router.Run(":8080")
}
