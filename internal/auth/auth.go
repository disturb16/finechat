package auth

import (
	"context"

	"github.com/disturb16/finechat/internal/auth/models"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	SaveUser(ctx context.Context, firstName, lastName, email, password string) error
	FindUserByEmail(ctx context.Context, email string) (*models.UserWithPassword, error)
	SaveFriend(ctx context.Context, email, friendEmail string) error
	ListFriends(ctx context.Context, email string) ([]*models.User, error)
}

type Service interface {
	RegisterUser(ctx context.Context, firstName, lastName, email, password string) error
	LoginUser(ctx context.Context, email, password string) (string, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	AddFriend(ctx context.Context, email, friendEmail string) error
	ListFriends(ctx context.Context, email string) ([]*models.User, error)
}

func NewRepository(db *sqlx.DB) Repository {
	return &AuthRepository{db: db}
}

func NewService(repo Repository) Service {
	return &AuthService{repo: repo}
}
