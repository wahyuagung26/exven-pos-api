package persistence

import (
	"context"

	"github.com/exven/pos-system/modules/subscription_plans/domain"
	"gorm.io/gorm"
)

type subscriptionPlanRepository struct {
	db *gorm.DB
}

func NewSubscriptionPlanRepository(db *gorm.DB) domain.SubscriptionPlanRepository {
	return &subscriptionPlanRepository{
		db: db,
	}
}

func (r *subscriptionPlanRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.SubscriptionPlan, int64, error) {
	var models []SubscriptionPlanModel
	var total int64

	// Get total count for active plans only
	if err := r.db.WithContext(ctx).Model(&SubscriptionPlanModel{}).
		Where("is_active = ?", true).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated active records only
	if err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Limit(limit).
		Offset(offset).
		Order("price ASC").
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	plans := make([]*domain.SubscriptionPlan, len(models))
	for i, model := range models {
		plans[i] = r.modelToDomain(&model)
	}

	return plans, total, nil
}

func (r *subscriptionPlanRepository) GetByID(ctx context.Context, id uint64) (*domain.SubscriptionPlan, error) {
	var model SubscriptionPlanModel

	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error; err != nil {
		return nil, err
	}

	return r.modelToDomain(&model), nil
}


func (r *subscriptionPlanRepository) modelToDomain(model *SubscriptionPlanModel) *domain.SubscriptionPlan {
	return &domain.SubscriptionPlan{
		ID:                      model.ID,
		Name:                    model.Name,
		Description:             model.Description,
		Price:                   model.Price,
		MaxOutlets:              model.MaxOutlets,
		MaxUsers:                model.MaxUsers,
		MaxProducts:             model.MaxProducts,
		MaxTransactionsPerMonth: model.MaxTransactionsPerMonth,
		Features:                model.GetFeaturesSlice(),
		IsActive:                model.IsActive,
		CreatedAt:               model.CreatedAt,
		UpdatedAt:               model.UpdatedAt,
	}
}