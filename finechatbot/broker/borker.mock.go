package broker

import "github.com/streadway/amqp"

type MockBroker struct{}

func (mb *MockBroker) Channel() (*amqp.Channel, error) {
	return nil, nil
}

func (mb *MockBroker) SendMessage(exchange, key string, messageType MessageType, payload interface{}) error {
	return nil
}

func (mb *MockBroker) Close() error {
	return nil
}
