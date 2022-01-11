package broker

import (
	"encoding/json"
	"fmt"

	"github.com/disturb16/finechat/configuration"
	"github.com/streadway/amqp"
)

// MessageBroker is the interface for the broker.
type MessageBroker interface {
	Channel() (*amqp.Channel, error)
	SendMessage(exchange, key string, messageType MessageType, payload interface{}) error
	Close() error
}

// Broker is the message broker.
type Broker struct {
	conn *amqp.Connection
}

// MessageType is the type of the message.
type MessageType string

// Message is the basic format of the message.
type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

const (
	// TypeReload indicates that the chatroom should reload the messages.
	TypeReload MessageType = "reload"
	// TypeStockRequest corresponds to the stock commands.
	TypeStockRequest MessageType = "stock_request"
	// TypeCommandError indicates that the command was not understood.
	TypeCommandError MessageType = "command_error"
)

// New creates a new message broker.
func New(config *configuration.Configuration) (*Broker, error) {
	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		config.RabbitMQ.User,
		config.RabbitMQ.Password,
		config.RabbitMQ.Host,
		config.RabbitMQ.Port,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	b := &Broker{
		conn: conn,
	}

	return b, nil
}

// Channel returns the amqp channel.
func (b *Broker) Channel() (*amqp.Channel, error) {
	return b.conn.Channel()
}

// DefaultExchange sets the exchange with default configuration.
func DefaultExchange(ch *amqp.Channel, exchange string) error {
	return ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		false,    // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
}

// DefaultQueue sets the queue with default configuration.
func DefaultQueue(ch *amqp.Channel, name string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

// DefaultConsumer sets the consumer with default configuration.
func DefaultConsumer(ch *amqp.Channel, q amqp.Queue) (<-chan amqp.Delivery, error) {
	return ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
}

// SendMessage sends a message to the broker.
func (b *Broker) SendMessage(exchange, key string, messageType MessageType, payload interface{}) error {
	ch, err := b.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	err = DefaultExchange(ch, exchange)
	if err != nil {
		return err
	}

	msg := Message{
		Type:    messageType,
		Payload: payload,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return ch.Publish(
		exchange, // exchange
		key,      // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
}

// Close closes the broker.
func (b *Broker) Close() error {
	return b.conn.Close()
}
