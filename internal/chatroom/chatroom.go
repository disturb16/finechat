package chatroom

import (
	"context"

	"github.com/disturb16/finechat/internal/chatroom/models"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	// SaveChatRoom saves a new chat room and sets the user as its owner.
	SaveChatRoom(ctx context.Context, name string, userID int64) error
	// SaveChatRoomMessage saves the message in the chat room.
	SaveChatRoomMessage(ctx context.Context, chatRoomId int64, userId int64, message string) error
	// GetChatRooms returns all chat rooms for the user.
	GetChatRooms(ctx context.Context, userId int64) ([]models.ChatRoom, error)
	// GetChatRoomMessages returns the last 50 messages for the chat room.
	GetChatRoomMessages(ctx context.Context, chatRoomId int64) ([]models.ChatRoomMessage, error)
}

type Service interface {
	CreateChatRoom(ctx context.Context, name string, userID int64) error
	ListChatRooms(ctx context.Context, userId int64) ([]models.ChatRoom, error)
	ListChatRoomMessages(ctx context.Context, chatRoomId int64) ([]models.ChatRoomMessage, error)
	PostChatRoomMessage(ctx context.Context, chatRoomId int64, userId int64, message string) error
}

func NewRepository(db *sqlx.DB) Repository {
	return &ChatRoomRepository{db: db}
}

func NewService(repo Repository) Service {
	return &ChatRoomService{repo: repo}
}
