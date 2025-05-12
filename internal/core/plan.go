package core

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Plan struct {
	ID            uuid.UUID `db:"id"`
	Name          string    `db:"name"`
	Price         int       `db:"price"`
	StripePriceID string    `db:"stripe_price_id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
