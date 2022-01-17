package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/disturb16/finechat/internal/api/dtos"
	"github.com/labstack/echo/v4"
)

func TestChatRoomsByUser(t *testing.T) {
	testCases := []struct {
		Name               string
		Email              string
		StatusCodeExpected int
	}{
		{
			Name:               "Get chat rooms by user",
			Email:              "john@doe.com",
			StatusCodeExpected: http.StatusOK,
		},
		{
			Name:               "Get chat rooms by user with invalid email",
			Email:              "john@doe",
			StatusCodeExpected: http.StatusBadRequest,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(r, w)
			c.SetParamNames("email")
			c.SetParamValues(tc.Email)

			// invoke handler
			err := h.chatRoomsByUser(c)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// check status code
			if w.Code != tc.StatusCodeExpected {
				t.Errorf("expected status code %d, got %d", tc.StatusCodeExpected, w.Code)
			}
		})
	}
}

func TestCreateChatRoom(t *testing.T) {
	testCases := []struct {
		Name               string
		Input              *dtos.CreateChatRoom
		StatusCodeExpected int
	}{
		{
			Name: "Create chat room",
			Input: &dtos.CreateChatRoom{
				Name:  "chat room",
				Email: "john@doe.com",
			},
			StatusCodeExpected: http.StatusCreated,
		},
		{
			Name: "Create chat room with invalid email",
			Input: &dtos.CreateChatRoom{
				Name:  "chat room",
				Email: "john@doe",
			},
			StatusCodeExpected: http.StatusBadRequest,
		},
		{
			Name:               "Create chat room with missing fields",
			Input:              &dtos.CreateChatRoom{},
			StatusCodeExpected: http.StatusBadRequest,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			body, _ := json.Marshal(tc.Input)
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
			w := httptest.NewRecorder()

			// set echo context
			e := echo.New()
			c := e.NewContext(r, w)

			// invoke handler
			err := h.createChatRoom(c)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// check status code
			if w.Code != tc.StatusCodeExpected {
				t.Errorf("expected status code %d, got %d", tc.StatusCodeExpected, w.Code)
			}
		})
	}
}

func TestCreateChatRoomMessage(t *testing.T) {
	testCases := []struct {
		Name               string
		ChatRoomID         string
		Input              *dtos.ChatRoomMessage
		StatusCodeExpected int
	}{
		{
			Name:       "Create chat room message",
			ChatRoomID: "1",
			Input: &dtos.ChatRoomMessage{
				Email:   "john@does.com",
				Message: "Hello",
			},
			StatusCodeExpected: http.StatusOK,
		},
		{
			Name:       "Create chat room command message",
			ChatRoomID: "1",
			Input: &dtos.ChatRoomMessage{
				Email:   "john@doe.com",
				Message: "/help",
			},
			StatusCodeExpected: http.StatusOK,
		},
		{
			Name:       "Create chat room message with invalid email",
			ChatRoomID: "1",
			Input: &dtos.ChatRoomMessage{
				Email:   "john@doe",
				Message: "Hello",
			},
			StatusCodeExpected: http.StatusBadRequest,
		},
		{
			Name:               "Create chat room message with invalid chat room id",
			ChatRoomID:         "s",
			Input:              &dtos.ChatRoomMessage{},
			StatusCodeExpected: http.StatusBadRequest,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			body, _ := json.Marshal(tc.Input)
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
			w := httptest.NewRecorder()

			// set echo context
			e := echo.New()
			c := e.NewContext(r, w)
			c.SetParamNames("chatRoomId")
			c.SetParamValues(tc.ChatRoomID)

			// invoke handler
			err := h.createChatRoomMessage(c)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// check status code
			if w.Code != tc.StatusCodeExpected {
				t.Errorf("expected status code %d, got %d", tc.StatusCodeExpected, w.Code)
			}
		})
	}
}

func TestChatRoomMessages(t *testing.T) {
	testCases := []struct {
		Name               string
		ChatRoomID         string
		StatusCodeExpected int
	}{
		{
			Name:               "Get chat room messages",
			ChatRoomID:         "1",
			StatusCodeExpected: http.StatusOK,
		},
		{
			Name:               "Get chat room messages with invalid chat room id",
			ChatRoomID:         "s",
			StatusCodeExpected: http.StatusBadRequest,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			// set echo context
			e := echo.New()
			c := e.NewContext(r, w)
			c.SetParamNames("chatRoomId")
			c.SetParamValues(tc.ChatRoomID)

			// invoke handler
			err := h.chatRoomMessages(c)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// check status code
			if w.Code != tc.StatusCodeExpected {
				t.Errorf("expected status code %d, got %d", tc.StatusCodeExpected, w.Code)
			}
		})
	}
}

func TestCreateChatRoomUser(t *testing.T) {
	testCases := []struct {
		Name               string
		ChatRoomID         string
		Input              *dtos.CreateChatRoomUser
		StatusCodeExpected int
	}{
		{
			Name:       "Create chat room user",
			ChatRoomID: "1",
			Input: &dtos.CreateChatRoomUser{
				Email: "john@doe.com",
			},
			StatusCodeExpected: http.StatusCreated,
		},
		{
			Name:       "Create chat room user with invalid email",
			ChatRoomID: "1",
			Input: &dtos.CreateChatRoomUser{
				Email: "john@doe",
			},
			StatusCodeExpected: http.StatusBadRequest,
		},
		{
			Name:               "Create chat room user with invalid chat room id",
			ChatRoomID:         "s",
			Input:              &dtos.CreateChatRoomUser{},
			StatusCodeExpected: http.StatusBadRequest,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			body, _ := json.Marshal(tc.Input)
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
			w := httptest.NewRecorder()

			// set echo context
			e := echo.New()
			c := e.NewContext(r, w)
			c.SetParamNames("chatRoomId")
			c.SetParamValues(tc.ChatRoomID)

			// invoke handler
			err := h.createChatRoomUser(c)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// check status code
			if w.Code != tc.StatusCodeExpected {
				t.Errorf("expected status code %d, got %d", tc.StatusCodeExpected, w.Code)
			}
		})
	}
}

func TestChatRoomUsers(t *testing.T) {
	testCases := []struct {
		Name               string
		ChatRoomID         string
		StatusCodeExpected int
	}{
		{
			Name:               "Get chat room users",
			ChatRoomID:         "1",
			StatusCodeExpected: http.StatusOK,
		},
		{
			Name:               "Get chat room users with invalid chat room id",
			ChatRoomID:         "s",
			StatusCodeExpected: http.StatusBadRequest,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			// set echo context
			e := echo.New()
			c := e.NewContext(r, w)
			c.SetParamNames("chatRoomId")
			c.SetParamValues(tc.ChatRoomID)

			// invoke handler
			err := h.chatRoomUsers(c)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// check status code
			if w.Code != tc.StatusCodeExpected {
				t.Errorf("expected status code %d, got %d", tc.StatusCodeExpected, w.Code)
			}
		})
	}
}

func TestRemoveChatRoomUser(t *testing.T) {
	testCases := []struct {
		Name               string
		ChatRoomID         string
		Input              *dtos.RemoveChatRoomUser
		StatusCodeExpected int
	}{
		{
			Name:       "Remove chat room user",
			ChatRoomID: "1",
			Input: &dtos.RemoveChatRoomUser{
				Email: "john@doe.com",
			},
			StatusCodeExpected: http.StatusOK,
		},
		{
			Name:       "Remove chat room user with invalid email",
			ChatRoomID: "1",
			Input: &dtos.RemoveChatRoomUser{
				Email: "john@doe",
			},
			StatusCodeExpected: http.StatusBadRequest,
		},
		{
			Name:               "Remove chat room user with invalid chat room id",
			ChatRoomID:         "s",
			Input:              &dtos.RemoveChatRoomUser{},
			StatusCodeExpected: http.StatusBadRequest,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodDelete, "/", nil)
			w := httptest.NewRecorder()

			// set echo context
			e := echo.New()
			c := e.NewContext(r, w)

			c.SetParamNames("chatRoomId", "email")
			c.SetParamValues(tc.ChatRoomID, tc.Input.Email)

			// invoke handler
			err := h.removeChatRoomUser(c)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// check status code
			if w.Code != tc.StatusCodeExpected {
				t.Errorf("expected status code %d, got %d", tc.StatusCodeExpected, w.Code)
			}
		})
	}
}
