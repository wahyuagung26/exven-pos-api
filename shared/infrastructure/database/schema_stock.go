package database

import (
	"time"
)

type MovementType string
type ReferenceType string

const (
	MovementTypeIn         MovementType = "in"
	MovementTypeOut        MovementType = "out"
	MovementTypeAdjustment MovementType = "adjustment"
	MovementTypeTransfer   MovementType = "transfer"

	ReferenceTypeSale       ReferenceType = "sale"
	ReferenceTypePurchase   ReferenceType = "purchase"
	ReferenceTypeAdjustment ReferenceType = "adjustment"
	ReferenceTypeTransfer   ReferenceType = "transfer"
	ReferenceTypeInitial    ReferenceType = "initial"
)

type StockMovement struct {
	ID            uint64        `gorm:"primaryKey;autoIncrement"`
	ProductID     uint64        `gorm:"not null;index:idx_stock_movements_product_outlet"`
	OutletID      uint64        `gorm:"not null;index:idx_stock_movements_product_outlet;index:idx_stock_movements_outlet_date"`
	MovementType  MovementType  `gorm:"not null"`
	Quantity      int           `gorm:"not null"`
	ReferenceType ReferenceType `gorm:"not null;index:idx_stock_movements_reference"`
	ReferenceID   *uint64       `gorm:"index:idx_stock_movements_reference"`
	Notes         string        `gorm:"type:text"`
	CreatedBy     uint64        `gorm:"not null"`
	CreatedAt     time.Time     `gorm:"autoCreateTime;index:idx_stock_movements_outlet_date"`

	Product   Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Outlet    Outlet  `gorm:"foreignKey:OutletID;constraint:OnDelete:CASCADE"`
	CreatedByUser User `gorm:"foreignKey:CreatedBy"`
}