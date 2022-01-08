package api

import (
	"log"
	"net/http"

	"github.com/disturb16/finechat/internal/api/dtos"
	"github.com/labstack/echo/v4"
)

func (h *Handler) chatRoomsByUser(c echo.Context) error {
	ctx := c.Request().Context()

	params := &dtos.ChatRoomsByUser{
		Email: c.Param("email"),
	}

	log.Println(params)

	err := h.validate.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, parseDTOError(err))
	}

	chatRooms, err := h.chatRoomService.ListChatRooms(ctx, params.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.JSON(http.StatusOK, chatRooms)
}

func (h *Handler) createChatRoom(c echo.Context) error {
	ctx := c.Request().Context()
	params := &dtos.CreateChatRoom{}

	err := decodeBody(c, params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest)
	}

	err = h.validate.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, parseDTOError(err))
	}

	// Get user information.
	user, err := h.authService.FindUserByEmail(ctx, params.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	// Save chatroom.
	err = h.chatRoomService.CreateChatRoom(ctx, params.Name, user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.JSON(http.StatusCreated, &messageResponse{"chatroom created"})
}
