package database

import (
	"time"
)

type ArchivedTransaction struct {
	ID                    uint64                `gorm:"primaryKey"`
	TenantID              uint64                `gorm:"not null;index:idx_archived_transactions_tenant_date;index:idx_archived_trans_tenant_date"`
	OutletID              uint64                `gorm:"not null"`
	CashierID             uint64                `gorm:"not null"`
	CustomerID            *uint64
	CustomerNameSnapshot  string                `gorm:"size:255"`
	CustomerPhoneSnapshot string                `gorm:"size:20"`
	CustomerEmailSnapshot string                `gorm:"size:255"`
	CashierNameSnapshot   string                `gorm:"size:255;not null"`
	OutletNameSnapshot    string                `gorm:"size:255;not null"`
	OutletCodeSnapshot    string                `gorm:"size:50;not null"`
	TransactionNumber     string                `gorm:"size:100;not null;index:idx_archived_transactions_number"`
	TransactionDate       time.Time             `gorm:"not null;index:idx_archived_transactions_tenant_date;index:idx_archived_trans_tenant_date"`
	Subtotal              float64               `gorm:"type:decimal(15,2);not null"`
	DiscountAmount        float64               `gorm:"type:decimal(15,2);default:0.00"`
	TaxAmount             float64               `gorm:"type:decimal(15,2);default:0.00"`
	TotalAmount           float64               `gorm:"type:decimal(15,2);not null"`
	PaidAmount            float64               `gorm:"type:decimal(15,2);not null"`
	ChangeAmount          float64               `gorm:"type:decimal(15,2);default:0.00"`
	PaymentMethod         PaymentMethodType     `gorm:"not null"`
	Status                TransactionStatusType `gorm:"default:'completed'"`
	Notes                 string                `gorm:"type:text"`
	OriginalCreatedAt     time.Time             `gorm:"not null"`
	OriginalUpdatedAt     time.Time             `gorm:"not null"`
	ArchivedAt            time.Time             `gorm:"autoCreateTime;index:idx_archived_transactions_archived_date"`
	ArchivedReason        string                `gorm:"size:100;default:'data_retention_policy'"`

	ArchivedTransactionItems   []ArchivedTransactionItem   `gorm:"foreignKey:TransactionID"`
	ArchivedTransactionPayments []ArchivedTransactionPayment `gorm:"foreignKey:TransactionID"`
}

func (ArchivedTransaction) TableName() string {
	return "archived_transactions"
}

type ArchivedTransactionItem struct {
	ID                      uint64    `gorm:"primaryKey"`
	TransactionID           uint64    `gorm:"not null;index:idx_archived_transaction_items_transaction;index:idx_archived_items_trans_product"`
	ProductID               uint64    `gorm:"not null;index:idx_archived_items_trans_product"`
	ProductNameSnapshot     string    `gorm:"size:255;not null"`
	ProductSKUSnapshot      string    `gorm:"size:100;not null;index:idx_archived_transaction_items_product_sku_snapshot"`
	ProductCategorySnapshot string    `gorm:"size:255"`
	ProductUnitSnapshot     string    `gorm:"size:50;default:'pcs'"`
	Quantity                int       `gorm:"not null"`
	UnitPrice               float64   `gorm:"type:decimal(12,2);not null"`
	CostPriceSnapshot       float64   `gorm:"type:decimal(12,2);default:0.00"`
	DiscountAmount          float64   `gorm:"type:decimal(12,2);default:0.00"`
	TotalPrice              float64   `gorm:"type:decimal(15,2);not null"`
	Notes                   string    `gorm:"type:text"`
	ArchivedAt              time.Time `gorm:"autoCreateTime;index:idx_archived_transaction_items_archived_date"`

	ArchivedTransaction ArchivedTransaction `gorm:"foreignKey:TransactionID"`
}

func (ArchivedTransactionItem) TableName() string {
	return "archived_transaction_items"
}

type ArchivedTransactionPayment struct {
	ID                uint64              `gorm:"primaryKey"`
	TransactionID     uint64              `gorm:"not null;index:idx_archived_transaction_payments_transaction"`
	PaymentMethod     PaymentMethodSingle `gorm:"not null"`
	Amount            float64             `gorm:"type:decimal(15,2);not null"`
	ReferenceNumber   string              `gorm:"size:100"`
	Notes             string              `gorm:"type:text"`
	OriginalCreatedAt time.Time           `gorm:"not null"`
	ArchivedAt        time.Time           `gorm:"autoCreateTime;index:idx_archived_transaction_payments_archived_date"`

	ArchivedTransaction ArchivedTransaction `gorm:"foreignKey:TransactionID"`
}

func (ArchivedTransactionPayment) TableName() string {
	return "archived_transaction_payments"
}