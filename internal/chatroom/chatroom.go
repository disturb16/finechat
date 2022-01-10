package chatroom

import (
	"context"
	"time"

	"github.com/disturb16/finechat/broker"
	"github.com/disturb16/finechat/internal/chatroom/models"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	// SaveChatRoom saves a new chat room and sets the user as its owner.
	SaveChatRoom(ctx context.Context, name string, userID int64) error
	// SaveChatRoomMessage saves the message in the chat room.
	SaveChatRoomMessage(ctx context.Context, chatRoomId int64, email string, message string, createdDate time.Time) error
	// GetChatRooms returns all chat rooms for the user.
	GetChatRooms(ctx context.Context, email string) ([]models.ChatRoom, error)
	// GetChatRoomMessages returns the last 50 messages for the chat room.
	GetChatRoomMessages(ctx context.Context, chatRoomId int64) ([]models.ChatRoomMessage, error)

	SaveChatRoomUser(ctx context.Context, chatRoomId int64, userId int64) error
	GetChatRoomUsers(ctx context.Context, chatRoomId int64) ([]models.ChatRoomUser, error)

	RemoveChatRoomGuest(ctx context.Context, chatRoomId int64, email string) error
}

type Service interface {
	CreateChatRoom(ctx context.Context, name string, userID int64) error
	ListChatRooms(ctx context.Context, email string) ([]models.ChatRoom, error)
	ListChatRoomMessages(ctx context.Context, chatRoomId int64) ([]models.ChatRoomMessage, error)
	PostChatRoomMessage(ctx context.Context, chatRoomId int64, email string, message string, createdDate time.Time) error
	AddChatRoomGuest(ctx context.Context, chatRoomId int64, userId int64) error
	ListChatRoomGuests(ctx context.Context, chatRoomId int64) ([]models.ChatRoomUser, error)
	RemoveChatRoomGuest(ctx context.Context, chatRoomId int64, email string) error
	SubscribeToChatroomSocket(ws *websocket.Conn, chatRoomId, email string) error
}

func NewRepository(db *sqlx.DB) Repository {
	return &ChatRoomRepository{db: db}
}

func NewService(repo Repository, b *broker.Broker) Service {
	return &ChatRoomService{
		repo:   repo,
		broker: b,
	}
}
