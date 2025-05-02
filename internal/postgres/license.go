package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/trysourcetool/onprem-portal/internal"
	"github.com/trysourcetool/onprem-portal/internal/core"
	"github.com/trysourcetool/onprem-portal/internal/database"
	"github.com/trysourcetool/onprem-portal/internal/errdefs"
)

var _ database.LicenseStore = (*licenseStore)(nil)

type licenseStore struct {
	db      internal.DB
	builder sq.StatementBuilderType
}

func newLicenseStore(db internal.DB) *licenseStore {
	return &licenseStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *licenseStore) Create(ctx context.Context, l *core.License) error {
	if _, err := s.builder.
		Insert(`"license"`).
		Columns(
			`"id"`,
			`"user_id"`,
			`"key"`,
		).
		Values(
			l.ID,
			l.UserID,
			l.Key,
		).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errdefs.ErrAlreadyExists(err)
		}
		return errdefs.ErrDatabase(err)
	}

	return nil
}
