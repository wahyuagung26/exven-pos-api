package database

import (
	"time"
)

type GenderType string

const (
	GenderMale   GenderType = "male"
	GenderFemale GenderType = "female"
)

type Customer struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement"`
	TenantID      uint64     `gorm:"not null;uniqueIndex:idx_customers_tenant_code;index:idx_customers_tenant_phone;index:idx_customers_tenant_email"`
	Code          string     `gorm:"size:50;uniqueIndex:idx_customers_tenant_code"`
	Name          string     `gorm:"size:255;not null"`
	Email         string     `gorm:"size:255;index:idx_customers_tenant_email"`
	Phone         string     `gorm:"size:20;index:idx_customers_tenant_phone"`
	Address       string     `gorm:"type:text"`
	City          string     `gorm:"size:100"`
	Province      string     `gorm:"size:100"`
	PostalCode    string     `gorm:"size:10"`
	BirthDate     *time.Time `gorm:"type:date"`
	Gender        GenderType `gorm:"type:gender_type"`
	LoyaltyPoints int        `gorm:"default:0"`
	TotalSpent    float64    `gorm:"type:decimal(15,2);default:0.00"`
	VisitCount    int        `gorm:"default:0"`
	LastVisitAt   *time.Time
	Notes         string    `gorm:"type:text"`
	IsActive      bool      `gorm:"default:true"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`

	Tenant       Tenant             `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
	Transactions []SalesTransaction `gorm:"foreignKey:CustomerID"`
}
