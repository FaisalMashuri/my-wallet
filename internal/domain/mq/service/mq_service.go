package service

import (
	"context"
	"github.com/FaisalMashuri/my-wallet/infrastructure"
	"github.com/FaisalMashuri/my-wallet/internal/domain/mq"
	"github.com/rabbitmq/amqp091-go"
	"time"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.mq.Channel.PublishWithContext(ctx, "user", key, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        payload,
	})
	if err != nil {
		return err
	}
	return nil
}
