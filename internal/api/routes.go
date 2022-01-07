package api

import (
	"log"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(h *Handler, e *echo.Echo) {
	e.Use(RequestLogger)
	e.Use(EnrichContext)

	apiGroup := e.Group("/api")

	apiGroup.POST("/users", h.RegisterUser)

	for _, r := range e.Routes() {
		log.Printf("[%s] %s", r.Method, r.Path)
	}
}
