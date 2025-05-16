package database

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/onprem-portal/internal/core"
)

type PlanStore interface {
	GetByID(ctx context.Context, id uuid.UUID) (*core.Plan, error)
	GetByStripePriceID(ctx context.Context, stripePriceID string) (*core.Plan, error)
	List(ctx context.Context) ([]*core.Plan, error)
}
