package api

import (
	"log"
	"net/http"

	"github.com/disturb16/finechat/broker"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func (h *Handler) subscribeToChatroomSocket(c echo.Context) error {
	var err error
	defer func() {
		log.Println(err)
	}()

	// Upgrade the connection to a WebSocket.
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	log.Println("New connection")

	// Create a channel
	ch, err := h.messageBroker.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	// Create a queue for this channel
	id := c.Param("chatRoomId")
	email := c.Param("email")

	emailKey := "." + email
	topic := "chatroom." + id

	err = broker.DefaultExchange(ch, topic)
	if err != nil {
		return err
	}

	q, err := broker.DefaultQueue(ch, "")
	if err != nil {
		return err
	}

	// Binding for all messages of the chatroom
	err = ch.QueueBind(
		q.Name, // queue name
		topic,  // routing key
		topic,  // exchange
		false,
		nil,
	)

	if err != nil {
		return err
	}

	defer ch.QueueUnbind(q.Name, topic, topic, nil)

	// Binding for messages only for this user
	err = ch.QueueBind(
		q.Name,         // queue name
		topic+emailKey, // routing key
		topic,          // exchange
		false,
		nil,
	)

	if err != nil {
		return err
	}

	defer ch.QueueUnbind(q.Name, topic+emailKey, topic, nil)

	msgs, err := broker.DefaultConsumer(ch, q)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// Send message to websocket corresponding to the chatroom
			ws.WriteMessage(websocket.TextMessage, d.Body)
		}
	}()

	go func() {
		for {
			// handle when client connection is lost.
			_, _, err := ws.ReadMessage()
			if err != nil {
				forever <- true
				break
			}
		}
	}()

	<-forever

	log.Println("Connection closed")
	return nil
}
