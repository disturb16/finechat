package api

import (
	"log"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(h *Handler, e *echo.Echo) {
	e.Use(RequestLogger)
	e.Use(EnrichContext)

	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	// ======= USERS ENDPOINTS =======
	usersAPI := e.Group("/api/users")
	usersAPI.POST("", h.RegisterUser)
	usersAPI.POST("/signin", h.Signin)
	usersAPI.GET("/:email/chatrooms", h.chatRoomsByUser)
	usersAPI.POST("/:email/friends", h.addFriend)
	usersAPI.GET("/:email/friends", h.userFriends)

	// ======= CHATROOM ENDPIONTS ===========
	chatroomsAPI := e.Group("/api/chatrooms")
	chatroomsAPI.POST("", h.createChatRoom)
	chatroomsAPI.POST("/:chatRoomId/messages", h.createChatRoomMessage)
	chatroomsAPI.GET("/:chatRoomId/messages", h.chatRoomMessages)
	chatroomsAPI.POST("/:chatRoomId/users", h.createChatRoomUser)
	chatroomsAPI.GET("/:chatRoomId/users", h.chatRoomUsers)

	for _, r := range e.Routes() {
		log.Printf("[%s] %s", r.Method, r.Path)
	}
}
