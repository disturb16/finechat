package api

import (
	"log"
	"net/http"

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

	return h.chatRoomService.SubscribeToChatroomSocket(ws, id, email)
}
