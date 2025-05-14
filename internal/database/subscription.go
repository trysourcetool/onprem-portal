package database

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/onprem-portal/internal/core"
)

type SubscriptionStore interface {
	GetByID(ctx context.Context, id uuid.UUID) (*core.Subscription, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*core.Subscription, error)
	GetByStripeSubscriptionID(ctx context.Context, stripeSubID string) (*core.Subscription, error)
	Create(ctx context.Context, s *core.Subscription) error
	Update(ctx context.Context, s *core.Subscription) error
	UpdateSeatCount(ctx context.Context, subID uuid.UUID, addSeatCount int64) (int64, error)
	DeleteByID(ctx context.Context, id uuid.UUID) error
	ListExpiredTrial(ctx context.Context, before time.Time) ([]*core.Subscription, error)
}
