package chatroom

import (
	"context"

	"github.com/disturb16/finechat/internal/chatroom/models"
)

type ChatRoomService struct {
	repo Repository
}

func (s *ChatRoomService) CreateChatRoom(ctx context.Context, name string, userID int64) error {
	return s.repo.SaveChatRoom(ctx, name, userID)
}

func (s *ChatRoomService) ListChatRooms(ctx context.Context, userId int64) ([]models.ChatRoom, error) {
	return s.repo.GetChatRooms(ctx, userId)
}

func (s *ChatRoomService) ListChatRoomMessages(ctx context.Context, chatRoomId int64) ([]models.ChatRoomMessage, error) {
	return s.repo.GetChatRoomMessages(ctx, chatRoomId)
}

func (s *ChatRoomService) PostChatRoomMessage(ctx context.Context, chatRoomId int64, userId int64, message string) error {
	return s.repo.SaveChatRoomMessage(ctx, chatRoomId, userId, message)
}
