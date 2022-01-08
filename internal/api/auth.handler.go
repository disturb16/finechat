package api

import (
	"net/http"

	"github.com/disturb16/finechat/internal/api/dtos"
	"github.com/labstack/echo/v4"
)

type messageResponse struct {
	Message string `json:"message"`
}

type AuthResponse struct {
	Token string `json:"token"`
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
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.JSON(http.StatusOK, &AuthResponse{Token: token})
}
