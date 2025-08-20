package domain

import (
	"time"
)

type Outlet struct {
	ID          uint64
	TenantID    uint64
	Name        string
	Code        string
	Description string
	Address     string
	City        string
	Province    string
	PostalCode  string
	Phone       string
	Email       string
	ManagerID   *uint64
	IsActive    bool
	Settings    map[string]interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Manager *Manager
}

type Manager struct {
	ID       uint64
	FullName string
	Email    string
	Phone    string
}