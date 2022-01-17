package dtos

type RemoveChatRoomUser struct {
	Email string `json:"email" validate:"required,email"`
}
