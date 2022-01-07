package models

type ChatRoom struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
