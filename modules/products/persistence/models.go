package persistence

import (
	"time"
	"encoding/json"
	"database/sql/driver"
	"errors"
	"github.com/exven/pos-system/modules/products/domain"
)

type ProductCategoryModel struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	TenantID    uint64    `gorm:"not null;index:idx_product_categories_tenant_parent"`
	ParentID    *uint64   `gorm:"index:idx_product_categories_tenant_parent"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	ImageURL    string    `gorm:"size:500"`
	SortOrder   int       `gorm:"default:0;index:idx_product_categories_sort_order"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (ProductCategoryModel) TableName() string {
	return "product_categories"
}

type JSONArrayModel []string

func (j JSONArrayModel) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONArrayModel) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

type JSONVariantsModel map[string]interface{}

func (j JSONVariantsModel) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONVariantsModel) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

type ProductModel struct {
	ID           uint64             `gorm:"primaryKey;autoIncrement"`
	TenantID     uint64             `gorm:"not null;uniqueIndex:idx_products_tenant_sku;index:idx_products_tenant_category"`
	CategoryID   *uint64            `gorm:"index:idx_products_tenant_category"`
	SKU          string             `gorm:"size:100;not null;uniqueIndex:idx_products_tenant_sku"`
	Barcode      string             `gorm:"size:100;index:idx_products_barcode"`
	Name         string             `gorm:"size:255;not null;index:idx_products_name"`
	Description  string             `gorm:"type:text"`
	Unit         string             `gorm:"size:50;default:'pcs'"`
	CostPrice    float64            `gorm:"type:decimal(12,2);default:0.00"`
	SellingPrice float64            `gorm:"type:decimal(12,2);not null"`
	MinStock     int                `gorm:"default:0"`
	TrackStock   bool               `gorm:"default:true"`
	IsActive     bool               `gorm:"default:true"`
	Images       JSONArrayModel     `gorm:"type:jsonb"`
	Variants     JSONVariantsModel  `gorm:"type:jsonb"`
	CreatedAt    time.Time          `gorm:"autoCreateTime"`
	UpdatedAt    time.Time          `gorm:"autoUpdateTime"`
}

func (ProductModel) TableName() string {
	return "products"
}

type ProductCategoryWithParentModel struct {
	ProductCategoryModel
	ParentName *string `gorm:"column:parent_name"`
}

type ProductWithCategoryModel struct {
	ProductModel
	CategoryName *string `gorm:"column:category_name"`
}

// Mapper functions

func (c *ProductCategoryModel) ToDomainProductCategory() *domain.ProductCategory {
	return &domain.ProductCategory{
		ID:          c.ID,
		TenantID:    c.TenantID,
		ParentID:    c.ParentID,
		Name:        c.Name,
		Description: c.Description,
		ImageURL:    c.ImageURL,
		SortOrder:   c.SortOrder,
		IsActive:    c.IsActive,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func (c *ProductCategoryModel) FromDomainProductCategory(category *domain.ProductCategory) {
	c.ID = category.ID
	c.TenantID = category.TenantID
	c.ParentID = category.ParentID
	c.Name = category.Name
	c.Description = category.Description
	c.ImageURL = category.ImageURL
	c.SortOrder = category.SortOrder
	c.IsActive = category.IsActive
	c.CreatedAt = category.CreatedAt
	c.UpdatedAt = category.UpdatedAt
}

func (p *ProductModel) ToDomainProduct() *domain.Product {
	images := make([]string, len(p.Images))
	for i, img := range p.Images {
		images[i] = img
	}
	
	variants := make(map[string]interface{})
	for k, v := range p.Variants {
		variants[k] = v
	}

	return &domain.Product{
		ID:           p.ID,
		TenantID:     p.TenantID,
		CategoryID:   p.CategoryID,
		SKU:          p.SKU,
		Barcode:      p.Barcode,
		Name:         p.Name,
		Description:  p.Description,
		Unit:         p.Unit,
		CostPrice:    p.CostPrice,
		SellingPrice: p.SellingPrice,
		MinStock:     p.MinStock,
		TrackStock:   p.TrackStock,
		IsActive:     p.IsActive,
		Images:       images,
		Variants:     variants,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func (p *ProductModel) FromDomainProduct(product *domain.Product) {
	p.ID = product.ID
	p.TenantID = product.TenantID
	p.CategoryID = product.CategoryID
	p.SKU = product.SKU
	p.Barcode = product.Barcode
	p.Name = product.Name
	p.Description = product.Description
	p.Unit = product.Unit
	p.CostPrice = product.CostPrice
	p.SellingPrice = product.SellingPrice
	p.MinStock = product.MinStock
	p.TrackStock = product.TrackStock
	p.IsActive = product.IsActive
	p.CreatedAt = product.CreatedAt
	p.UpdatedAt = product.UpdatedAt

	images := make(JSONArrayModel, len(product.Images))
	for i, img := range product.Images {
		images[i] = img
	}
	p.Images = images

	variants := make(JSONVariantsModel)
	for k, v := range product.Variants {
		variants[k] = v
	}
	p.Variants = variants
}