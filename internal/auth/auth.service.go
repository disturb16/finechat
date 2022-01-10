package auth

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/disturb16/finechat/internal/auth/models"
	"github.com/disturb16/finechat/tokenparser"
	"golang.org/x/crypto/bcrypt"
)

// AuthService is the service that handles authentication.
type AuthService struct {
	repo Repository
}

var (
	// ErrUserExists is returned when the user already exists.
	ErrUserExists error = errors.New("user already exists")
	// ErrUserNotFound is returned when the user is not found.
	ErrUserNotFound           error = errors.New("user not found")
	ErrInvalidUserCredentials error = errors.New("invalid user credentials")
)

// RegisterUser registers a new user.
func (s *AuthService) RegisterUser(ctx context.Context, firstName, lastName, email, password string) error {
	u, _ := s.repo.FindUserByEmail(ctx, email)
	if u != nil {
		return ErrUserExists
	}

	hashedpass, err := bcrypt.GenerateFromPassword([]byte(password), hasCost)
	if err != nil {
		return err
	}

	return s.repo.SaveUser(ctx, firstName, lastName, email, string(hashedpass))
}

// LoginUser authenticates an user and returns the jwt.
func (s *AuthService) LoginUser(ctx context.Context, email, password string) (string, error) {
	u, err := s.repo.FindUserByEmail(ctx, email)
	if u == nil || err != nil {
		log.Println(err)
		return "", ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		log.Println(err)
		return "", ErrInvalidUserCredentials
	}

	name := fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	return tokenparser.CreateAuthToken(u.Email, name)
}

func (s *AuthService) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	u, err := s.repo.FindUserByEmail(ctx, email)
	if u == nil || err != nil {
		log.Println(err)
		return nil, ErrUserNotFound
	}

	return u.ToUser(), nil
}

func (s *AuthService) AddFriend(ctx context.Context, email, friendEmail string) error {
	return s.repo.SaveFriend(ctx, email, friendEmail)
}

func (s *AuthService) ListFriends(ctx context.Context, email string) ([]*models.User, error) {
	return s.repo.ListFriends(ctx, email)
}
