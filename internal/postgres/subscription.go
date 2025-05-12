package postgres

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/onprem-portal/internal"
	"github.com/trysourcetool/onprem-portal/internal/core"
	"github.com/trysourcetool/onprem-portal/internal/database"
)

type subscriptionStore struct {
	db      internal.DB
	builder sq.StatementBuilderType
}

func newSubscriptionStore(db internal.DB) database.SubscriptionStore {
	return &subscriptionStore{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *subscriptionStore) GetByID(ctx context.Context, id uuid.UUID) (*core.Subscription, error) {
	query, args, err := s.builder.
		Select(`s."id"`, `s."user_id"`, `s."plan_id"`, `s."status"`, `s."stripe_customer_id"`, `s."stripe_subscription_id"`, `s."trial_start"`, `s."trial_end"`, `s."created_at"`, `s."updated_at"`).
		From(`"subscription" s`).
		Where(sq.Eq{`s."id"`: id}).
		ToSql()
	if err != nil {
		return nil, err
	}
	var sub core.Subscription
	if err := s.db.GetContext(ctx, &sub, query, args...); err != nil {
		return nil, err
	}
	return &sub, nil
}

func (s *subscriptionStore) GetByUserID(ctx context.Context, userID uuid.UUID) (*core.Subscription, error) {
	query, args, err := s.builder.
		Select(`s."id"`, `s."user_id"`, `s."plan_id"`, `s."status"`, `s."stripe_customer_id"`, `s."stripe_subscription_id"`, `s."trial_start"`, `s."trial_end"`, `s."created_at"`, `s."updated_at"`).
		From(`"subscription" s`).
		Where(sq.Eq{`s."user_id"`: userID}).
		ToSql()
	if err != nil {
		return nil, err
	}
	var sub core.Subscription
	if err := s.db.GetContext(ctx, &sub, query, args...); err != nil {
		return nil, err
	}
	return &sub, nil
}

func (s *subscriptionStore) Create(ctx context.Context, sub *core.Subscription) error {
	_, err := s.builder.
		Insert(`"subscription"`).
		Columns(`"id"`, `"user_id"`, `"plan_id"`, `"status"`, `"stripe_customer_id"`, `"stripe_subscription_id"`, `"trial_start"`, `"trial_end"`).
		Values(sub.ID, sub.UserID, sub.PlanID, sub.Status, sub.StripeCustomerID, sub.StripeSubscriptionID, sub.TrialStart, sub.TrialEnd).
		RunWith(s.db).
		ExecContext(ctx)
	return err
}

func (s *subscriptionStore) Update(ctx context.Context, sub *core.Subscription) error {
	_, err := s.builder.
		Update(`"subscription"`).
		Set(`"plan_id"`, sub.PlanID).
		Set(`"status"`, sub.Status).
		Set(`"stripe_customer_id"`, sub.StripeCustomerID).
		Set(`"stripe_subscription_id"`, sub.StripeSubscriptionID).
		Set(`"trial_start"`, sub.TrialStart).
		Set(`"trial_end"`, sub.TrialEnd).
		Where(sq.Eq{`"id"`: sub.ID}).
		RunWith(s.db).
		ExecContext(ctx)
	return err
}

func (s *subscriptionStore) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.builder.
		Delete(`"subscription"`).
		Where(sq.Eq{`"id"`: id}).
		RunWith(s.db).
		ExecContext(ctx)
	return err
}

func (s *subscriptionStore) ListExpiredTrial(ctx context.Context, before time.Time) ([]*core.Subscription, error) {
	query, args, err := s.builder.
		Select(`s."id"`, `s."user_id"`, `s."plan_id"`, `s."status"`, `s."stripe_customer_id"`, `s."stripe_subscription_id"`, `s."trial_start"`, `s."trial_end"`, `s."created_at"`, `s."updated_at"`).
		From(`"subscription" s`).
		Where(sq.And{
			sq.Eq{`s."status"`: int(core.SubscriptionStatusTrial)},
			sq.Lt{`s."trial_end"`: before},
		}).
		ToSql()
	if err != nil {
		return nil, err
	}
	var subs []*core.Subscription
	if err := s.db.SelectContext(ctx, &subs, query, args...); err != nil {
		return nil, err
	}
	return subs, nil
}

func (s *subscriptionStore) GetByStripeSubscriptionID(ctx context.Context, stripeSubID string) (*core.Subscription, error) {
	query, args, err := s.builder.
		Select(`s."id"`, `s."user_id"`, `s."plan_id"`, `s."status"`, `s."stripe_customer_id"`, `s."stripe_subscription_id"`, `s."trial_start"`, `s."trial_end"`, `s."created_at"`, `s."updated_at"`).
		From(`"subscription" s`).
		Where(sq.Eq{`s."stripe_subscription_id"`: stripeSubID}).
		ToSql()
	if err != nil {
		return nil, err
	}
	var sub core.Subscription
	if err := s.db.GetContext(ctx, &sub, query, args...); err != nil {
		return nil, err
	}
	return &sub, nil
}
