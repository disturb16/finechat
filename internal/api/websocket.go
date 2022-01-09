package api

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
)

var (
	upgrader = websocket.Upgrader{}
)

func (h *Handler) socket(c echo.Context) error {

	var err error
	defer func() {
		log.Println(err)
	}()
	// Upgrade from http request to websocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	defer ws.Close()
	log.Println("New connection")

	// Connect  to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@192.168.0.11:5672/")
	if err != nil {
		return err
	}

	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	//TODO: get chatroom name.
	id := c.Param("chatRoomId")
	topic := "chatroom." + id

	err = ch.ExchangeDeclare(
		topic,   // chatroom name
		"topic", // type
		true,    // durable
		true,    // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		topic,  // exchange
		false,
		nil,
	)

	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		// handle client connection lost.
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				forever <- true
				break
			}
		}
	}()

	go func() {

		for d := range msgs {
			// Send message to websocket corresponding to the chatroom

			ws.WriteMessage(websocket.TextMessage, d.Body)
		}
	}()

	<-forever

	log.Println("Connection closed")
	return nil
}
