package services

import (
	"context"

	"github.com/exven/pos-system/modules/subscription_plans/domain"
)

type subscriptionPlanService struct {
	repo domain.SubscriptionPlanRepository
}

func NewSubscriptionPlanService(repo domain.SubscriptionPlanRepository) domain.SubscriptionPlanService {
	return &subscriptionPlanService{
		repo: repo,
	}
}

func (s *subscriptionPlanService) GetAll(ctx context.Context, limit, offset int) ([]*domain.SubscriptionPlan, int64, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *subscriptionPlanService) GetByID(ctx context.Context, id uint64) (*domain.SubscriptionPlan, error) {
	return s.repo.GetByID(ctx, id)
}

