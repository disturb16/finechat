package chatroom

import (
	"context"
	"time"

	"github.com/disturb16/finechat/internal/chatroom/models"
)

type MockRepository struct{}

func (mr *MockRepository) SaveChatRoom(ctx context.Context, name string, userID int64) error {
	return nil
}

func (mr *MockRepository) SaveChatRoomMessage(ctx context.Context, chatRoomId int64, email string, message string, createdDate time.Time) error {
	return nil
}

func (mr *MockRepository) GetChatRooms(ctx context.Context, email string) ([]models.ChatRoom, error) {
	return nil, nil
}

func (mr *MockRepository) GetChatRoomMessages(ctx context.Context, chatRoomId int64) ([]models.ChatRoomMessage, error) {
	return nil, nil
}

func (mr *MockRepository) SaveChatRoomUser(ctx context.Context, chatRoomId int64, userId int64) error {
	return nil
}

func (mr *MockRepository) GetChatRoomUsers(ctx context.Context, chatRoomId int64) ([]models.ChatRoomUser, error) {
	return nil, nil
}

func (mr *MockRepository) RemoveChatRoomGuest(ctx context.Context, chatRoomId int64, email string) error {
	return nil
}
