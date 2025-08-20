package database

import (
	"time"
)

type PaymentMethodType string
type TransactionStatusType string
type PaymentMethodSingle string

const (
	PaymentMethodCash     PaymentMethodType = "cash"
	PaymentMethodCard     PaymentMethodType = "card"
	PaymentMethodTransfer PaymentMethodType = "transfer"
	PaymentMethodEwallet  PaymentMethodType = "ewallet"
	PaymentMethodMultiple PaymentMethodType = "multiple"

	TransactionStatusCompleted TransactionStatusType = "completed"
	TransactionStatusCancelled TransactionStatusType = "cancelled"
	TransactionStatusRefunded  TransactionStatusType = "refunded"

	PaymentSingleCash     PaymentMethodSingle = "cash"
	PaymentSingleCard     PaymentMethodSingle = "card"
	PaymentSingleTransfer PaymentMethodSingle = "transfer"
	PaymentSingleEwallet  PaymentMethodSingle = "ewallet"
)

type SalesTransaction struct {
	ID                    uint64                `gorm:"primaryKey;autoIncrement"`
	TenantID              uint64                `gorm:"not null;index:idx_transactions_tenant_outlet_date;index:idx_sales_tenant_date_status"`
	OutletID              uint64                `gorm:"not null;index:idx_transactions_tenant_outlet_date;index:idx_sales_outlet_date_total"`
	CashierID             uint64                `gorm:"not null;index:idx_transactions_cashier_date"`
	CustomerID            *uint64               `gorm:"constraint:OnDelete:SET NULL"`
	CustomerNameSnapshot  string                `gorm:"size:255;index:idx_transactions_customer_snapshot"`
	CustomerPhoneSnapshot string                `gorm:"size:20;index:idx_transactions_customer_snapshot"`
	CustomerEmailSnapshot string                `gorm:"size:255"`
	CashierNameSnapshot   string                `gorm:"size:255;not null;index:idx_transactions_cashier_snapshot"`
	OutletNameSnapshot    string                `gorm:"size:255;not null"`
	OutletCodeSnapshot    string                `gorm:"size:50;not null;index:idx_transactions_outlet_snapshot"`
	TransactionNumber     string                `gorm:"size:100;not null;uniqueIndex"`
	TransactionDate       time.Time             `gorm:"default:CURRENT_TIMESTAMP;index:idx_transactions_tenant_outlet_date;index:idx_transactions_cashier_date;index:idx_sales_tenant_date_status;index:idx_sales_outlet_date_total"`
	Subtotal              float64               `gorm:"type:decimal(15,2);not null"`
	DiscountAmount        float64               `gorm:"type:decimal(15,2);default:0.00"`
	TaxAmount             float64               `gorm:"type:decimal(15,2);default:0.00"`
	TotalAmount           float64               `gorm:"type:decimal(15,2);not null;index:idx_sales_outlet_date_total"`
	PaidAmount            float64               `gorm:"type:decimal(15,2);not null"`
	ChangeAmount          float64               `gorm:"type:decimal(15,2);default:0.00"`
	PaymentMethod         PaymentMethodType     `gorm:"not null"`
	Status                TransactionStatusType `gorm:"default:'completed';index:idx_sales_tenant_date_status"`
	Notes                 string                `gorm:"type:text"`
	CreatedAt             time.Time             `gorm:"autoCreateTime"`
	UpdatedAt             time.Time             `gorm:"autoUpdateTime"`

	Tenant              Tenant               `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
	Outlet              Outlet               `gorm:"foreignKey:OutletID"`
	Cashier             User                 `gorm:"foreignKey:CashierID"`
	Customer            *Customer            `gorm:"foreignKey:CustomerID;constraint:OnDelete:SET NULL"`
	TransactionItems    []TransactionItem    `gorm:"foreignKey:TransactionID"`
	TransactionPayments []TransactionPayment `gorm:"foreignKey:TransactionID"`
}

func (SalesTransaction) TableName() string {
	return "transactions"
}

type TransactionItem struct {
	ID                      uint64  `gorm:"primaryKey;autoIncrement"`
	TransactionID           uint64  `gorm:"not null;index:idx_transaction_items_transaction"`
	ProductID               uint64  `gorm:"not null;index:idx_transaction_items_product"`
	ProductNameSnapshot     string  `gorm:"size:255;not null;index:idx_transaction_items_product_name_snapshot"`
	ProductSKUSnapshot      string  `gorm:"size:100;not null;index:idx_transaction_items_product_sku_snapshot"`
	ProductCategorySnapshot string  `gorm:"size:255"`
	ProductUnitSnapshot     string  `gorm:"size:50;default:'pcs'"`
	Quantity                int     `gorm:"not null"`
	UnitPrice               float64 `gorm:"type:decimal(12,2);not null"`
	CostPriceSnapshot       float64 `gorm:"type:decimal(12,2);default:0.00"`
	DiscountAmount          float64 `gorm:"type:decimal(12,2);default:0.00"`
	TotalPrice              float64 `gorm:"type:decimal(15,2);not null"`
	Notes                   string  `gorm:"type:text"`

	Transaction SalesTransaction `gorm:"foreignKey:TransactionID;constraint:OnDelete:CASCADE"`
	Product     Product          `gorm:"foreignKey:ProductID"`
}

type TransactionPayment struct {
	ID              uint64              `gorm:"primaryKey;autoIncrement"`
	TransactionID   uint64              `gorm:"not null;index:idx_transaction_payments_transaction"`
	PaymentMethod   PaymentMethodSingle `gorm:"not null"`
	Amount          float64             `gorm:"type:decimal(15,2);not null"`
	ReferenceNumber string              `gorm:"size:100"`
	Notes           string              `gorm:"type:text"`
	CreatedAt       time.Time           `gorm:"autoCreateTime"`

	Transaction SalesTransaction `gorm:"foreignKey:TransactionID;constraint:OnDelete:CASCADE"`
}
