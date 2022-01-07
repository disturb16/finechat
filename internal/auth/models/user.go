package models

import "time"

type User struct {
	ID        int64  `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

type UserWithPassword struct {
	ID          int64     `db:"id"`
	FirstName   string    `db:"first_name"`
	LastName    string    `db:"last_name"`
	Email       string    `db:"email"`
	Password    string    `db:"password"`
	CreatedDate time.Time `db:"created_date"`
}

func (u *UserWithPassword) ToUser() *User {
	return &User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
}
