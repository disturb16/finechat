package chatroom

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/disturb16/finechat/broker"
	"github.com/disturb16/finechat/internal/chatroom/models"
	"github.com/disturb16/finechat/internal/finechatbot"

	"github.com/gorilla/websocket"
)

type ChatRoomService struct {
	repo   Repository
	broker *broker.Broker
}

func (s *ChatRoomService) CreateChatRoom(ctx context.Context, name string, userID int64) error {
	return s.repo.SaveChatRoom(ctx, name, userID)
}

func (s *ChatRoomService) ListChatRooms(ctx context.Context, email string) ([]models.ChatRoom, error) {
	return s.repo.GetChatRooms(ctx, email)
}

func (s *ChatRoomService) ListChatRoomMessages(ctx context.Context, chatRoomId int64) ([]models.ChatRoomMessage, error) {
	return s.repo.GetChatRoomMessages(ctx, chatRoomId)
}

func (s *ChatRoomService) PostChatRoomMessage(ctx context.Context, chatRoomId int64, email string, message string, createdDate time.Time) error {

	// If the message starts with / then it is a command
	// else is a normal message.
	if strings.HasPrefix(message, "/") {
		payload := &finechatbot.StockCommand{
			Email:      email,
			ChatRoomID: chatRoomId,
			Message:    message,
		}

		return s.broker.SendMessage(
			finechatbot.StockCommandTopic,
			"",
			broker.TypeStockRequest,
			payload,
		)
	}

	// Save message to database.
	err := s.repo.SaveChatRoomMessage(ctx, chatRoomId, email, message, createdDate)
	if err != nil {
		return err
	}

	// Send message to notify other users in the chatroom.
	exchange := fmt.Sprintf("chatroom.%d", chatRoomId)
	s.broker.SendMessage(exchange, exchange, broker.TypeReload, "")
	return nil
}

func (s *ChatRoomService) AddChatRoomGuest(ctx context.Context, chatRoomId int64, userId int64) error {
	return s.repo.SaveChatRoomUser(ctx, chatRoomId, userId)
}

func (s *ChatRoomService) ListChatRoomGuests(ctx context.Context, chatRoomId int64) ([]models.ChatRoomUser, error) {
	return s.repo.GetChatRoomUsers(ctx, chatRoomId)
}

func (s *ChatRoomService) RemoveChatRoomGuest(ctx context.Context, chatRoomId int64, email string) error {
	return s.repo.RemoveChatRoomGuest(ctx, chatRoomId, email)
}

func (s *ChatRoomService) SubscribeToChatroomSocket(ws *websocket.Conn, chatRoomId, email string) error {
	// Create a channel
	ch, err := s.broker.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	topic := "chatroom." + chatRoomId

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

	// Binding for messages of the user.
	key := topic + "." + email
	err = ch.QueueBind(
		q.Name, // queue name
		key,    // routing key
		topic,  // exchange
		false,
		nil,
	)

	if err != nil {
		return err
	}

	defer ch.QueueUnbind(q.Name, key, topic, nil)

	msgs, err := broker.DefaultConsumer(ch, q)
	if err != nil {
		return err
	}

	taskDone := make(chan bool)

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
				taskDone <- true
				break
			}
		}
	}()

	<-taskDone
	return nil
}
