package domain

import (
	"context"
)

type ProductCategoryRepository interface {
	Create(ctx context.Context, category *ProductCategory) error
	Update(ctx context.Context, category *ProductCategory) error
	Delete(ctx context.Context, tenantID, categoryID uint64) error
	FindByID(ctx context.Context, tenantID, categoryID uint64) (*ProductCategory, error)
	FindAll(ctx context.Context, tenantID uint64, query ProductCategoryQuery) ([]*ProductCategory, int64, error)
	FindByParentID(ctx context.Context, tenantID uint64, parentID *uint64) ([]*ProductCategory, error)
	CheckIfExists(ctx context.Context, tenantID uint64, name string, excludeID *uint64) (bool, error)
	HasChildren(ctx context.Context, tenantID, categoryID uint64) (bool, error)
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, tenantID, productID uint64) error
	FindByID(ctx context.Context, tenantID, productID uint64) (*Product, error)
	FindAll(ctx context.Context, tenantID uint64, limit, offset int) ([]*Product, int64, error)
	FindBySKU(ctx context.Context, tenantID uint64, sku string) (*Product, error)
	FindByBarcode(ctx context.Context, tenantID uint64, barcode string) (*Product, error)
	FindByCategory(ctx context.Context, tenantID, categoryID uint64, limit, offset int) ([]*Product, int64, error)
	CheckSKUExists(ctx context.Context, tenantID uint64, sku string, excludeID *uint64) (bool, error)
	Count(ctx context.Context, tenantID uint64) (int64, error)
}

type ProductCategoryService interface {
	Create(ctx context.Context, tenantID uint64, req CreateProductCategoryRequest) (*ProductCategory, error)
	Update(ctx context.Context, tenantID, categoryID uint64, req UpdateProductCategoryRequest) (*ProductCategory, error)
	Delete(ctx context.Context, tenantID, categoryID uint64) error
	GetByID(ctx context.Context, tenantID, categoryID uint64) (*ProductCategory, error)
	GetAll(ctx context.Context, tenantID uint64, query ProductCategoryQuery) ([]*ProductCategory, int64, error)
	GetHierarchy(ctx context.Context, tenantID uint64) ([]*ProductCategory, error)
}

type ProductService interface {
	Create(ctx context.Context, tenantID uint64, req CreateProductRequest) (*Product, error)
	Update(ctx context.Context, tenantID, productID uint64, req UpdateProductRequest) (*Product, error)
	Delete(ctx context.Context, tenantID, productID uint64) error
	GetByID(ctx context.Context, tenantID, productID uint64) (*Product, error)
	GetAll(ctx context.Context, tenantID uint64, limit, offset int) ([]*Product, int64, error)
	GetBySKU(ctx context.Context, tenantID uint64, sku string) (*Product, error)
	GetByBarcode(ctx context.Context, tenantID uint64, barcode string) (*Product, error)
	GetByCategory(ctx context.Context, tenantID, categoryID uint64, limit, offset int) ([]*Product, int64, error)
}
