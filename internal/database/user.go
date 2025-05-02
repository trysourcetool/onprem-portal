package database

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/onprem-portal/internal/core"
)

type UserStore interface {
	GetByID(context.Context, uuid.UUID) (*core.User, error)
	GetByRefreshTokenHash(context.Context, string) (*core.User, error)
	GetByEmail(context.Context, string) (*core.User, error)
	GetByGoogleID(context.Context, string) (*core.User, error)
	Create(context.Context, *core.User) error
	Update(context.Context, *core.User) error
	IsEmailExists(context.Context, string) (bool, error)
}
