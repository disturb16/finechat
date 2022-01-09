package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/disturb16/finechat/broker"
	"github.com/disturb16/finechat/internal/api/dtos"
	"github.com/disturb16/finechat/logger"
	"github.com/labstack/echo/v4"
)

type ChatRoomMessageResponse struct {
	Message     string    `json:"message"`
	User        string    `json:"user"`
	CreatedDate time.Time `json:"created_date"`
}

func (h *Handler) chatRoomsByUser(c echo.Context) error {
	ctx := c.Request().Context()

	params := &dtos.ChatRoomsByUser{
		Email: c.Param("email"),
	}

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

func (h *Handler) createChatRoomMessage(c echo.Context) error {
	ctx := c.Request().Context()
	params := &dtos.ChatRoomMessage{}

	chatRoomIDStr := c.Param("chatRoomId")
	chatRoomID, err := strconv.ParseInt(chatRoomIDStr, 10, 64)
	if err != nil {
		logger.Println(ctx, "error parsing chatroom id", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest)
	}

	err = decodeBody(c, params)
	if err != nil {
		logger.Println(ctx, err)
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

	// Save message.
	err = h.chatRoomService.PostChatRoomMessage(ctx, chatRoomID, user.ID, params.Message, params.CreatedDate)
	if err != nil {
		logger.Println(ctx, "error posting chatroom message", err)
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	exchange := "chatroom." + chatRoomIDStr
	err = h.messageBroker.SendMessage(exchange, exchange, broker.TypeReload, "")
	if err != nil {
		logger.Println(ctx, "error sending message to broker", err)
	}

	return c.JSON(http.StatusOK, &messageResponse{"message posted"})
}

func (h *Handler) chatRoomMessages(c echo.Context) error {
	ctx := c.Request().Context()

	chatRoomIDStr := c.Param("chatRoomId")
	chatRoomID, err := strconv.ParseInt(chatRoomIDStr, 10, 64)
	if err != nil {
		logger.Println(ctx, "error parsing chatroom id", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest)
	}

	mm, err := h.chatRoomService.ListChatRoomMessages(ctx, chatRoomID)
	if err != nil {
		logger.Println(ctx, "error listing chatroom messages", err)
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	messages := []ChatRoomMessageResponse{}
	for _, m := range mm {
		messages = append(messages, ChatRoomMessageResponse{
			Message:     m.Message,
			User:        m.User,
			CreatedDate: m.CreatedDate,
		})
	}

	return c.JSON(http.StatusOK, messages)
}

func (h *Handler) createChatRoomUser(c echo.Context) error {
	ctx := c.Request().Context()

	chatRoomIDStr := c.Param("chatRoomId")
	chatRoomID, err := strconv.ParseInt(chatRoomIDStr, 10, 64)
	if err != nil {
		logger.Println(ctx, "error parsing chatroom id", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest)
	}

	params := &dtos.CreateChatRoomUser{}

	err = decodeBody(c, params)
	if err != nil {
		logger.Println(ctx, err)
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
	err = h.chatRoomService.AddChatRoomGuest(ctx, chatRoomID, user.ID)

	return c.JSON(http.StatusCreated, nil)
}

func (h *Handler) chatRoomUsers(c echo.Context) error {
	ctx := c.Request().Context()

	chatRoomIDStr := c.Param("chatRoomId")
	chatRoomID, err := strconv.ParseInt(chatRoomIDStr, 10, 64)
	if err != nil {
		logger.Println(ctx, "error parsing chatroom id", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest)
	}

	users, err := h.chatRoomService.ListChatRoomGuests(ctx, chatRoomID)
	if err != nil {
		logger.Println(ctx, "error listing chatroom users", err)
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.JSON(http.StatusOK, users)
}

func (h *Handler) removeChatRoomUser(c echo.Context) error {
	ctx := c.Request().Context()

	params := &dtos.RemoveChatRoomUser{
		Email: c.Param("email"),
	}

	chatRoomIDStr := c.Param("chatRoomId")
	chatRoomID, err := strconv.ParseInt(chatRoomIDStr, 10, 64)
	if err != nil {
		logger.Println(ctx, "error parsing chatroom id", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest)
	}

	err = h.validate.Struct(params)
	if err != nil {
		logger.Println(ctx, err)
		return c.JSON(http.StatusBadRequest, parseDTOError(err))
	}

	err = h.chatRoomService.RemoveChatRoomGuest(ctx, chatRoomID, params.Email)
	if err != nil {
		logger.Println(ctx, "error removing chatroom user", err)
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.NoContent(http.StatusOK)
}
