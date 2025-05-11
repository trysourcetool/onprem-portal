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
	ID                   uuid.UUID
	UserID               uuid.UUID
	PlanID               uuid.UUID
	Status               SubscriptionStatus
	StripeCustomerID     string
	StripeSubscriptionID string
	TrialStart           time.Time
	TrialEnd             time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
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
