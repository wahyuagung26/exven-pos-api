package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/exven/pos-system/modules/outlets/domain"
	"gorm.io/gorm"
)

type outletRepository struct {
	db *gorm.DB
}

func NewOutletRepository(db *gorm.DB) domain.OutletRepository {
	return &outletRepository{db: db}
}

func (r *outletRepository) Create(ctx context.Context, outlet *domain.Outlet) error {
	model := &OutletModel{}
	model.FromDomainOutlet(outlet)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return errors.New("outlet with this code already exists")
		}
		return fmt.Errorf("failed to create outlet: %w", err)
	}

	outlet.ID = model.ID
	outlet.CreatedAt = model.CreatedAt
	outlet.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *outletRepository) GetByID(ctx context.Context, tenantID, outletID uint64) (*domain.Outlet, error) {
	var model OutletWithManagerModel

	err := r.db.WithContext(ctx).
		Table("outlets o").
		Select("o.*, u.full_name as manager_full_name, u.email as manager_email, u.phone as manager_phone").
		Joins("LEFT JOIN users u ON o.manager_id = u.id").
		Where("o.id = ? AND o.tenant_id = ?", outletID, tenantID).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("outlet not found")
		}
		return nil, fmt.Errorf("failed to find outlet: %w", err)
	}

	return model.ToDomainOutletWithManager(), nil
}

func (r *outletRepository) GetByCode(ctx context.Context, tenantID uint64, code string) (*domain.Outlet, error) {
	var model OutletWithManagerModel

	err := r.db.WithContext(ctx).
		Table("outlets o").
		Select("o.*, u.full_name as manager_full_name, u.email as manager_email, u.phone as manager_phone").
		Joins("LEFT JOIN users u ON o.manager_id = u.id").
		Where("o.code = ? AND o.tenant_id = ?", code, tenantID).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("outlet not found")
		}
		return nil, fmt.Errorf("failed to find outlet by code: %w", err)
	}

	return model.ToDomainOutletWithManager(), nil
}

func (r *outletRepository) GetAll(ctx context.Context, tenantID uint64, query domain.OutletQuery) ([]*domain.Outlet, int64, error) {
	var models []OutletWithManagerModel
	var total int64

	// Build base query for counting
	countQuery := r.db.WithContext(ctx).
		Model(&OutletModel{}).
		Where("tenant_id = ?", tenantID)

	// Build base query for fetching
	fetchQuery := r.db.WithContext(ctx).
		Table("outlets o").
		Select("o.*, u.full_name as manager_full_name, u.email as manager_email, u.phone as manager_phone").
		Joins("LEFT JOIN users u ON o.manager_id = u.id").
		Where("o.tenant_id = ?", tenantID)

	// Apply filters
	if query.Name != "" {
		nameFilter := "%" + query.Name + "%"
		countQuery = countQuery.Where("name ILIKE ?", nameFilter)
		fetchQuery = fetchQuery.Where("o.name ILIKE ?", nameFilter)
	}

	if query.Code != "" {
		codeFilter := "%" + query.Code + "%"
		countQuery = countQuery.Where("code ILIKE ?", codeFilter)
		fetchQuery = fetchQuery.Where("o.code ILIKE ?", codeFilter)
	}

	if query.City != "" {
		cityFilter := "%" + query.City + "%"
		countQuery = countQuery.Where("city ILIKE ?", cityFilter)
		fetchQuery = fetchQuery.Where("o.city ILIKE ?", cityFilter)
	}

	if query.Province != "" {
		provinceFilter := "%" + query.Province + "%"
		countQuery = countQuery.Where("province ILIKE ?", provinceFilter)
		fetchQuery = fetchQuery.Where("o.province ILIKE ?", provinceFilter)
	}

	if query.ManagerID != nil {
		countQuery = countQuery.Where("manager_id = ?", *query.ManagerID)
		fetchQuery = fetchQuery.Where("o.manager_id = ?", *query.ManagerID)
	}

	if query.IsActive != nil {
		countQuery = countQuery.Where("is_active = ?", *query.IsActive)
		fetchQuery = fetchQuery.Where("o.is_active = ?", *query.IsActive)
	}

	// Count total records
	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count outlets: %w", err)
	}

	// Apply sorting
	orderBy := "o.created_at DESC"
	if query.Sort != "" {
		direction := "ASC"
		if query.Order != "" && strings.ToUpper(query.Order) == "DESC" {
			direction = "DESC"
		}

		switch query.Sort {
		case "name":
			orderBy = fmt.Sprintf("o.name %s", direction)
		case "code":
			orderBy = fmt.Sprintf("o.code %s", direction)
		case "city":
			orderBy = fmt.Sprintf("o.city %s", direction)
		case "created_at":
			orderBy = fmt.Sprintf("o.created_at %s", direction)
		}
	}

	// Apply pagination and fetch
	offset := (query.Page - 1) * query.Limit
	err = fetchQuery.
		Order(orderBy).
		Limit(query.Limit).
		Offset(offset).
		Find(&models).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find outlets: %w", err)
	}

	outlets := make([]*domain.Outlet, len(models))
	for i, model := range models {
		outlets[i] = model.ToDomainOutletWithManager()
	}

	return outlets, total, nil
}

func (r *outletRepository) Update(ctx context.Context, outlet *domain.Outlet) error {
	model := &OutletModel{}
	model.FromDomainOutlet(outlet)

	result := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", outlet.ID, outlet.TenantID).
		Updates(model)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") || strings.Contains(result.Error.Error(), "unique constraint") {
			return errors.New("outlet with this code already exists")
		}
		return fmt.Errorf("failed to update outlet: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("outlet not found")
	}

	return nil
}

func (r *outletRepository) Delete(ctx context.Context, tenantID, outletID uint64) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", outletID, tenantID).
		Delete(&OutletModel{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete outlet: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("outlet not found")
	}

	return nil
}

func (r *outletRepository) IsCodeExists(ctx context.Context, tenantID uint64, code string, excludeID *uint64) (bool, error) {
	query := r.db.WithContext(ctx).
		Model(&OutletModel{}).
		Where("tenant_id = ? AND code = ?", tenantID, code)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check code existence: %w", err)
	}

	return count > 0, nil
}