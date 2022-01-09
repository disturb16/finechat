package models

import "time"

type ChatRoomMessage struct {
	Message     string    `db:"message"`
	UserID      int64     `db:"user_id"`
	User        string    `db:"user"`
	CreatedDate time.Time `db:"created_date"`
}
