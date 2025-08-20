package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type JSONPermissions []string

func (j JSONPermissions) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONPermissions) Scan(value interface{}) error {
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

type JSONSettings map[string]interface{}

func (j JSONSettings) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONSettings) Scan(value interface{}) error {
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

type Role struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"size:50;uniqueIndex;not null"`
	DisplayName string    `gorm:"size:100;not null"`
	Description string    `gorm:"type:text"`
	Permissions JSONPermissions `gorm:"type:jsonb"`
	IsSystem    bool      `gorm:"default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	Users []User `gorm:"foreignKey:RoleID"`
}

type User struct {
	ID               uint64     `gorm:"primaryKey;autoIncrement"`
	TenantID         uint64     `gorm:"not null;uniqueIndex:idx_tenant_username;uniqueIndex:idx_tenant_email;index:idx_tenant_active"`
	RoleID           uint64     `gorm:"not null"`
	Username         string     `gorm:"size:100;not null;uniqueIndex:idx_tenant_username"`
	Email            string     `gorm:"size:255;not null;uniqueIndex:idx_tenant_email"`
	PasswordHash     string     `gorm:"size:255;not null"`
	FullName         string     `gorm:"size:255;not null"`
	Phone            string     `gorm:"size:20"`
	AvatarURL        string     `gorm:"size:500"`
	IsActive         bool       `gorm:"default:true;index:idx_tenant_active"`
	LastLoginAt      *time.Time
	EmailVerifiedAt  *time.Time
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`

	Tenant      Tenant       `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
	Role        Role         `gorm:"foreignKey:RoleID"`
	UserOutlets []UserOutlet `gorm:"foreignKey:UserID"`
}

type Outlet struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	TenantID    uint64    `gorm:"not null;uniqueIndex:idx_tenant_code;index:idx_tenant_outlet_active"`
	Name        string    `gorm:"size:255;not null"`
	Code        string    `gorm:"size:50;not null;uniqueIndex:idx_tenant_code"`
	Description string    `gorm:"type:text"`
	Address     string    `gorm:"type:text"`
	City        string    `gorm:"size:100"`
	Province    string    `gorm:"size:100"`
	PostalCode  string    `gorm:"size:10"`
	Phone       string    `gorm:"size:20"`
	Email       string    `gorm:"size:255"`
	ManagerID   *uint64
	IsActive    bool      `gorm:"default:true;index:idx_tenant_outlet_active"`
	Settings    JSONSettings `gorm:"type:jsonb"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Tenant      Tenant       `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
	Manager     *User        `gorm:"foreignKey:ManagerID;constraint:OnDelete:SET NULL"`
	UserOutlets []UserOutlet `gorm:"foreignKey:OutletID"`
}

type UserOutlet struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	UserID    uint64    `gorm:"not null;uniqueIndex:idx_user_outlet"`
	OutletID  uint64    `gorm:"not null;uniqueIndex:idx_user_outlet;index:idx_outlet_active"`
	IsActive  bool      `gorm:"default:true;index:idx_outlet_active"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Outlet Outlet `gorm:"foreignKey:OutletID;constraint:OnDelete:CASCADE"`
}