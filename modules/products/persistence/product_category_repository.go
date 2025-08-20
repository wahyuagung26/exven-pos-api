package persistence

import (
	"context"
	"fmt"
	"strings"

	"github.com/exven/pos-system/modules/products/domain"
	"gorm.io/gorm"
)

type ProductCategoryRepository struct {
	db *gorm.DB
}

func NewProductCategoryRepository(db *gorm.DB) *ProductCategoryRepository {
	return &ProductCategoryRepository{db: db}
}

func (r *ProductCategoryRepository) Create(ctx context.Context, category *domain.ProductCategory) error {
	categoryModel := &ProductCategoryModel{}
	categoryModel.FromDomainProductCategory(category)

	err := r.db.WithContext(ctx).Create(categoryModel).Error
	if err != nil {
		return err
	}

	category.ID = categoryModel.ID
	category.CreatedAt = categoryModel.CreatedAt
	category.UpdatedAt = categoryModel.UpdatedAt

	return nil
}

func (r *ProductCategoryRepository) Update(ctx context.Context, category *domain.ProductCategory) error {
	categoryModel := &ProductCategoryModel{}
	categoryModel.FromDomainProductCategory(category)
	return r.db.WithContext(ctx).Save(categoryModel).Error
}

func (r *ProductCategoryRepository) Delete(ctx context.Context, tenantID, categoryID uint64) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, categoryID).
		Delete(&ProductCategoryModel{}).Error
}

func (r *ProductCategoryRepository) FindByID(ctx context.Context, tenantID, categoryID uint64) (*domain.ProductCategory, error) {
	var categoryModel ProductCategoryModel
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, categoryID).
		First(&categoryModel).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product category not found")
		}
		return nil, err
	}

	category := categoryModel.ToDomainProductCategory()

	// Load parent if exists
	if category.ParentID != nil {
		var parentModel ProductCategoryModel
		if err := r.db.WithContext(ctx).
			Where("tenant_id = ? AND id = ?", tenantID, *category.ParentID).
			First(&parentModel).Error; err == nil {
			category.Parent = parentModel.ToDomainProductCategory()
		}
	}

	return category, nil
}

func (r *ProductCategoryRepository) FindAll(ctx context.Context, tenantID uint64, query domain.ProductCategoryQuery) ([]*domain.ProductCategory, int64, error) {
	var categoryModels []ProductCategoryModel
	var total int64

	db := r.db.WithContext(ctx).Model(&ProductCategoryModel{}).Where("tenant_id = ?", tenantID)

	// Apply filters
	if query.ParentID != nil {
		db = db.Where("parent_id = ?", *query.ParentID)
	}

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sort := "sort_order"
	order := "ASC"

	if query.Sort != "" {
		switch query.Sort {
		case "name", "created_at", "updated_at", "sort_order":
			sort = query.Sort
		}
	}

	if query.Order != "" && strings.ToUpper(query.Order) == "DESC" {
		order = "DESC"
	}

	db = db.Order(fmt.Sprintf("%s %s", sort, order))

	// Apply pagination
	if query.Limit > 0 {
		db = db.Limit(query.Limit)
	}

	if query.Page > 0 {
		offset := (query.Page - 1) * query.Limit
		db = db.Offset(offset)
	}

	err := db.Find(&categoryModels).Error
	if err != nil {
		return nil, 0, err
	}

	categories := make([]*domain.ProductCategory, len(categoryModels))
	for i, categoryModel := range categoryModels {
		categories[i] = categoryModel.ToDomainProductCategory()
	}

	return categories, total, nil
}

func (r *ProductCategoryRepository) FindByParentID(ctx context.Context, tenantID uint64, parentID *uint64) ([]*domain.ProductCategory, error) {
	var categoryModels []ProductCategoryModel

	db := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)

	if parentID == nil {
		db = db.Where("parent_id IS NULL")
	} else {
		db = db.Where("parent_id = ?", *parentID)
	}

	err := db.Where("is_active = ?", true).
		Order("sort_order ASC, name ASC").
		Find(&categoryModels).Error

	if err != nil {
		return nil, err
	}

	categories := make([]*domain.ProductCategory, len(categoryModels))
	for i, categoryModel := range categoryModels {
		categories[i] = categoryModel.ToDomainProductCategory()
	}

	return categories, nil
}

func (r *ProductCategoryRepository) CheckIfExists(ctx context.Context, tenantID uint64, name string, excludeID *uint64) (bool, error) {
	var count int64

	db := r.db.WithContext(ctx).
		Model(&ProductCategoryModel{}).
		Where("tenant_id = ? AND name = ?", tenantID, name)

	if excludeID != nil {
		db = db.Where("id != ?", *excludeID)
	}

	err := db.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *ProductCategoryRepository) HasChildren(ctx context.Context, tenantID, categoryID uint64) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&ProductCategoryModel{}).
		Where("tenant_id = ? AND parent_id = ?", tenantID, categoryID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
