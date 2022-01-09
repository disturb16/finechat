package broker

import (
	"encoding/json"
	"fmt"

	"github.com/disturb16/finechat/configuration"
	"github.com/streadway/amqp"
)

type Broker struct {
	amqpURL string
	conn    *amqp.Connection
}

type MessageType string

type Message struct {
	Type    MessageType `json:"type"`
	Payload string      `json:"payload"`
}

const (
	TypeReload       MessageType = "reload"
	TypeStockRequest MessageType = "stock_request"
)

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
		amqpURL: url,
		conn:    conn,
	}

	return b, nil
}

func (b *Broker) Channel() (*amqp.Channel, error) {
	return b.conn.Channel()
}

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

func (b *Broker) SendMessage(exchange, key string, messageType MessageType, message string) error {
	ch, err := b.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	err = DefaultExchange(ch, exchange)
	if err != nil {
		return err
	}

	m := &Message{
		Type:    messageType,
		Payload: message,
	}

	body, err := json.Marshal(m)
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
