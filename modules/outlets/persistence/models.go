package persistence

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/exven/pos-system/modules/outlets/domain"
	"time"
)

type JSONSettingsModel map[string]interface{}

func (j JSONSettingsModel) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONSettingsModel) Scan(value interface{}) error {
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

type OutletModel struct {
	ID          uint64             `gorm:"primaryKey;autoIncrement"`
	TenantID    uint64             `gorm:"not null;uniqueIndex:idx_tenant_code;index:idx_tenant_outlet_active"`
	Name        string             `gorm:"size:255;not null"`
	Code        string             `gorm:"size:50;not null;uniqueIndex:idx_tenant_code"`
	Description string             `gorm:"type:text"`
	Address     string             `gorm:"type:text"`
	City        string             `gorm:"size:100"`
	Province    string             `gorm:"size:100"`
	PostalCode  string             `gorm:"size:10"`
	Phone       string             `gorm:"size:20"`
	Email       string             `gorm:"size:255"`
	ManagerID   *uint64            `gorm:"column:manager_id"`
	IsActive    bool               `gorm:"default:true;index:idx_tenant_outlet_active"`
	Settings    JSONSettingsModel  `gorm:"type:jsonb"`
	CreatedAt   time.Time          `gorm:"autoCreateTime"`
	UpdatedAt   time.Time          `gorm:"autoUpdateTime"`
}

func (OutletModel) TableName() string {
	return "outlets"
}

type OutletWithManagerModel struct {
	OutletModel
	ManagerFullName *string `gorm:"column:manager_full_name"`
	ManagerEmail    *string `gorm:"column:manager_email"`
	ManagerPhone    *string `gorm:"column:manager_phone"`
}

// Mapper functions

func (o *OutletModel) ToDomainOutlet() *domain.Outlet {
	settings := make(map[string]interface{})
	for k, v := range o.Settings {
		settings[k] = v
	}

	return &domain.Outlet{
		ID:          o.ID,
		TenantID:    o.TenantID,
		Name:        o.Name,
		Code:        o.Code,
		Description: o.Description,
		Address:     o.Address,
		City:        o.City,
		Province:    o.Province,
		PostalCode:  o.PostalCode,
		Phone:       o.Phone,
		Email:       o.Email,
		ManagerID:   o.ManagerID,
		IsActive:    o.IsActive,
		Settings:    settings,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

func (o *OutletModel) FromDomainOutlet(outlet *domain.Outlet) {
	o.ID = outlet.ID
	o.TenantID = outlet.TenantID
	o.Name = outlet.Name
	o.Code = outlet.Code
	o.Description = outlet.Description
	o.Address = outlet.Address
	o.City = outlet.City
	o.Province = outlet.Province
	o.PostalCode = outlet.PostalCode
	o.Phone = outlet.Phone
	o.Email = outlet.Email
	o.ManagerID = outlet.ManagerID
	o.IsActive = outlet.IsActive
	o.CreatedAt = outlet.CreatedAt
	o.UpdatedAt = outlet.UpdatedAt

	settings := make(JSONSettingsModel)
	for k, v := range outlet.Settings {
		settings[k] = v
	}
	o.Settings = settings
}

func (o *OutletWithManagerModel) ToDomainOutletWithManager() *domain.Outlet {
	outlet := o.OutletModel.ToDomainOutlet()
	
	if o.ManagerFullName != nil {
		outlet.Manager = &domain.Manager{
			ID:       *o.ManagerID,
			FullName: *o.ManagerFullName,
			Email:    "",
			Phone:    "",
		}
		
		if o.ManagerEmail != nil {
			outlet.Manager.Email = *o.ManagerEmail
		}
		
		if o.ManagerPhone != nil {
			outlet.Manager.Phone = *o.ManagerPhone
		}
	}
	
	return outlet
}