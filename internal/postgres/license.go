package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
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

func (s *licenseStore) GetByUserID(ctx context.Context, userID uuid.UUID) (*core.License, error) {
	query, args, err := s.builder.
		Select(s.columns()...).
		From(`"license" l`).
		Where(sq.Eq{`l."user_id"`: userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var l core.License
	if err := s.db.GetContext(ctx, &l, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrLicenseNotFound(err)
		}
		return nil, err
	}

	return &l, nil
}

func (s *licenseStore) Create(ctx context.Context, l *core.License) error {
	if _, err := s.builder.
		Insert(`"license"`).
		Columns(
			`"id"`,
			`"user_id"`,
			`"key_hash"`,
			`"key_ciphertext"`,
			`"key_nonce"`,
		).
		Values(
			l.ID,
			l.UserID,
			l.KeyHash,
			l.KeyCiphertext,
			l.KeyNonce,
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

func (s *licenseStore) columns() []string {
	return []string{
		`l."id"`,
		`l."user_id"`,
		`l."key_hash"`,
		`l."key_ciphertext"`,
		`l."key_nonce"`,
	}
}
