package core

import "time"

type License struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Key       string    `db:"key"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
