package main

import (
	"log"
	"os"

	middlewere "middlewere"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	db, err := middlewere.OpenConn()
	failOnError(err, "Failed to connect to PostgreSQL")
	defer db.Close()

	conn, err := amqp.Dial(rabbitURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"group",
		true,
		false,
		false,
		false,
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		},
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			if err := middlewere.ProcessDelivery(db, d); err != nil {
				log.Printf("failed to process delivery: %v", err)
				_ = d.Nack(false, true)
				continue
			}

			if err := d.Ack(false); err != nil {
				log.Printf("failed to ack message: %v", err)
				continue
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
