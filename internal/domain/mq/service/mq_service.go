package service

import (
	"context"
	"fmt"
	"github.com/FaisalMashuri/my-wallet/infrastructure"
	"github.com/FaisalMashuri/my-wallet/internal/domain/mq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type mqService struct {
	mq *infrastructure.RabbitMQ
}

func NewMqService(mq *infrastructure.RabbitMQ) mq.MQService {
	return &mqService{
		mq: mq,
	}
}

func (s *mqService) SendData(key string, payload []byte) error {
	ctx := context.Background()
	err := s.mq.Channel.PublishWithContext(ctx, "user", key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        payload,
	})
	if err != nil {
		fmt.Println("ERROR publish topic : ", err)
		return err
	}
	return nil
}
