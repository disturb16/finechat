package api

import (
	"net/http"

	"github.com/disturb16/finechat/internal/api/dtos"
	"github.com/disturb16/finechat/internal/auth"
	"github.com/disturb16/finechat/logger"
	"github.com/labstack/echo/v4"
)

type messageResponse struct {
	Message string `json:"message"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type FriendResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *Handler) RegisterUser(c echo.Context) error {
	ctx := c.Request().Context()
	params := &dtos.RegisterUser{}

	err := decodeBody(c, params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest)
	}

	err = h.validate.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, parseDTOError(err))
	}

	err = h.authService.RegisterUser(ctx, params.FirstName, params.LastName, params.Email, params.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.JSON(http.StatusCreated, &messageResponse{"user created"})
}

func (h *Handler) Signin(c echo.Context) error {
	ctx := c.Request().Context()
	params := &dtos.Sigin{}

	err := decodeBody(c, params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest)
	}

	err = h.validate.Struct(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, parseDTOError(err))
	}

	token, err := h.authService.LoginUser(ctx, params.Email, params.Password)
	if err != nil {
		if err == auth.ErrUserNotFound || err == auth.ErrInvalidUserCredentials {
			return c.JSON(http.StatusBadRequest, auth.ErrUserNotFound)
		}

		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.JSON(http.StatusOK, &AuthResponse{Token: token})
}

func (h *Handler) addFriend(c echo.Context) error {
	ctx := c.Request().Context()
	params := &dtos.AddFriend{
		Email: c.Param("email"),
	}

	err := decodeBody(c, params)
	if err != nil {
		logger.Println(ctx, err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest)
	}

	err = h.validate.Struct(params)
	if err != nil {
		logger.Println(ctx, err)
		return c.JSON(http.StatusBadRequest, parseDTOError(err))
	}

	err = h.authService.AddFriend(ctx, params.Email, params.FriendEmail)
	if err != nil {
		logger.Println(ctx, err)
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.JSON(http.StatusCreated, &messageResponse{"friend added"})
}

func (h *Handler) userFriends(c echo.Context) error {
	ctx := c.Request().Context()
	email := c.Param("email")

	friends, err := h.authService.ListFriends(ctx, email)
	if err != nil {
		logger.Println(ctx, err)
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	friendList := []FriendResponse{}

	for _, friend := range friends {
		friendList = append(friendList, FriendResponse{
			Email: friend.Email,
			Name:  friend.FirstName + " " + friend.LastName,
		})
	}

	return c.JSON(http.StatusOK, friendList)
}
