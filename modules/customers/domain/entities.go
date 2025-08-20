package domain

import (
	"time"
)

type Customer struct {
	ID            uint64
	TenantID      uint64
	Code          string
	Name          string
	Email         string
	Phone         string
	Address       string
	City          string
	Province      string
	PostalCode    string
	BirthDate     *time.Time
	Gender        string
	LoyaltyPoints int
	TotalSpent    float64
	VisitCount    int
	LastVisitAt   *time.Time
	Notes         string
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}