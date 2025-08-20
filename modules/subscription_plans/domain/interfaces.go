package domain

import (
	"context"
)

type SubscriptionPlanRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]*SubscriptionPlan, int64, error)
	GetByID(ctx context.Context, id uint64) (*SubscriptionPlan, error)
}

type SubscriptionPlanService interface {
	GetAll(ctx context.Context, limit, offset int) ([]*SubscriptionPlan, int64, error)
	GetByID(ctx context.Context, id uint64) (*SubscriptionPlan, error)
}