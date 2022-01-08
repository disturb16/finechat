package chatroom

import (
	"context"

	"github.com/disturb16/finechat/internal/chatroom/models"
	"github.com/jmoiron/sqlx"
)

type ChatRoomRepository struct {
	db *sqlx.DB
}

const (
	qrySaveChatRoom        string = `call saveChatRoom(?,?)`
	qryGetChatRoomsByUser  string = `call getChatRoomsByUser(?)`
	qryChatRoomMessages    string = `call getChatRoomMessages(?)`
	qrySaveChatRoomMessage string = `call saveChatRoomMessage(?, ?, ?)`
	qrySaveChatRoomUser    string = `call saveChatRoomUser(?, ?)`
)

func (r *ChatRoomRepository) SaveChatRoom(ctx context.Context, name string, userID int64) error {
	_, err := r.db.ExecContext(ctx, qrySaveChatRoom, name, userID)
	return err
}

func (r *ChatRoomRepository) GetChatRooms(ctx context.Context, email string) ([]models.ChatRoom, error) {
	chatRooms := []models.ChatRoom{}
	err := r.db.SelectContext(ctx, &chatRooms, qryGetChatRoomsByUser, email)
	return chatRooms, err
}

func (r *ChatRoomRepository) GetChatRoomMessages(ctx context.Context, chatRoomId int64) ([]models.ChatRoomMessage, error) {
	messages := []models.ChatRoomMessage{}
	err := r.db.SelectContext(ctx, &messages, qryChatRoomMessages, chatRoomId)
	return messages, err
}

func (r *ChatRoomRepository) SaveChatRoomMessage(ctx context.Context, chatRoomId int64, userId int64, message string) error {
	_, err := r.db.ExecContext(ctx, qrySaveChatRoomMessage, chatRoomId, userId, message)
	return err
}

func (r *ChatRoomRepository) SaveChatRoomUser(ctx context.Context, chatRoomId int64, userId int64) error {
	_, err := r.db.ExecContext(ctx, qrySaveChatRoomUser, chatRoomId, userId)
	return err
}
