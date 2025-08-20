package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
	SubscriptionStatusPending   SubscriptionStatus = "pending"
)

type JSONFeatures []string

func (j JSONFeatures) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONFeatures) Scan(value interface{}) error {
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

type SubscriptionPlan struct {
	ID                      uint64       `gorm:"primaryKey;autoIncrement"`
	Name                    string       `gorm:"size:100;not null"`
	Description             string       `gorm:"type:text"`
	Price                   float64      `gorm:"type:decimal(12,2);default:0.00"`
	MaxOutlets              int          `gorm:"not null;default:1"`
	MaxUsers                int          `gorm:"not null;default:1"`
	MaxProducts             *int         `gorm:"default:null"`
	MaxTransactionsPerMonth *int         `gorm:"default:null"`
	Features                JSONFeatures `gorm:"type:jsonb"`
	IsActive                bool         `gorm:"default:true"`
	CreatedAt               time.Time    `gorm:"autoCreateTime"`
	UpdatedAt               time.Time    `gorm:"autoUpdateTime"`
}

func (SubscriptionPlan) TableName() string {
	return "subscription_plans"
}

type Tenant struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement"`
	Name       string     `gorm:"size:255;not null"`
	BusinessType string   `gorm:"size:100"`
	Email      string     `gorm:"size:255;uniqueIndex;not null"`
	Phone      string     `gorm:"size:20"`
	Address    string     `gorm:"type:text"`
	City       string     `gorm:"size:100"`
	Province   string     `gorm:"size:100"`
	PostalCode string     `gorm:"size:10"`
	TaxNumber  string     `gorm:"size:50"`
	LogoURL    string     `gorm:"size:500"`
	Timezone   string     `gorm:"size:50;default:'Asia/Jakarta'"`
	Currency   string     `gorm:"size:3;default:'IDR'"`
	IsActive   bool       `gorm:"default:true;index"`
	TrialEndsAt *time.Time
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	Subscriptions []TenantSubscription `gorm:"foreignKey:TenantID"`
	Users         []User               `gorm:"foreignKey:TenantID"`
	Outlets       []Outlet             `gorm:"foreignKey:TenantID"`
}

type TenantSubscription struct {
	ID                 uint64    `gorm:"primaryKey;autoIncrement"`
	TenantID           uint64    `gorm:"not null;index:idx_tenant_status"`
	SubscriptionPlanID uint64    `gorm:"not null"`
	Status             SubscriptionStatus `gorm:"default:'pending';index:idx_tenant_status"`
	StartsAt           time.Time `gorm:"not null"`
	EndsAt             time.Time `gorm:"not null;index"`
	AutoRenew          bool      `gorm:"default:true"`
	PaymentMethod      string    `gorm:"size:50"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`

	Tenant           Tenant           `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
	SubscriptionPlan SubscriptionPlan `gorm:"foreignKey:SubscriptionPlanID"`
}