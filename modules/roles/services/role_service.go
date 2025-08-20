package services

import (
	"context"

	"github.com/exven/pos-system/modules/roles/domain"
)

type roleService struct {
	repo domain.RoleRepository
}

func NewRoleService(repo domain.RoleRepository) domain.RoleService {
	return &roleService{
		repo: repo,
	}
}

func (s *roleService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Role, int64, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *roleService) GetByID(ctx context.Context, id uint64) (*domain.Role, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *roleService) GetByName(ctx context.Context, name string) (*domain.Role, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *roleService) GetSystemRoles(ctx context.Context) ([]*domain.Role, error) {
	return s.repo.GetSystemRoles(ctx)
}