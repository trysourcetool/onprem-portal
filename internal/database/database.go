package database

import (
	"context"
)

type Stores interface {
	License() LicenseStore
	User() UserStore
	Plan() PlanStore
	Subscription() SubscriptionStore
}

type DB interface {
	Stores
	WithTx(ctx context.Context, fn func(tx Tx) error) error
}

type Tx interface {
	Stores
}
