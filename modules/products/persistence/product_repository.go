package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/exven/pos-system/modules/products/domain"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	model := &ProductModel{}
	model.FromDomainProduct(product)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return errors.New("product with this SKU already exists")
		}
		return fmt.Errorf("failed to create product: %w", err)
	}

	product.ID = model.ID
	product.CreatedAt = model.CreatedAt
	product.UpdatedAt = model.UpdatedAt

	return nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	model := &ProductModel{}
	model.FromDomainProduct(product)

	result := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", product.ID, product.TenantID).
		Updates(model)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") || strings.Contains(result.Error.Error(), "unique constraint") {
			return errors.New("product with this SKU already exists")
		}
		return fmt.Errorf("failed to update product: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, tenantID, productID uint64) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", productID, tenantID).
		Delete(&ProductModel{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete product: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (r *productRepository) FindByID(ctx context.Context, tenantID, productID uint64) (*domain.Product, error) {
	var model ProductWithCategoryModel

	err := r.db.WithContext(ctx).
		Table("products p").
		Select("p.*, pc.name as category_name").
		Joins("LEFT JOIN product_categories pc ON p.category_id = pc.id").
		Where("p.id = ? AND p.tenant_id = ?", productID, tenantID).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	product := model.ProductModel.ToDomainProduct()

	// Add category information if exists
	if model.CategoryName != nil {
		product.Category = &domain.ProductCategory{
			ID:   *model.CategoryID,
			Name: *model.CategoryName,
		}
	}

	return product, nil
}

func (r *productRepository) FindAll(ctx context.Context, tenantID uint64, limit, offset int) ([]*domain.Product, int64, error) {
	var models []ProductWithCategoryModel
	var total int64

	// Count total records
	err := r.db.WithContext(ctx).
		Table("products").
		Where("tenant_id = ?", tenantID).
		Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Get products with category information
	err = r.db.WithContext(ctx).
		Table("products p").
		Select("p.*, pc.name as category_name").
		Joins("LEFT JOIN product_categories pc ON p.category_id = pc.id").
		Where("p.tenant_id = ?", tenantID).
		Order("p.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find products: %w", err)
	}

	products := make([]*domain.Product, len(models))
	for i, model := range models {
		products[i] = model.ProductModel.ToDomainProduct()

		// Add category information if exists
		if model.CategoryName != nil {
			products[i].Category = &domain.ProductCategory{
				ID:   *model.CategoryID,
				Name: *model.CategoryName,
			}
		}
	}

	return products, total, nil
}

func (r *productRepository) FindBySKU(ctx context.Context, tenantID uint64, sku string) (*domain.Product, error) {
	var model ProductWithCategoryModel

	err := r.db.WithContext(ctx).
		Table("products p").
		Select("p.*, pc.name as category_name").
		Joins("LEFT JOIN product_categories pc ON p.category_id = pc.id").
		Where("p.sku = ? AND p.tenant_id = ?", sku, tenantID).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to find product by SKU: %w", err)
	}

	product := model.ProductModel.ToDomainProduct()

	// Add category information if exists
	if model.CategoryName != nil {
		product.Category = &domain.ProductCategory{
			ID:   *model.CategoryID,
			Name: *model.CategoryName,
		}
	}

	return product, nil
}

func (r *productRepository) FindByBarcode(ctx context.Context, tenantID uint64, barcode string) (*domain.Product, error) {
	var model ProductWithCategoryModel

	err := r.db.WithContext(ctx).
		Table("products p").
		Select("p.*, pc.name as category_name").
		Joins("LEFT JOIN product_categories pc ON p.category_id = pc.id").
		Where("p.barcode = ? AND p.tenant_id = ?", barcode, tenantID).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, fmt.Errorf("failed to find product by barcode: %w", err)
	}

	product := model.ProductModel.ToDomainProduct()

	// Add category information if exists
	if model.CategoryName != nil {
		product.Category = &domain.ProductCategory{
			ID:   *model.CategoryID,
			Name: *model.CategoryName,
		}
	}

	return product, nil
}

func (r *productRepository) FindByCategory(ctx context.Context, tenantID, categoryID uint64, limit, offset int) ([]*domain.Product, int64, error) {
	var models []ProductWithCategoryModel
	var total int64

	// Count total records
	err := r.db.WithContext(ctx).
		Table("products").
		Where("tenant_id = ? AND category_id = ?", tenantID, categoryID).
		Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products by category: %w", err)
	}

	// Get products with category information
	err = r.db.WithContext(ctx).
		Table("products p").
		Select("p.*, pc.name as category_name").
		Joins("LEFT JOIN product_categories pc ON p.category_id = pc.id").
		Where("p.tenant_id = ? AND p.category_id = ?", tenantID, categoryID).
		Order("p.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to find products by category: %w", err)
	}

	products := make([]*domain.Product, len(models))
	for i, model := range models {
		products[i] = model.ProductModel.ToDomainProduct()

		// Add category information if exists
		if model.CategoryName != nil {
			products[i].Category = &domain.ProductCategory{
				ID:   *model.CategoryID,
				Name: *model.CategoryName,
			}
		}
	}

	return products, total, nil
}

func (r *productRepository) CheckSKUExists(ctx context.Context, tenantID uint64, sku string, excludeID *uint64) (bool, error) {
	query := r.db.WithContext(ctx).
		Model(&ProductModel{}).
		Where("tenant_id = ? AND sku = ?", tenantID, sku)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check SKU existence: %w", err)
	}

	return count > 0, nil
}

func (r *productRepository) Count(ctx context.Context, tenantID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&ProductModel{}).
		Where("tenant_id = ?", tenantID).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return count, nil
}
