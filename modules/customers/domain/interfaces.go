package domain

import "context"

type CustomerRepository interface {
	Create(ctx context.Context, customer *Customer) error
	GetByID(ctx context.Context, tenantID, customerID uint64) (*Customer, error)
	GetByCode(ctx context.Context, tenantID uint64, code string) (*Customer, error)
	GetByPhone(ctx context.Context, tenantID uint64, phone string) (*Customer, error)
	GetByEmail(ctx context.Context, tenantID uint64, email string) (*Customer, error)
	GetAll(ctx context.Context, tenantID uint64, query CustomerQuery) ([]*Customer, int64, error)
	Update(ctx context.Context, customer *Customer) error
	Delete(ctx context.Context, tenantID, customerID uint64) error
	IsCodeExists(ctx context.Context, tenantID uint64, code string, excludeID *uint64) (bool, error)
	IsPhoneExists(ctx context.Context, tenantID uint64, phone string, excludeID *uint64) (bool, error)
	IsEmailExists(ctx context.Context, tenantID uint64, email string, excludeID *uint64) (bool, error)
	UpdateStats(ctx context.Context, customerID uint64, totalSpent float64, visitCount int) error
}

type CustomerService interface {
	Create(ctx context.Context, tenantID uint64, req CreateCustomerRequest) (*Customer, error)
	GetByID(ctx context.Context, tenantID, customerID uint64) (*Customer, error)
	GetByCode(ctx context.Context, tenantID uint64, code string) (*Customer, error)
	GetByPhone(ctx context.Context, tenantID uint64, phone string) (*Customer, error)
	GetByEmail(ctx context.Context, tenantID uint64, email string) (*Customer, error)
	GetAll(ctx context.Context, tenantID uint64, query CustomerQuery) ([]*Customer, int64, error)
	Update(ctx context.Context, tenantID, customerID uint64, req UpdateCustomerRequest) (*Customer, error)
	Delete(ctx context.Context, tenantID, customerID uint64) error
}