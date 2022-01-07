package api

import (
	"github.com/disturb16/finechat/internal/auth"
	"github.com/disturb16/finechat/internal/chatroom"
	"gopkg.in/go-playground/validator.v9"
)

type Handler struct {
	authService     auth.Service
	chatRoomService chatroom.Service
	validate        *validator.Validate
}

// New creates a new api handler
func NewHandler(authService auth.Service, chatRoomService chatroom.Service) *Handler {

	v := validator.New()

	handler := &Handler{
		authService:     authService,
		chatRoomService: chatRoomService,
		validate:        v,
	}

	return handler
}
