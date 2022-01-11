package chatroom

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/disturb16/finechat/broker"
)

var s Service

func TestMain(m *testing.M) {
	s = NewService(&MockRepository{}, &broker.MockBroker{})
	code := m.Run()
	os.Exit(code)
}

func TestCreateChatRoom(t *testing.T) {
	testCases := []struct {
		Name          string
		ChatRoomName  string
		UserID        int64
		ExpectedError error
	}{
		{
			Name:          "Create chat room",
			ChatRoomName:  "Test chat room",
			UserID:        1,
			ExpectedError: nil,
		},
	}

	ctx := context.Background()

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			err := s.CreateChatRoom(ctx, tc.ChatRoomName, tc.UserID)
			if err != tc.ExpectedError {
				t.Errorf("Expected error %v, got %v", tc.ExpectedError, err)
			}
		})
	}
}

func TestPostChatRoomMessage(t *testing.T) {
	testCases := []struct {
		Name          string
		ChatRoomID    int64
		Email         string
		Message       string
		ExpectedError error
	}{
		{
			Name:          "Post chat room message",
			ChatRoomID:    1,
			Email:         "john@doe.com",
			Message:       "Hello world",
			ExpectedError: nil,
		},
		{
			Name:          "Post chat command",
			ChatRoomID:    1,
			Email:         "john@doe.com",
			Message:       "/stock=aapl.us",
			ExpectedError: nil,
		},
	}

	ctx := context.Background()

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			err := s.PostChatRoomMessage(ctx, tc.ChatRoomID, tc.Email, tc.Message, time.Now())
			if err != tc.ExpectedError {
				t.Errorf("Expected error %v, got %v", tc.ExpectedError, err)
			}
		})
	}
}
