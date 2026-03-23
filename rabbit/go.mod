module rabbit

go 1.26.1

require (
	github.com/rabbitmq/amqp091-go v1.10.0
	middlewere v0.0.0
)

require github.com/lib/pq v1.12.0 // indirect

replace middlewere => ../middlewere
