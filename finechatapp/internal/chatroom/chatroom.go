package chatroom

import (
	"context"
	"time"

	"github.com/disturb16/finechat/broker"
	"github.com/disturb16/finechat/internal/chatroom/models"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

// Repository is the chatroom interactor with the database.
type Repository interface {
	// SaveChatRoom saves a new chat room and sets the user as its owner.
	SaveChatRoom(ctx context.Context, name string, userID int64) error
	// SaveChatRoomMessage saves the message in the chat room.
	SaveChatRoomMessage(ctx context.Context, chatRoomId int64, email string, message string, createdDate time.Time) error
	// GetChatRooms returns all chat rooms for the user.
	GetChatRooms(ctx context.Context, email string) ([]models.ChatRoom, error)
	// GetChatRoomMessages returns the last 50 messages for the chat room.
	GetChatRoomMessages(ctx context.Context, chatRoomId int64) ([]models.ChatRoomMessage, error)
	// SaveChatRoomUser saves the user as a member of the chat room.
	SaveChatRoomUser(ctx context.Context, chatRoomId int64, userId int64) error
	// GetChatRoomUsers returns all users in the chat room.
	GetChatRoomUsers(ctx context.Context, chatRoomId int64) ([]models.ChatRoomUser, error)
	// RemoveChatRoomUser deletes the user from the chat room.
	RemoveChatRoomGuest(ctx context.Context, chatRoomId int64, email string) error
}

// Service is the chatroom business logic layer.
type Service interface {
	// CreateChatRoom creates a new chat room.
	CreateChatRoom(ctx context.Context, name string, userID int64) error
	// ListChatRooms returns all chat rooms for the user.
	ListChatRooms(ctx context.Context, email string) ([]models.ChatRoom, error)
	// ListChatRoomMessages returns the last 50 messages for the chat room.
	ListChatRoomMessages(ctx context.Context, chatRoomId int64) ([]models.ChatRoomMessage, error)
	// PostChatRoomMessage saves the message in the chat room or process the chat command.
	PostChatRoomMessage(ctx context.Context, chatRoomId int64, email string, message string, createdDate time.Time) error
	// AddChatRoomGuest adds a guest to the chat room.
	AddChatRoomGuest(ctx context.Context, chatRoomId int64, userId int64) error
	// ListChatRoomGuests returns all guests in the chat room.
	ListChatRoomGuests(ctx context.Context, chatRoomId int64) ([]models.ChatRoomUser, error)
	// RemoveChatRoomGuest deletes the guest from the chat room.
	RemoveChatRoomGuest(ctx context.Context, chatRoomId int64, email string) error
	// SubscribeToChatroomSocket subscribes to the chatroom websocket.
	SubscribeToChatroomSocket(ws *websocket.Conn, chatRoomId, email string) error
}

// NewRepository returns a new chatroom repository.
func NewRepository(db *sqlx.DB) Repository {
	return &ChatRoomRepository{db: db}
}

// NewService returns a new chatroom service.
func NewService(repo Repository, b broker.MessageBroker) Service {
	return &ChatRoomService{
		repo:   repo,
		broker: b,
	}
}
