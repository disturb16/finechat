package auth

import (
	"context"

	"github.com/disturb16/finechat/internal/auth/models"
)

type MockRepository struct{}

func (mr *MockRepository) SaveUser(ctx context.Context, firstName, lastName, email, password string) error {
	return nil
}

func (mr *MockRepository) FindUserByEmail(ctx context.Context, email string) (*models.UserWithPassword, error) {
	if email == "john_doe@email.com" {
		return &models.UserWithPassword{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
			Email:     email,
			Password:  "$2a$12$7oK5/y.FZHmj/mpPFe7vp.H0nqlDcfjP3wCJuPJEBMe.3peyYUeee",
		}, nil
	}

	return nil, nil
}

func (mr *MockRepository) SaveFriend(ctx context.Context, email, friendEmail string) error {
	return nil
}

func (mr *MockRepository) ListFriends(ctx context.Context, email string) ([]*models.User, error) {
	return nil, nil
}
