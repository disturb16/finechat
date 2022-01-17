package auth

import (
	"context"

	"github.com/disturb16/finechat/internal/auth/models"
	"github.com/jmoiron/sqlx"
)

// Repository is the auth interactor with the database.
type Repository interface {
	// SaveUser saves a user to the database.
	SaveUser(ctx context.Context, firstName, lastName, email, password string) error
	// FindUserByEmail finds a user by email.
	FindUserByEmail(ctx context.Context, email string) (*models.UserWithPassword, error)
	// SaveFriend saves a friend to the database.
	SaveFriend(ctx context.Context, email, friendEmail string) error
	// ListFriends lists all friends of a user.
	ListFriends(ctx context.Context, email string) ([]*models.User, error)
}

// Service is the auth business logic layer.
type Service interface {
	// RegisterUser creates a new user.
	RegisterUser(ctx context.Context, firstName, lastName, email, password string) error
	// LoginUser logs in a user and returns a jwt token.
	LoginUser(ctx context.Context, email, password string) (string, error)
	// FindUserByEmail finds a user by email.
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	// AddFriend adds a friend to a user.
	AddFriend(ctx context.Context, email, friendEmail string) error
	// ListFriends lists all friends of a user.
	ListFriends(ctx context.Context, email string) ([]*models.User, error)
}

// NewRepository returns a new auth repository.
func NewRepository(db *sqlx.DB) Repository {
	return &AuthRepository{db: db}
}

// NewService returns a new auth service.
func NewService(repo Repository) Service {
	return &AuthService{repo: repo}
}
