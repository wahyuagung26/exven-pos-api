package domain

import (
	"context"
)

type RoleRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]*Role, int64, error)
	GetByID(ctx context.Context, id uint64) (*Role, error)
	GetByName(ctx context.Context, name string) (*Role, error)
	GetSystemRoles(ctx context.Context) ([]*Role, error)
}

type RoleService interface {
	GetAll(ctx context.Context, limit, offset int) ([]*Role, int64, error)
	GetByID(ctx context.Context, id uint64) (*Role, error)
	GetByName(ctx context.Context, name string) (*Role, error)
	GetSystemRoles(ctx context.Context) ([]*Role, error)
}