package models

type ChatRoomMessage struct {
	Message  string `db:"message"`
	UserID   int64  `db:"user_id"`
	UserName string `db:"user_name"`
}
