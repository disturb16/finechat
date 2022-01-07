package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) RegisterUser(c echo.Context) error {
	return c.JSON(http.StatusCreated, nil)
}
