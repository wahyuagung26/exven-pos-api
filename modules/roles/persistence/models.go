package persistence

import (
	"time"
	"encoding/json"
)

type RoleModel struct {
	ID          uint64    `gorm:"primaryKey;column:id"`
	Name        string    `gorm:"column:name;unique"`
	DisplayName string    `gorm:"column:display_name"`
	Description string    `gorm:"column:description"`
	Permissions string    `gorm:"column:permissions"` // JSON string
	IsSystem    bool      `gorm:"column:is_system"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (RoleModel) TableName() string {
	return "roles"
}

func (m *RoleModel) GetPermissionsSlice() []string {
	if m.Permissions == "" {
		return []string{}
	}
	
	var permissions []string
	if err := json.Unmarshal([]byte(m.Permissions), &permissions); err != nil {
		return []string{}
	}
	
	return permissions
}