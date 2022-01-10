package auth

import (
	"context"
	"os"
	"testing"
)

var s Service

func TestMain(m *testing.M) {
	s = NewService(&MockRepository{})

	code := m.Run()
	os.Exit(code)
}

func TestRegisterUser(t *testing.T) {
	testCases := []struct {
		Name          string
		FirstName     string
		LastName      string
		Email         string
		Password      string
		ExpectedError error
	}{
		{
			Name:          "Register a new user",
			FirstName:     "John",
			LastName:      "Doe",
			Email:         "john@doe.com",
			Password:      "password",
			ExpectedError: nil,
		},
		{
			Name:          "Register an existing user",
			FirstName:     "John",
			LastName:      "Doe",
			Email:         "john_doe@email.com",
			Password:      "password",
			ExpectedError: ErrUserExists,
		},
	}

	ctx := context.Background()

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			err := s.RegisterUser(ctx, tc.FirstName, tc.LastName, tc.Email, tc.Password)
			if err != tc.ExpectedError {
				t.Errorf("Expected error %v, got %v", tc.ExpectedError, err)
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	testCases := []struct {
		Name          string
		Email         string
		Password      string
		ExpectedError error
	}{
		{
			Name:          "Login a user",
			Email:         "john_doe@email.com",
			Password:      "password",
			ExpectedError: nil,
		},
		{
			Name:          "Login with wrong password",
			Email:         "john_doe@email.com",
			Password:      "wrong_password",
			ExpectedError: ErrInvalidUserCredentials,
		},
		{
			Name:          "Login a non-existing user",
			Email:         "sd@sdsd.com",
			Password:      "password",
			ExpectedError: ErrUserNotFound,
		},
	}

	ctx := context.Background()

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			_, err := s.LoginUser(ctx, tc.Email, tc.Password)
			if err != tc.ExpectedError {
				t.Errorf("Expected error %v, got %v", tc.ExpectedError, err)
			}
		})
	}
}

func TestFindUserByEmail(t *testing.T) {
	testCases := []struct {
		Name          string
		Email         string
		ExpectedError error
	}{
		{
			Name:          "Find a user",
			Email:         "john_doe@email.com",
			ExpectedError: nil,
		},
		{
			Name:          "Find a non-existing user",
			Email:         "non-exist@email.com",
			ExpectedError: ErrUserNotFound,
		},
	}

	ctx := context.Background()

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			_, err := s.FindUserByEmail(ctx, tc.Email)
			if err != tc.ExpectedError {
				t.Errorf("Expected error %v, got %v", tc.ExpectedError, err)
			}
		})
	}
}
