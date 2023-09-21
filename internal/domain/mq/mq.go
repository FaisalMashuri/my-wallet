package mq

type MessageQueueEmitter struct {
	Title   string
	Payload map[string]interface{}
}

type MQService interface {
	SendData(string, []byte) error
}
