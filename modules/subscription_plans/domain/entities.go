package domain

import (
	"time"
)

type SubscriptionPlan struct {
	ID                      uint64
	Name                    string
	Description             string
	Price                   float64
	MaxOutlets              int
	MaxUsers                int
	MaxProducts             *int
	MaxTransactionsPerMonth *int
	Features                []string
	IsActive                bool
	CreatedAt               time.Time
	UpdatedAt               time.Time
}