package auth

import (
	"context"
	"errors"

	"github.com/disturb16/finechat/internal/auth/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *AuthRepository
}

var ErrUserExists error = errors.New("user already exists")

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

	passToSave := string(hashedpass)

	return s.repo.SaveUser(ctx, firstName, lastName, email, passToSave)
}

// LoginUser authenticates an user.
func (s *AuthService) LoginUser(ctx context.Context, email, password string) (*models.User, error) {

	u, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return u.ToUser(), nil
}
