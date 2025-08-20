package persistence

import (
	"time"

	"github.com/exven/pos-system/modules/customers/domain"
)

type CustomerModel struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement"`
	TenantID      uint64     `gorm:"not null;uniqueIndex:idx_tenant_code;index:idx_tenant_customer_active"`
	Code          string     `gorm:"size:50;uniqueIndex:idx_tenant_code"`
	Name          string     `gorm:"size:255;not null;index:idx_customer_name"`
	Email         string     `gorm:"size:255;index:idx_customer_email"`
	Phone         string     `gorm:"size:20;index:idx_customer_phone"`
	Address       string     `gorm:"type:text"`
	City          string     `gorm:"size:100"`
	Province      string     `gorm:"size:100"`
	PostalCode    string     `gorm:"size:10"`
	BirthDate     *time.Time `gorm:"column:birth_date"`
	Gender        string     `gorm:"size:10"`
	LoyaltyPoints int        `gorm:"default:0"`
	TotalSpent    float64    `gorm:"type:decimal(15,2);default:0.00"`
	VisitCount    int        `gorm:"default:0"`
	LastVisitAt   *time.Time `gorm:"column:last_visit_at"`
	Notes         string     `gorm:"type:text"`
	IsActive      bool       `gorm:"default:true;index:idx_tenant_customer_active"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime"`
}

func (CustomerModel) TableName() string {
	return "customers"
}

// Mapper functions

func (c *CustomerModel) ToDomainCustomer() *domain.Customer {
	return &domain.Customer{
		ID:            c.ID,
		TenantID:      c.TenantID,
		Code:          c.Code,
		Name:          c.Name,
		Email:         c.Email,
		Phone:         c.Phone,
		Address:       c.Address,
		City:          c.City,
		Province:      c.Province,
		PostalCode:    c.PostalCode,
		BirthDate:     c.BirthDate,
		Gender:        c.Gender,
		LoyaltyPoints: c.LoyaltyPoints,
		TotalSpent:    c.TotalSpent,
		VisitCount:    c.VisitCount,
		LastVisitAt:   c.LastVisitAt,
		Notes:         c.Notes,
		IsActive:      c.IsActive,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}

func (c *CustomerModel) FromDomainCustomer(customer *domain.Customer) {
	c.ID = customer.ID
	c.TenantID = customer.TenantID
	c.Code = customer.Code
	c.Name = customer.Name
	c.Email = customer.Email
	c.Phone = customer.Phone
	c.Address = customer.Address
	c.City = customer.City
	c.Province = customer.Province
	c.PostalCode = customer.PostalCode
	c.BirthDate = customer.BirthDate
	c.Gender = customer.Gender
	c.LoyaltyPoints = customer.LoyaltyPoints
	c.TotalSpent = customer.TotalSpent
	c.VisitCount = customer.VisitCount
	c.LastVisitAt = customer.LastVisitAt
	c.Notes = customer.Notes
	c.IsActive = customer.IsActive
	c.CreatedAt = customer.CreatedAt
	c.UpdatedAt = customer.UpdatedAt
}