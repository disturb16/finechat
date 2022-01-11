package auth

import (
	"context"

	"github.com/disturb16/finechat/internal/auth/models"
)

type MockAuthService struct{}

func (mas *MockAuthService) RegisterUser(ctx context.Context, firstName, lastName, email, password string) error {
	return nil
}

func (mas *MockAuthService) LoginUser(ctx context.Context, email, password string) (string, error) {
	if email != "john@doe.com" || password != "password" {
		return "", ErrInvalidUserCredentials
	}

	return "token", nil
}

func (mas *MockAuthService) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email != "john@doe.com" {
		return nil, ErrInvalidUserCredentials
	}

	return &models.User{ID: 1, Email: email}, nil
}

func (mas *MockAuthService) AddFriend(ctx context.Context, email, friendEmail string) error {
	return nil

}
func (mas *MockAuthService) ListFriends(ctx context.Context, email string) ([]*models.User, error) {
	return nil, nil
}
