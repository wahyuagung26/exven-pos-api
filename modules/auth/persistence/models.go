package persistence

import (
	"time"

	"github.com/exven/pos-system/modules/auth/domain"
)

// UserModel maps to the database users table
type UserModel struct {
	ID              uint64 `gorm:"primaryKey;autoIncrement"`
	TenantID        uint64 `gorm:"not null;uniqueIndex:idx_tenant_email;index:idx_tenant_active"`
	RoleID          uint64 `gorm:"not null"`
	Email           string `gorm:"size:255;not null;uniqueIndex:idx_tenant_email"`
	PasswordHash    string `gorm:"size:255;not null"`
	FullName        string `gorm:"size:255;not null"`
	Phone           string `gorm:"size:20"`
	AvatarURL       string `gorm:"size:500"`
	IsActive        bool   `gorm:"default:true;index:idx_tenant_active"`
	LastLoginAt     *time.Time
	EmailVerifiedAt *time.Time
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`

	Role   RoleModel   `gorm:"foreignKey:RoleID"`
	Tenant TenantModel `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
}

func (UserModel) TableName() string {
	return "users"
}

// RoleModel maps to the database roles table
type RoleModel struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"size:50;uniqueIndex;not null"`
	DisplayName string    `gorm:"size:100;not null"`
	Description string    `gorm:"type:text"`
	Permissions string    `gorm:"type:json"`
	IsSystem    bool      `gorm:"default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (RoleModel) TableName() string {
	return "roles"
}

// TenantModel maps to the database tenants table
type TenantModel struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"size:255;not null"`
	BusinessType string `gorm:"size:100"`
	Email        string `gorm:"size:255;uniqueIndex;not null"`
	Phone        string `gorm:"size:20"`
	IsActive     bool   `gorm:"default:true;index"`
	TrialEndsAt  *time.Time
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (TenantModel) TableName() string {
	return "tenants"
}

// ToDomainUser converts UserModel to domain.User
func (u *UserModel) ToDomainUser() *domain.User {
	var role *domain.Role
	if u.Role.ID != 0 {
		role = u.Role.ToDomainRole()
	}

	var tenant *domain.Tenant
	if u.Tenant.ID != 0 {
		tenant = u.Tenant.ToDomainTenant()
	}

	return &domain.User{
		ID:              u.ID,
		TenantID:        u.TenantID,
		RoleID:          u.RoleID,
		Email:           u.Email,
		PasswordHash:    u.PasswordHash,
		FullName:        u.FullName,
		Phone:           u.Phone,
		AvatarURL:       u.AvatarURL,
		IsActive:        u.IsActive,
		LastLoginAt:     u.LastLoginAt,
		EmailVerifiedAt: u.EmailVerifiedAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
		Role:            role,
		Tenant:          tenant,
	}
}

// FromDomainUser converts domain.User to UserModel
func (u *UserModel) FromDomainUser(user *domain.User) {
	u.ID = user.ID
	u.TenantID = user.TenantID
	u.RoleID = user.RoleID
	u.Email = user.Email
	u.PasswordHash = user.PasswordHash
	u.FullName = user.FullName
	u.Phone = user.Phone
	u.AvatarURL = user.AvatarURL
	u.IsActive = user.IsActive
	u.LastLoginAt = user.LastLoginAt
	u.EmailVerifiedAt = user.EmailVerifiedAt
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
}

// ToDomainRole converts RoleModel to domain.Role
func (r *RoleModel) ToDomainRole() *domain.Role {
	return &domain.Role{
		ID:          r.ID,
		Name:        r.Name,
		DisplayName: r.DisplayName,
		Description: r.Description,
		Permissions: []string{}, // Parse JSON permissions if needed
		IsSystem:    r.IsSystem,
		CreatedAt:   r.CreatedAt,
	}
}

// ToDomainTenant converts TenantModel to domain.Tenant
func (t *TenantModel) ToDomainTenant() *domain.Tenant {
	return &domain.Tenant{
		ID:           t.ID,
		Name:         t.Name,
		BusinessType: t.BusinessType,
		Email:        t.Email,
		Phone:        t.Phone,
		IsActive:     t.IsActive,
		TrialEndsAt:  t.TrialEndsAt,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}
