package database

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/onprem-portal/internal/core"
)

type UserStore interface {
	GetByID(ctx context.Context, id uuid.UUID) (*core.User, error)
	GetByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*core.User, error)
	GetByEmail(ctx context.Context, email string) (*core.User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*core.User, error)
	Create(ctx context.Context, user *core.User) error
	Update(ctx context.Context, user *core.User) error
	IsEmailExists(ctx context.Context, email string) (bool, error)
}
