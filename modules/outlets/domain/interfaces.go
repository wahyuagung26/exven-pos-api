package domain

import "context"

type OutletRepository interface {
	Create(ctx context.Context, outlet *Outlet) error
	GetByID(ctx context.Context, tenantID, outletID uint64) (*Outlet, error)
	GetByCode(ctx context.Context, tenantID uint64, code string) (*Outlet, error)
	GetAll(ctx context.Context, tenantID uint64, query OutletQuery) ([]*Outlet, int64, error)
	Update(ctx context.Context, outlet *Outlet) error
	Delete(ctx context.Context, tenantID, outletID uint64) error
	IsCodeExists(ctx context.Context, tenantID uint64, code string, excludeID *uint64) (bool, error)
}

type OutletService interface {
	Create(ctx context.Context, tenantID uint64, req CreateOutletRequest) (*Outlet, error)
	GetByID(ctx context.Context, tenantID, outletID uint64) (*Outlet, error)
	GetByCode(ctx context.Context, tenantID uint64, code string) (*Outlet, error)
	GetAll(ctx context.Context, tenantID uint64, query OutletQuery) ([]*Outlet, int64, error)
	Update(ctx context.Context, tenantID, outletID uint64, req UpdateOutletRequest) (*Outlet, error)
	Delete(ctx context.Context, tenantID, outletID uint64) error
}