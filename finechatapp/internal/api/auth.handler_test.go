package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/disturb16/finechat/broker"
	"github.com/disturb16/finechat/internal/api/dtos"
	"github.com/disturb16/finechat/internal/auth"
	"github.com/disturb16/finechat/internal/chatroom"
	"github.com/labstack/echo/v4"
)

var h *Handler

func TestMain(m *testing.M) {
	h = NewHandler(
		&auth.MockAuthService{},
		&chatroom.MockChatRoomService{},
		&broker.MockBroker{},
	)

	code := m.Run()
	os.Exit(code)
}

func TestRegisterUser(t *testing.T) {
	testCases := []struct {
		Name               string
		Input              *dtos.RegisterUser
		StatusCodeExpected int
	}{
		{
			Name: "register user",
			Input: &dtos.RegisterUser{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@doe.com",
				Password:  "password",
			},
			StatusCodeExpected: http.StatusCreated,
		},
		{
			Name: "register user with invalid email",
			Input: &dtos.RegisterUser{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@doe",
			},
			StatusCodeExpected: http.StatusBadRequest,
		},
		{
			Name:               "register user with missing fields",
			Input:              &dtos.RegisterUser{},
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
			err := h.RegisterUser(c)
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

func TestSignin(t *testing.T) {
	testCases := []struct {
		Name               string
		Input              *dtos.Sigin
		StatusCodeExpected int
	}{
		{
			Name: "signin user",
			Input: &dtos.Sigin{
				Email:    "john@doe.com",
				Password: "password",
			},
			StatusCodeExpected: http.StatusOK,
		},
		{
			Name: "signin user with invalid email",
			Input: &dtos.Sigin{
				Email:    "john@doe",
				Password: "password",
			},
			StatusCodeExpected: http.StatusBadRequest,
		},
		{
			Name:               "signin user with empty input",
			Input:              &dtos.Sigin{},
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
			err := h.Signin(c)
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
