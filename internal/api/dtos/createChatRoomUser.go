package dtos

type CreateChatRoomUser struct {
	Email string `json:"email" validate:"required,email"`
}
