package dtos

type ChatRoomsByUser struct {
	Email string `json:"email" validate:"required"`
}
