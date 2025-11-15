package database

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func queuerInstance() *amqp091.Channel {
	rabbitMQ := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp091.Dial(rabbitMQ)
	if err != nil {
		fmt.Printf("Rabbit MQ Server failed: %s", err.Error())
		panic("Job Queuer Failed")
	}
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Rabbit MQ Server failed: %s", err.Error())
		panic("Job Queuer Channel Connection Failed")
	}
	_, err = ch.QueueDeclare(
		"jobs", // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	ch.Qos(1, 0, false)
	if err != nil {
		fmt.Printf("Rabbit MQ Cannot declare queue: %s", err.Error())
		panic("Job Queuer Channel Connection Failed")
	}
	return ch
}

var QueuerChannel = queuerInstance()

func ListQueues(conn *amqp091.Connection) {

}
