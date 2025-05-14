package core

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type SubscriptionStatus int

const (
	SubscriptionStatusUnknown SubscriptionStatus = iota
	SubscriptionStatusTrial
	SubscriptionStatusActive
	SubscriptionStatusCanceled
	SubscriptionStatusPastDue

	subscriptionStatusUnknown  = "unknown"
	subscriptionStatusTrial    = "trial"
	subscriptionStatusActive   = "active"
	subscriptionStatusCanceled = "canceled"
	subscriptionStatusPastDue  = "past_due"
)

func (s SubscriptionStatus) String() string {
	subscriptionStatuses := []string{
		subscriptionStatusUnknown,
		subscriptionStatusTrial,
		subscriptionStatusActive,
		subscriptionStatusCanceled,
		subscriptionStatusPastDue,
	}

	if int(s) < 0 || int(s) >= len(subscriptionStatuses) {
		return subscriptionStatusUnknown
	}

	return subscriptionStatuses[s]
}

func SubscriptionStatusFromString(s string) SubscriptionStatus {
	statusMap := map[string]SubscriptionStatus{
		subscriptionStatusTrial:    SubscriptionStatusTrial,
		subscriptionStatusActive:   SubscriptionStatusActive,
		subscriptionStatusCanceled: SubscriptionStatusCanceled,
		subscriptionStatusPastDue:  SubscriptionStatusPastDue,
	}
	if status, ok := statusMap[s]; ok {
		return status
	}
	return SubscriptionStatusUnknown
}

type Subscription struct {
	ID                   uuid.UUID          `db:"id"`
	UserID               uuid.UUID          `db:"user_id"`
	PlanID               *uuid.UUID         `db:"plan_id"`
	Status               SubscriptionStatus `db:"status"`
	StripeCustomerID     string             `db:"stripe_customer_id"`
	StripeSubscriptionID string             `db:"stripe_subscription_id"`
	TrialStart           time.Time          `db:"trial_start"`
	TrialEnd             time.Time          `db:"trial_end"`
	SeatCount            int64              `db:"seat_count"`
	CreatedAt            time.Time          `db:"created_at"`
	UpdatedAt            time.Time          `db:"updated_at"`
}

func (s *Subscription) IsTrial() bool {
	return s.Status == SubscriptionStatusTrial
}

func (s *Subscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive
}

func (s *Subscription) IsCanceled() bool {
	return s.Status == SubscriptionStatusCanceled
}

func (s *Subscription) IsPastDue() bool {
	return s.Status == SubscriptionStatusPastDue
}
