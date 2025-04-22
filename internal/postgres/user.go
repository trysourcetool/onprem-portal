package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/trysourcetool/onprem-portal/internal"
	"github.com/trysourcetool/onprem-portal/internal/core"
	"github.com/trysourcetool/onprem-portal/internal/database"
)

var _ database.UserStore = (*userStore)(nil)

type userStore struct {
	db      internal.DB
	builder sq.StatementBuilderType
}

func newUserStore(db internal.DB) *userStore {
	return &userStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *userStore) GetByEmail(ctx context.Context, email string) (*core.User, error) {
	return nil, nil
}

func (s *userStore) GetByGoogleID(ctx context.Context, googleID string) (*core.User, error) {
	return nil, nil
}

func (s *userStore) Create(ctx context.Context, user *core.User) error {
	return nil
}

func (s *userStore) Update(ctx context.Context, user *core.User) error {
	return nil
}
