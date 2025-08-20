package persistence

import (
	"context"
	"fmt"

	"github.com/exven/pos-system/modules/auth/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	userModel := &UserModel{}
	userModel.FromDomainUser(user)

	err := r.db.WithContext(ctx).Create(userModel).Error
	if err != nil {
		return err
	}

	user.ID = userModel.ID
	user.CreatedAt = userModel.CreatedAt
	user.UpdatedAt = userModel.UpdatedAt

	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	userModel := &UserModel{}
	userModel.FromDomainUser(user)
	return r.db.WithContext(ctx).Save(userModel).Error
}

func (r *UserRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&UserModel{}, id).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*domain.User, error) {
	var userModel UserModel
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Tenant").
		First(&userModel, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return userModel.ToDomainUser(), nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, tenantID uint64, username string) (*domain.User, error) {
	var userModel UserModel
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Tenant").
		Where("tenant_id = ? AND username = ?", tenantID, username).
		First(&userModel).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return userModel.ToDomainUser(), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, tenantID uint64, email string) (*domain.User, error) {
	var userModel UserModel
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Tenant").
		Where("tenant_id = ? AND email = ?", tenantID, email).
		First(&userModel).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return userModel.ToDomainUser(), nil
}

func (r *UserRepository) FindAll(ctx context.Context, tenantID uint64, limit, offset int) ([]*domain.User, error) {
	var userModels []UserModel
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("tenant_id = ?", tenantID).
		Limit(limit).
		Offset(offset).
		Find(&userModels).Error

	if err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(userModels))
	for i, userModel := range userModels {
		users[i] = userModel.ToDomainUser()
	}

	return users, nil
}

func (r *UserRepository) FindByUsernameGlobal(ctx context.Context, username string) (*domain.User, error) {
	var userModel UserModel
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Tenant").
		Where("username = ?", username).
		First(&userModel).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return userModel.ToDomainUser(), nil
}

func (r *UserRepository) FindByEmailGlobal(ctx context.Context, email string) (*domain.User, error) {
	var userModel UserModel
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Tenant").
		Where("email = ?", email).
		First(&userModel).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return userModel.ToDomainUser(), nil
}

func (r *UserRepository) Count(ctx context.Context, tenantID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("tenant_id = ?", tenantID).
		Count(&count).Error

	return count, err
}
