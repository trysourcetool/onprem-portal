package database

import (
	"context"

	"github.com/trysourcetool/onprem-portal/internal/core"
)

type LicenseStore interface {
	Create(context.Context, *core.License) error
}
