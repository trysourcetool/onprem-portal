package core

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Plan struct {
	ID            uuid.UUID
	Name          string
	Price         int
	StripePriceID string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
