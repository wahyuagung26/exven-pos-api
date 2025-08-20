package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/exven/pos-system/modules/customers/domain"
	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) domain.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	model := &CustomerModel{}
	model.FromDomainCustomer(customer)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return errors.New("customer with this code already exists")
		}
		return fmt.Errorf("failed to create customer: %w", err)
	}

	customer.ID = model.ID
	customer.CreatedAt = model.CreatedAt
	customer.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *customerRepository) GetByID(ctx context.Context, tenantID, customerID uint64) (*domain.Customer, error) {
	var model CustomerModel

	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ? AND is_active = ?", customerID, tenantID, true).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer not found")
		}
		return nil, fmt.Errorf("failed to find customer: %w", err)
	}

	return model.ToDomainCustomer(), nil
}

func (r *customerRepository) GetByCode(ctx context.Context, tenantID uint64, code string) (*domain.Customer, error) {
	var model CustomerModel

	err := r.db.WithContext(ctx).
		Where("code = ? AND tenant_id = ? AND is_active = ?", code, tenantID, true).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer not found")
		}
		return nil, fmt.Errorf("failed to find customer by code: %w", err)
	}

	return model.ToDomainCustomer(), nil
}

func (r *customerRepository) GetByPhone(ctx context.Context, tenantID uint64, phone string) (*domain.Customer, error) {
	var model CustomerModel

	err := r.db.WithContext(ctx).
		Where("phone = ? AND tenant_id = ? AND is_active = ?", phone, tenantID, true).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer not found")
		}
		return nil, fmt.Errorf("failed to find customer by phone: %w", err)
	}

	return model.ToDomainCustomer(), nil
}

func (r *customerRepository) GetByEmail(ctx context.Context, tenantID uint64, email string) (*domain.Customer, error) {
	var model CustomerModel

	err := r.db.WithContext(ctx).
		Where("email = ? AND tenant_id = ? AND is_active = ?", email, tenantID, true).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer not found")
		}
		return nil, fmt.Errorf("failed to find customer by email: %w", err)
	}

	return model.ToDomainCustomer(), nil
}

func (r *customerRepository) GetAll(ctx context.Context, tenantID uint64, query domain.CustomerQuery) ([]*domain.Customer, int64, error) {
	var models []CustomerModel
	var total int64

	// Build base query for counting
	countQuery := r.db.WithContext(ctx).
		Model(&CustomerModel{}).
		Where("tenant_id = ?", tenantID)

	// Build base query for fetching
	fetchQuery := r.db.WithContext(ctx).
		Model(&CustomerModel{}).
		Where("tenant_id = ?", tenantID)

	// Apply filters
	if query.Name != "" {
		nameFilter := "%" + query.Name + "%"
		countQuery = countQuery.Where("name ILIKE ?", nameFilter)
		fetchQuery = fetchQuery.Where("name ILIKE ?", nameFilter)
	}

	if query.Code != "" {
		codeFilter := "%" + query.Code + "%"
		countQuery = countQuery.Where("code ILIKE ?", codeFilter)
		fetchQuery = fetchQuery.Where("code ILIKE ?", codeFilter)
	}

	if query.Email != "" {
		emailFilter := "%" + query.Email + "%"
		countQuery = countQuery.Where("email ILIKE ?", emailFilter)
		fetchQuery = fetchQuery.Where("email ILIKE ?", emailFilter)
	}

	if query.Phone != "" {
		phoneFilter := "%" + query.Phone + "%"
		countQuery = countQuery.Where("phone ILIKE ?", phoneFilter)
		fetchQuery = fetchQuery.Where("phone ILIKE ?", phoneFilter)
	}

	if query.City != "" {
		cityFilter := "%" + query.City + "%"
		countQuery = countQuery.Where("city ILIKE ?", cityFilter)
		fetchQuery = fetchQuery.Where("city ILIKE ?", cityFilter)
	}

	if query.Province != "" {
		provinceFilter := "%" + query.Province + "%"
		countQuery = countQuery.Where("province ILIKE ?", provinceFilter)
		fetchQuery = fetchQuery.Where("province ILIKE ?", provinceFilter)
	}

	if query.Gender != "" {
		countQuery = countQuery.Where("gender = ?", query.Gender)
		fetchQuery = fetchQuery.Where("gender = ?", query.Gender)
	}

	if query.IsActive != nil {
		countQuery = countQuery.Where("is_active = ?", *query.IsActive)
		fetchQuery = fetchQuery.Where("is_active = ?", *query.IsActive)
	} else {
		// Default filter for active customers only
		countQuery = countQuery.Where("is_active = ?", true)
		fetchQuery = fetchQuery.Where("is_active = ?", true)
	}

	// Count total records
	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customers: %w", err)
	}

	// Apply sorting
	orderBy := "created_at DESC"
	if query.Sort != "" {
		direction := "ASC"
		if query.Order != "" && strings.ToUpper(query.Order) == "DESC" {
			direction = "DESC"
		}

		switch query.Sort {
		case "name":
			orderBy = fmt.Sprintf("name %s", direction)
		case "code":
			orderBy = fmt.Sprintf("code %s", direction)
		case "email":
			orderBy = fmt.Sprintf("email %s", direction)
		case "phone":
			orderBy = fmt.Sprintf("phone %s", direction)
		case "city":
			orderBy = fmt.Sprintf("city %s", direction)
		case "total_spent":
			orderBy = fmt.Sprintf("total_spent %s", direction)
		case "visit_count":
			orderBy = fmt.Sprintf("visit_count %s", direction)
		case "last_visit_at":
			orderBy = fmt.Sprintf("last_visit_at %s", direction)
		case "created_at":
			orderBy = fmt.Sprintf("created_at %s", direction)
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
		return nil, 0, fmt.Errorf("failed to find customers: %w", err)
	}

	customers := make([]*domain.Customer, len(models))
	for i, model := range models {
		customers[i] = model.ToDomainCustomer()
	}

	return customers, total, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	model := &CustomerModel{}
	model.FromDomainCustomer(customer)

	result := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", customer.ID, customer.TenantID).
		Updates(model)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") || strings.Contains(result.Error.Error(), "unique constraint") {
			return errors.New("customer with this code already exists")
		}
		return fmt.Errorf("failed to update customer: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("customer not found")
	}

	return nil
}

func (r *customerRepository) Delete(ctx context.Context, tenantID, customerID uint64) error {
	result := r.db.WithContext(ctx).
		Model(&CustomerModel{}).
		Where("id = ? AND tenant_id = ?", customerID, tenantID).
		Update("is_active", false)

	if result.Error != nil {
		return fmt.Errorf("failed to delete customer: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("customer not found")
	}

	return nil
}

func (r *customerRepository) IsCodeExists(ctx context.Context, tenantID uint64, code string, excludeID *uint64) (bool, error) {
	if code == "" {
		return false, nil
	}

	query := r.db.WithContext(ctx).
		Model(&CustomerModel{}).
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

func (r *customerRepository) IsPhoneExists(ctx context.Context, tenantID uint64, phone string, excludeID *uint64) (bool, error) {
	if phone == "" {
		return false, nil
	}

	query := r.db.WithContext(ctx).
		Model(&CustomerModel{}).
		Where("tenant_id = ? AND phone = ?", tenantID, phone)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check phone existence: %w", err)
	}

	return count > 0, nil
}

func (r *customerRepository) IsEmailExists(ctx context.Context, tenantID uint64, email string, excludeID *uint64) (bool, error) {
	if email == "" {
		return false, nil
	}

	query := r.db.WithContext(ctx).
		Model(&CustomerModel{}).
		Where("tenant_id = ? AND email = ?", tenantID, email)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}

func (r *customerRepository) UpdateStats(ctx context.Context, customerID uint64, totalSpent float64, visitCount int) error {
	result := r.db.WithContext(ctx).
		Model(&CustomerModel{}).
		Where("id = ?", customerID).
		Updates(map[string]interface{}{
			"total_spent":   totalSpent,
			"visit_count":   visitCount,
			"last_visit_at": "NOW()",
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update customer stats: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("customer not found")
	}

	return nil
}