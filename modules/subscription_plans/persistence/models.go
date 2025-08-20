package persistence

import (
	"time"
	"encoding/json"
)

type SubscriptionPlanModel struct {
	ID                      uint64    `gorm:"primaryKey;column:id"`
	Name                    string    `gorm:"column:name"`
	Description             string    `gorm:"column:description"`
	Price                   float64   `gorm:"column:price"`
	MaxOutlets              int       `gorm:"column:max_outlets"`
	MaxUsers                int       `gorm:"column:max_users"`
	MaxProducts             *int      `gorm:"column:max_products"`
	MaxTransactionsPerMonth *int      `gorm:"column:max_transactions_per_month"`
	Features                string    `gorm:"column:features"` // JSON string
	IsActive                bool      `gorm:"column:is_active"`
	CreatedAt               time.Time `gorm:"column:created_at"`
	UpdatedAt               time.Time `gorm:"column:updated_at"`
}

func (SubscriptionPlanModel) TableName() string {
	return "subscription_plans"
}

func (m *SubscriptionPlanModel) GetFeaturesSlice() []string {
	if m.Features == "" {
		return []string{}
	}
	
	var features []string
	if err := json.Unmarshal([]byte(m.Features), &features); err != nil {
		return []string{}
	}
	
	return features
}