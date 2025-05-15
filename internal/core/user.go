package core

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

type User struct {
	ID                  uuid.UUID  `db:"id"`
	Email               string     `db:"email"`
	FirstName           string     `db:"first_name"`
	LastName            string     `db:"last_name"`
	RefreshTokenHash    string     `db:"refresh_token_hash"`
	GoogleID            string     `db:"google_id"`
	ScheduledDeletionAt *time.Time `db:"scheduled_deletion_at"`
	CreatedAt           time.Time  `db:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at"`
}

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func (u *User) SetScheduledDeletionAt() {
	const userDeletionDelay = 30 * 24 * time.Hour
	t := time.Now().Add(userDeletionDelay)
	u.ScheduledDeletionAt = &t
}
