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
	hasCost           int    = 12
	qrySaveUser       string = `call saveUser(?, ?, ?, ?)`
	qryFindByEmail    string = `call getUserByEmail(?)`
	qrySaveFriend     string = `call saveFriend(?, ?)`
	qryGetUserFriends string = `call getUserFriends(?)`
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

	err := r.db.GetContext(ctx, u, qryFindByEmail, email)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *AuthRepository) SaveFriend(
	ctx context.Context,
	email, friendEmail string,
) error {

	_, err := r.db.ExecContext(
		ctx,
		qrySaveFriend,
		email,
		friendEmail,
	)

	return err
}

func (r *AuthRepository) ListFriends(
	ctx context.Context,
	email string,
) ([]*models.User, error) {

	friends := []*models.User{}

	err := r.db.SelectContext(
		ctx,
		&friends,
		qryGetUserFriends,
		email,
	)

	return friends, err
}
