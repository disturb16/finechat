package chatroom

import (
	"context"
	"fmt"
	"time"

	"github.com/disturb16/finechat/broker"
	"github.com/disturb16/finechat/internal/chatroom/models"
	"github.com/disturb16/finechat/stockhelper"
	"github.com/gorilla/websocket"
)

type ChatRoomService struct {
	repo   Repository
	broker *broker.Broker
}

func messageIsCommand(message string) bool {
	return message[0] == '/'
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

	exchange := fmt.Sprintf("chatroom.%d", chatRoomId)

	if messageIsCommand(message) {
		return s.processStockCommand(exchange, chatRoomId, email, message)
	}

	err := s.repo.SaveChatRoomMessage(ctx, chatRoomId, email, message, createdDate)
	if err != nil {
		return err
	}

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

	emailKey := "." + email
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

func (s *ChatRoomService) processStockCommand(exchange string, chatRoomId int64, email, message string) error {
	symbol, err := stockhelper.GetSymbol(message)
	if err != nil {
		return err
	}

	stockShare, err := stockhelper.GetShare(symbol)
	if err != nil {
		return err
	}

	// key := fmt.Sprintf("chatroom.%d.%s", chatRoomId, email)
	payload := fmt.Sprintf("%s quote is $%s per share", symbol, stockShare)

	return s.broker.SendMessage(exchange, exchange, broker.TypeStockRequest, payload)
}
