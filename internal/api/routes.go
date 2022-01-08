package api

import (
	"log"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(h *Handler, e *echo.Echo) {
	e.Use(RequestLogger)
	e.Use(EnrichContext)

	usersAPI := e.Group("/api/users")
	usersAPI.POST("", h.RegisterUser)
	usersAPI.POST("/signin", h.Signin)
	usersAPI.GET("/:email/chatrooms", h.chatRoomsByUser)

	chatroomsAPI := e.Group("/api/chatrooms")
	chatroomsAPI.POST("", h.createChatRoom)

	for _, r := range e.Routes() {
		log.Printf("[%s] %s", r.Method, r.Path)
	}
}
