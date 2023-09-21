package infrastructure

import (
	"fmt"
	"github.com/FaisalMashuri/my-wallet/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

var MessageBroker = NewRabbitMQ()

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   *amqp.Queue
}

func NewRabbitMQ() *RabbitMQ {
	conString := fmt.Sprintf("amqps://%s:%s@%s/%s", config.AppConfig.RabbitMQConfig.User, config.AppConfig.RabbitMQConfig.Password, config.AppConfig.RabbitMQConfig.Host, config.AppConfig.RabbitMQConfig.User)
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
	return &RabbitMQ{
		Conn:    conn,
		Queue:   nil,
		Channel: ch,
	}
}
