package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidRequest = echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	ErrInternalServer = echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	ErrBeerTaken      = echo.NewHTTPError(http.StatusBadRequest, "beer's id already taken")
	ErrBeerNotFound   = echo.NewHTTPError(http.StatusNotFound, "beer not found")
)
