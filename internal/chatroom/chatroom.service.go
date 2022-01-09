package chatroom

import (
	"context"
	"time"

	"github.com/disturb16/finechat/internal/chatroom/models"
)

type ChatRoomService struct {
	repo Repository
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

func (s *ChatRoomService) PostChatRoomMessage(ctx context.Context, chatRoomId int64, userId int64, message string, createdDate time.Time) error {
	return s.repo.SaveChatRoomMessage(ctx, chatRoomId, userId, message, createdDate)
}

func (s *ChatRoomService) AddChatRoomGuest(ctx context.Context, chatRoomId int64, userId int64) error {
	return s.repo.SaveChatRoomUser(ctx, chatRoomId, userId)
}

func (s *ChatRoomService) ListChatRoomGuests(ctx context.Context, chatRoomId int64) ([]models.ChatRoomUser, error) {
	return s.repo.GetChatRoomUsers(ctx, chatRoomId)
}
