package auth

import (
	"context"

	"github.com/disturb16/finechat/internal/auth/models"
	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	db *sqlx.DB
}

const (
	hasCost        int    = 20
	qrySaveUser    string = `call saveUser(?, ?, ?, ?)`
	qryFindByEmail string = `call getUserByEmail(?)`
)

func (r *AuthRepository) SaveUser(
	ctx context.Context,
	firstName, lastName, email, password string) error {

	_, err := r.db.ExecContext(
		ctx,
		qrySaveUser,
		firstName,
		lastName,
		email,
		password,
	)

	return err
}

func (r *AuthRepository) FindUserByEmail(
	ctx context.Context,
	email string,
) (*models.UserWithPassword, error) {

	u := &models.UserWithPassword{}

	err := r.db.GetContext(nil, u, qryFindByEmail, email)
	return u, err
}
