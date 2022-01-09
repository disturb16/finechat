package models

type ChatRoomUser struct {
	Email string `json:"email" db:"email"`
	Name  string `json:"name" db:"name"`
}
