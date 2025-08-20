package domain

import (
	"time"
)

type User struct {
	ID              uint64
	TenantID        uint64
	RoleID          uint64
	Email           string
	PasswordHash    string
	FullName        string
	Phone           string
	AvatarURL       string
	IsActive        bool
	LastLoginAt     *time.Time
	EmailVerifiedAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Role   *Role
	Tenant *Tenant
}

type Role struct {
	ID          uint64
	Name        string
	DisplayName string
	Description string
	Permissions []string
	IsSystem    bool
	CreatedAt   time.Time
}

type Tenant struct {
	ID           uint64
	Name         string
	BusinessType string
	Email        string
	Phone        string
	IsActive     bool
	TrialEndsAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Session struct {
	ID           string
	UserID       uint64
	TenantID     uint64
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

type LoginCredentials struct {
	Email    string
	Password string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

type TokenClaims struct {
	UserID    uint64
	TenantID  uint64
	Email     string
	RoleID    uint64
	ExpiresAt time.Time
}
