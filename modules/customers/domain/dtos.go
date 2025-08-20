package domain

import "time"

type CreateCustomerRequest struct {
	Code       string     `json:"code" validate:"omitempty,max=50"`
	Name       string     `json:"name" validate:"required,min=1,max=255"`
	Email      string     `json:"email" validate:"omitempty,email,max=255"`
	Phone      string     `json:"phone" validate:"omitempty,max=20"`
	Address    string     `json:"address"`
	City       string     `json:"city" validate:"max=100"`
	Province   string     `json:"province" validate:"max=100"`
	PostalCode string     `json:"postal_code" validate:"max=10"`
	BirthDate  *time.Time `json:"birth_date"`
	Gender     string     `json:"gender" validate:"omitempty,oneof=male female"`
	Notes      string     `json:"notes"`
}

type UpdateCustomerRequest struct {
	Code       string     `json:"code" validate:"omitempty,max=50"`
	Name       string     `json:"name" validate:"required,min=1,max=255"`
	Email      string     `json:"email" validate:"omitempty,email,max=255"`
	Phone      string     `json:"phone" validate:"omitempty,max=20"`
	Address    string     `json:"address"`
	City       string     `json:"city" validate:"max=100"`
	Province   string     `json:"province" validate:"max=100"`
	PostalCode string     `json:"postal_code" validate:"max=10"`
	BirthDate  *time.Time `json:"birth_date"`
	Gender     string     `json:"gender" validate:"omitempty,oneof=male female"`
	Notes      string     `json:"notes"`
	IsActive   bool       `json:"is_active"`
}

type CustomerResponse struct {
	ID            uint64     `json:"id"`
	TenantID      uint64     `json:"tenant_id"`
	Code          string     `json:"code"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	Phone         string     `json:"phone"`
	Address       string     `json:"address"`
	City          string     `json:"city"`
	Province      string     `json:"province"`
	PostalCode    string     `json:"postal_code"`
	BirthDate     *time.Time `json:"birth_date"`
	Gender        string     `json:"gender"`
	LoyaltyPoints int        `json:"loyalty_points"`
	TotalSpent    float64    `json:"total_spent"`
	VisitCount    int        `json:"visit_count"`
	LastVisitAt   *time.Time `json:"last_visit_at"`
	Notes         string     `json:"notes"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     string     `json:"created_at"`
	UpdatedAt     string     `json:"updated_at"`
}

type CustomerListResponse struct {
	Customers []CustomerResponse `json:"customers"`
	Total     int64              `json:"total"`
	Page      int                `json:"page"`
	Limit     int                `json:"limit"`
}

type CustomerQuery struct {
	Name       string  `query:"name"`
	Code       string  `query:"code"`
	Email      string  `query:"email"`
	Phone      string  `query:"phone"`
	City       string  `query:"city"`
	Province   string  `query:"province"`
	Gender     string  `query:"gender"`
	IsActive   *bool   `query:"is_active"`
	Page       int     `query:"page"`
	Limit      int     `query:"limit"`
	Sort       string  `query:"sort"`
	Order      string  `query:"order"`
}