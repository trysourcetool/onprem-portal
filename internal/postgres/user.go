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

func (s *userStore) GetByID(ctx context.Context, id uuid.UUID) (*core.User, error) {
	query, args, err := s.builder.
		Select(s.columns()...).
		From(`"user" u`).
		Where(sq.Eq{`u."id"`: id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var u core.User
	if err := s.db.GetContext(ctx, &u, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserNotFound(err)
		}
		return nil, err
	}

	return &u, nil
}

func (s *userStore) GetByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*core.User, error) {
	query, args, err := s.builder.
		Select(s.columns()...).
		From(`"user" u`).
		Where(sq.Eq{`u."refresh_token_hash"`: refreshTokenHash}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var u core.User
	if err := s.db.GetContext(ctx, &u, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserNotFound(err)
		}
		return nil, err
	}

	return &u, nil
}

func (s *userStore) GetByEmail(ctx context.Context, email string) (*core.User, error) {
	query, args, err := s.builder.
		Select(s.columns()...).
		From(`"user" u`).
		Where(sq.Eq{`u."email"`: email}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var u core.User
	if err := s.db.GetContext(ctx, &u, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserNotFound(err)
		}
		return nil, err
	}

	return &u, nil
}

func (s *userStore) GetByGoogleID(ctx context.Context, googleID string) (*core.User, error) {
	query, args, err := s.builder.
		Select(s.columns()...).
		From(`"user" u`).
		Where(sq.Eq{`u."google_id"`: googleID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var u core.User
	if err := s.db.GetContext(ctx, &u, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.ErrUserNotFound(err)
		}
		return nil, err
	}

	return &u, nil
}

func (s *userStore) Create(ctx context.Context, u *core.User) error {
	if _, err := s.builder.
		Insert(`"user"`).
		Columns(
			`"id"`,
			`"email"`,
			`"first_name"`,
			`"last_name"`,
			`"refresh_token_hash"`,
			`"google_id"`,
		).
		Values(
			u.ID,
			u.Email,
			u.FirstName,
			u.LastName,
			u.RefreshTokenHash,
			u.GoogleID,
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

func (s *userStore) Update(ctx context.Context, u *core.User) error {
	if _, err := s.builder.
		Update(`"user"`).
		Set(`"email"`, u.Email).
		Set(`"first_name"`, u.FirstName).
		Set(`"last_name"`, u.LastName).
		Set(`"refresh_token_hash"`, u.RefreshTokenHash).
		Set(`"google_id"`, u.GoogleID).
		Where(sq.Eq{`"id"`: u.ID}).
		RunWith(s.db).
		ExecContext(ctx); err != nil {
		return errdefs.ErrDatabase(err)
	}

	return nil
}

func (s *userStore) IsEmailExists(ctx context.Context, email string) (bool, error) {
	if _, err := s.GetByEmail(ctx, email); err != nil {
		if errdefs.IsUserNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *userStore) columns() []string {
	return []string{
		`u."id"`,
		`u."email"`,
		`u."first_name"`,
		`u."last_name"`,
		`u."google_id"`,
		`u."refresh_token_hash"`,
		`u."created_at"`,
		`u."updated_at"`,
	}
}
