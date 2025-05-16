package database

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/onprem-portal/internal/core"
)

type LicenseStore interface {
	GetByUserID(context.Context, uuid.UUID) (*core.License, error)
	GetByKeyHash(context.Context, string) (*core.License, error)
	Create(context.Context, *core.License) error
}
