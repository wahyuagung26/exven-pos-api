package domain

import (
	"time"
)

type ProductCategory struct {
	ID          uint64
	TenantID    uint64
	ParentID    *uint64
	Name        string
	Description string
	ImageURL    string
	SortOrder   int
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Parent        *ProductCategory
	SubCategories []*ProductCategory
}

type Product struct {
	ID           uint64
	TenantID     uint64
	CategoryID   *uint64
	SKU          string
	Barcode      string
	Name         string
	Description  string
	Unit         string
	CostPrice    float64
	SellingPrice float64
	MinStock     int
	TrackStock   bool
	IsActive     bool
	Images       []string
	Variants     map[string]interface{}
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Category *ProductCategory
}
