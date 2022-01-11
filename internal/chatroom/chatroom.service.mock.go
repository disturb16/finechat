package chatroom

import (
	"context"
	"time"

	"github.com/disturb16/finechat/internal/chatroom/models"
	"github.com/gorilla/websocket"
)

type MockChatRoomService struct{}

func (mcs *MockChatRoomService) CreateChatRoom(ctx context.Context, name string, userID int64) error {
	return nil
}

func (mcs *MockChatRoomService) ListChatRooms(ctx context.Context, email string) ([]models.ChatRoom, error) {
	return nil, nil
}

func (mcs *MockChatRoomService) ListChatRoomMessages(ctx context.Context, chatRoomId int64) ([]models.ChatRoomMessage, error) {
	return nil, nil
}

func (mcs *MockChatRoomService) PostChatRoomMessage(ctx context.Context, chatRoomId int64, email string, message string, createdDate time.Time) error {
	return nil
}

func (mcs *MockChatRoomService) AddChatRoomGuest(ctx context.Context, chatRoomId int64, userId int64) error {
	return nil
}

func (mcs *MockChatRoomService) ListChatRoomGuests(ctx context.Context, chatRoomId int64) ([]models.ChatRoomUser, error) {
	return nil, nil
}

func (mcs *MockChatRoomService) RemoveChatRoomGuest(ctx context.Context, chatRoomId int64, email string) error {
	return nil
}

func (mcs *MockChatRoomService) SubscribeToChatroomSocket(ws *websocket.Conn, chatRoomId, email string) error {
	return nil
}
