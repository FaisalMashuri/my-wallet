package infrastructure

import (
	"fmt"
	"github.com/FaisalMashuri/my-wallet/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   *amqp.Queue
}

func NewRabbitMQ(cfg config.Config) *RabbitMQ {
	conString := cfg.DSN_MQ

	conn, err := amqp.Dial(conString)
	if err != nil {
		return nil
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil
	}
	log.Println("RabbitMQ connection established")
	err = ch.ExchangeDeclare("user", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("RabbitMQ connection established")
	return &RabbitMQ{
		Conn:    conn,
		Queue:   nil,
		Channel: ch,
	}
}
