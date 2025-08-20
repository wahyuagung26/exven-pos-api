package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/exven/pos-system/modules/outlets/domain"
)

type outletService struct {
	outletRepo domain.OutletRepository
}

func NewOutletService(outletRepo domain.OutletRepository) domain.OutletService {
	return &outletService{
		outletRepo: outletRepo,
	}
}

func (s *outletService) Create(ctx context.Context, tenantID uint64, req domain.CreateOutletRequest) (*domain.Outlet, error) {
	// Validate code uniqueness
	exists, err := s.outletRepo.IsCodeExists(ctx, tenantID, req.Code, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("outlet with this code already exists")
	}

	// Create outlet entity
	outlet := &domain.Outlet{
		TenantID:    tenantID,
		Name:        strings.TrimSpace(req.Name),
		Code:        strings.TrimSpace(req.Code),
		Description: strings.TrimSpace(req.Description),
		Address:     strings.TrimSpace(req.Address),
		City:        strings.TrimSpace(req.City),
		Province:    strings.TrimSpace(req.Province),
		PostalCode:  strings.TrimSpace(req.PostalCode),
		Phone:       strings.TrimSpace(req.Phone),
		Email:       strings.TrimSpace(req.Email),
		ManagerID:   req.ManagerID,
		IsActive:    true,
		Settings:    req.Settings,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Initialize settings map if nil
	if outlet.Settings == nil {
		outlet.Settings = make(map[string]interface{})
	}

	err = s.outletRepo.Create(ctx, outlet)
	if err != nil {
		return nil, err
	}

	// Return outlet with manager information if exists
	return s.outletRepo.GetByID(ctx, tenantID, outlet.ID)
}

func (s *outletService) Update(ctx context.Context, tenantID, outletID uint64, req domain.UpdateOutletRequest) (*domain.Outlet, error) {
	// Check if outlet exists
	existingOutlet, err := s.outletRepo.GetByID(ctx, tenantID, outletID)
	if err != nil {
		return nil, err
	}

	// Validate code uniqueness (excluding current outlet)
	exists, err := s.outletRepo.IsCodeExists(ctx, tenantID, req.Code, &outletID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("outlet with this code already exists")
	}

	// Update outlet entity
	existingOutlet.Name = strings.TrimSpace(req.Name)
	existingOutlet.Code = strings.TrimSpace(req.Code)
	existingOutlet.Description = strings.TrimSpace(req.Description)
	existingOutlet.Address = strings.TrimSpace(req.Address)
	existingOutlet.City = strings.TrimSpace(req.City)
	existingOutlet.Province = strings.TrimSpace(req.Province)
	existingOutlet.PostalCode = strings.TrimSpace(req.PostalCode)
	existingOutlet.Phone = strings.TrimSpace(req.Phone)
	existingOutlet.Email = strings.TrimSpace(req.Email)
	existingOutlet.ManagerID = req.ManagerID
	existingOutlet.IsActive = req.IsActive
	existingOutlet.Settings = req.Settings
	existingOutlet.UpdatedAt = time.Now()

	// Initialize settings map if nil
	if existingOutlet.Settings == nil {
		existingOutlet.Settings = make(map[string]interface{})
	}

	err = s.outletRepo.Update(ctx, existingOutlet)
	if err != nil {
		return nil, err
	}

	// Return updated outlet with manager information
	return s.outletRepo.GetByID(ctx, tenantID, outletID)
}

func (s *outletService) Delete(ctx context.Context, tenantID, outletID uint64) error {
	// Check if outlet exists
	_, err := s.outletRepo.GetByID(ctx, tenantID, outletID)
	if err != nil {
		return err
	}

	return s.outletRepo.Delete(ctx, tenantID, outletID)
}

func (s *outletService) GetByID(ctx context.Context, tenantID, outletID uint64) (*domain.Outlet, error) {
	return s.outletRepo.GetByID(ctx, tenantID, outletID)
}

func (s *outletService) GetByCode(ctx context.Context, tenantID uint64, code string) (*domain.Outlet, error) {
	if strings.TrimSpace(code) == "" {
		return nil, errors.New("code cannot be empty")
	}

	return s.outletRepo.GetByCode(ctx, tenantID, strings.TrimSpace(code))
}

func (s *outletService) GetAll(ctx context.Context, tenantID uint64, query domain.OutletQuery) ([]*domain.Outlet, int64, error) {
	// Set default pagination if not provided
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 50
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	return s.outletRepo.GetAll(ctx, tenantID, query)
}