package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/disturb16/finechat/tokenparser"
	"golang.org/x/crypto/bcrypt"
)

// AuthService is the service that handles authentication.
type AuthService struct {
	repo Repository
}

// ErrUserExists is returned when the user already exists.
var ErrUserExists error = errors.New("user already exists")

// RegisterUser registers a new user.
func (s *AuthService) RegisterUser(ctx context.Context, firstName, lastName, email, password string) error {
	_, err := s.repo.FindUserByEmail(ctx, email)
	if err == nil {
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
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	return tokenparser.CreateAuthToken(u.Email, name)
}
