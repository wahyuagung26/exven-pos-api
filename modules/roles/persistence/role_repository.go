package persistence

import (
	"context"

	"github.com/exven/pos-system/modules/roles/domain"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) domain.RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Role, int64, error) {
	var models []RoleModel
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&RoleModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	roles := make([]*domain.Role, len(models))
	for i, model := range models {
		roles[i] = r.modelToDomain(&model)
	}

	return roles, total, nil
}

func (r *roleRepository) GetByID(ctx context.Context, id uint64) (*domain.Role, error) {
	var model RoleModel

	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error; err != nil {
		return nil, err
	}

	return r.modelToDomain(&model), nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	var model RoleModel

	if err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&model).Error; err != nil {
		return nil, err
	}

	return r.modelToDomain(&model), nil
}

func (r *roleRepository) GetSystemRoles(ctx context.Context) ([]*domain.Role, error) {
	var models []RoleModel

	if err := r.db.WithContext(ctx).
		Where("is_system = ?", true).
		Order("name ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	roles := make([]*domain.Role, len(models))
	for i, model := range models {
		roles[i] = r.modelToDomain(&model)
	}

	return roles, nil
}

func (r *roleRepository) modelToDomain(model *RoleModel) *domain.Role {
	return &domain.Role{
		ID:          model.ID,
		Name:        model.Name,
		DisplayName: model.DisplayName,
		Description: model.Description,
		Permissions: model.GetPermissionsSlice(),
		IsSystem:    model.IsSystem,
		CreatedAt:   model.CreatedAt,
	}
}