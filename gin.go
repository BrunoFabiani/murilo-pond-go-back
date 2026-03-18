package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    amqp "github.com/rabbitmq/amqp091-go"
)



func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
// album represents data about a record album.
type people struct {
    ID     string  `json:"id"`
    Name  string  `json:"name"`
}

// albums slice to seed record album data.
var group = []people{
    
}

var conn *amqp.Connection
var ch *amqp.Channel
var q amqp.Queue

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {

    
    
    var newperson people

    // Call BindJSON to bind the received JSON to
    // newAlbum.
    if err := c.BindJSON(&newperson); err != nil {
        return
    }

    body, err := json.Marshal(newperson)
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

    // Add the new album to the slice.
    group = append(group, newperson)
    c.IndentedJSON(http.StatusCreated, newperson)
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
    router.POST("/group", postAlbums)
    
    router.Run(":8080")
}
