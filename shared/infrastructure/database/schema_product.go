package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type ProductCategory struct {
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

	Tenant        Tenant            `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
	Parent        *ProductCategory  `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL"`
	SubCategories []ProductCategory `gorm:"foreignKey:ParentID"`
	Products      []Product         `gorm:"foreignKey:CategoryID"`
}

type JSONArray []string

func (j JSONArray) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONArray) Scan(value interface{}) error {
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

type JSONVariants map[string]interface{}

func (j JSONVariants) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONVariants) Scan(value interface{}) error {
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

type Product struct {
	ID           uint64       `gorm:"primaryKey;autoIncrement"`
	TenantID     uint64       `gorm:"not null;uniqueIndex:idx_products_tenant_sku;index:idx_products_tenant_category"`
	CategoryID   *uint64      `gorm:"index:idx_products_tenant_category"`
	SKU          string       `gorm:"size:100;not null;uniqueIndex:idx_products_tenant_sku"`
	Barcode      string       `gorm:"size:100;index:idx_products_barcode"`
	Name         string       `gorm:"size:255;not null;index:idx_products_name"`
	Description  string       `gorm:"type:text"`
	Unit         string       `gorm:"size:50;default:'pcs'"`
	CostPrice    float64      `gorm:"type:decimal(12,2);default:0.00"`
	SellingPrice float64      `gorm:"type:decimal(12,2);not null"`
	MinStock     int          `gorm:"default:0"`
	TrackStock   bool         `gorm:"default:true"`
	IsActive     bool         `gorm:"default:true"`
	Images       JSONArray    `gorm:"type:jsonb"`
	Variants     JSONVariants `gorm:"type:jsonb"`
	CreatedAt    time.Time    `gorm:"autoCreateTime"`
	UpdatedAt    time.Time    `gorm:"autoUpdateTime"`

	Tenant           Tenant            `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
	Category         *ProductCategory  `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL"`
	ProductStocks    []ProductStock    `gorm:"foreignKey:ProductID"`
	TransactionItems []TransactionItem `gorm:"foreignKey:ProductID"`
	StockMovements   []StockMovement   `gorm:"foreignKey:ProductID"`
}

type ProductStock struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement"`
	ProductID        uint64    `gorm:"not null;uniqueIndex:idx_product_stocks_product_outlet"`
	OutletID         uint64    `gorm:"not null;uniqueIndex:idx_product_stocks_product_outlet;index:idx_product_stocks_outlet_quantity"`
	Quantity         int       `gorm:"not null;default:0;index:idx_product_stocks_outlet_quantity"`
	ReservedQuantity int       `gorm:"default:0"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`

	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Outlet  Outlet  `gorm:"foreignKey:OutletID;constraint:OnDelete:CASCADE"`
}
