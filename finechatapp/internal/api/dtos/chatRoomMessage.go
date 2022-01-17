package dtos

import "time"

type ChatRoomMessage struct {
	Email       string    `json:"email" validate:"required,email"`
	Message     string    `json:"message" validate:"required"`
	CreatedDate time.Time `json:"created_date"`
}
