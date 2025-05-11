package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/onprem-portal/internal"
	"github.com/trysourcetool/onprem-portal/internal/core"
	"github.com/trysourcetool/onprem-portal/internal/database"
)

type planStore struct {
	db      internal.DB
	builder sq.StatementBuilderType
}

func newPlanStore(db internal.DB) database.PlanStore {
	return &planStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *planStore) GetByID(ctx context.Context, id uuid.UUID) (*core.Plan, error) {
	query, args, err := s.builder.
		Select("id", "name", "price", "stripe_price_id", "created_at", "updated_at").
		From(`plan`).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}
	var p core.Plan
	if err := s.db.GetContext(ctx, &p, query, args...); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *planStore) List(ctx context.Context) ([]*core.Plan, error) {
	query, args, err := s.builder.
		Select("id", "name", "price", "stripe_price_id", "created_at", "updated_at").
		From(`plan`).
		ToSql()
	if err != nil {
		return nil, err
	}
	var plans []*core.Plan
	if err := s.db.SelectContext(ctx, &plans, query, args...); err != nil {
		return nil, err
	}
	return plans, nil
}
