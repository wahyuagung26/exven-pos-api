package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/exven/pos-system/modules/customers/domain"
)

type customerService struct {
	customerRepo domain.CustomerRepository
}

func NewCustomerService(customerRepo domain.CustomerRepository) domain.CustomerService {
	return &customerService{
		customerRepo: customerRepo,
	}
}

func (s *customerService) Create(ctx context.Context, tenantID uint64, req domain.CreateCustomerRequest) (*domain.Customer, error) {
	// Validate code uniqueness if provided
	if req.Code != "" {
		exists, err := s.customerRepo.IsCodeExists(ctx, tenantID, req.Code, nil)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("customer with this code already exists")
		}
	}

	// Validate phone uniqueness if provided
	if req.Phone != "" {
		exists, err := s.customerRepo.IsPhoneExists(ctx, tenantID, req.Phone, nil)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("customer with this phone already exists")
		}
	}

	// Validate email uniqueness if provided
	if req.Email != "" {
		exists, err := s.customerRepo.IsEmailExists(ctx, tenantID, req.Email, nil)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("customer with this email already exists")
		}
	}

	// Generate customer code if not provided
	code := strings.TrimSpace(req.Code)
	if code == "" {
		code = s.generateCustomerCode(strings.TrimSpace(req.Name))
		// Ensure generated code is unique
		for {
			exists, err := s.customerRepo.IsCodeExists(ctx, tenantID, code, nil)
			if err != nil {
				return nil, err
			}
			if !exists {
				break
			}
			code = s.generateCustomerCodeWithSuffix(strings.TrimSpace(req.Name))
		}
	}

	// Create customer entity
	customer := &domain.Customer{
		TenantID:      tenantID,
		Code:          code,
		Name:          strings.TrimSpace(req.Name),
		Email:         strings.TrimSpace(req.Email),
		Phone:         strings.TrimSpace(req.Phone),
		Address:       strings.TrimSpace(req.Address),
		City:          strings.TrimSpace(req.City),
		Province:      strings.TrimSpace(req.Province),
		PostalCode:    strings.TrimSpace(req.PostalCode),
		BirthDate:     req.BirthDate,
		Gender:        req.Gender,
		LoyaltyPoints: 0,
		TotalSpent:    0.0,
		VisitCount:    0,
		LastVisitAt:   nil,
		Notes:         strings.TrimSpace(req.Notes),
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := s.customerRepo.Create(ctx, customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *customerService) Update(ctx context.Context, tenantID, customerID uint64, req domain.UpdateCustomerRequest) (*domain.Customer, error) {
	// Check if customer exists
	existingCustomer, err := s.customerRepo.GetByID(ctx, tenantID, customerID)
	if err != nil {
		return nil, err
	}

	// Validate code uniqueness if provided (excluding current customer)
	if req.Code != "" {
		exists, err := s.customerRepo.IsCodeExists(ctx, tenantID, req.Code, &customerID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("customer with this code already exists")
		}
	}

	// Validate phone uniqueness if provided (excluding current customer)
	if req.Phone != "" {
		exists, err := s.customerRepo.IsPhoneExists(ctx, tenantID, req.Phone, &customerID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("customer with this phone already exists")
		}
	}

	// Validate email uniqueness if provided (excluding current customer)
	if req.Email != "" {
		exists, err := s.customerRepo.IsEmailExists(ctx, tenantID, req.Email, &customerID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("customer with this email already exists")
		}
	}

	// Update customer entity
	existingCustomer.Code = strings.TrimSpace(req.Code)
	existingCustomer.Name = strings.TrimSpace(req.Name)
	existingCustomer.Email = strings.TrimSpace(req.Email)
	existingCustomer.Phone = strings.TrimSpace(req.Phone)
	existingCustomer.Address = strings.TrimSpace(req.Address)
	existingCustomer.City = strings.TrimSpace(req.City)
	existingCustomer.Province = strings.TrimSpace(req.Province)
	existingCustomer.PostalCode = strings.TrimSpace(req.PostalCode)
	existingCustomer.BirthDate = req.BirthDate
	existingCustomer.Gender = req.Gender
	existingCustomer.Notes = strings.TrimSpace(req.Notes)
	existingCustomer.IsActive = req.IsActive
	existingCustomer.UpdatedAt = time.Now()

	err = s.customerRepo.Update(ctx, existingCustomer)
	if err != nil {
		return nil, err
	}

	return existingCustomer, nil
}

func (s *customerService) Delete(ctx context.Context, tenantID, customerID uint64) error {
	// Check if customer exists
	_, err := s.customerRepo.GetByID(ctx, tenantID, customerID)
	if err != nil {
		return err
	}

	return s.customerRepo.Delete(ctx, tenantID, customerID)
}

func (s *customerService) GetByID(ctx context.Context, tenantID, customerID uint64) (*domain.Customer, error) {
	return s.customerRepo.GetByID(ctx, tenantID, customerID)
}

func (s *customerService) GetByCode(ctx context.Context, tenantID uint64, code string) (*domain.Customer, error) {
	if strings.TrimSpace(code) == "" {
		return nil, errors.New("code cannot be empty")
	}

	return s.customerRepo.GetByCode(ctx, tenantID, strings.TrimSpace(code))
}

func (s *customerService) GetByPhone(ctx context.Context, tenantID uint64, phone string) (*domain.Customer, error) {
	if strings.TrimSpace(phone) == "" {
		return nil, errors.New("phone cannot be empty")
	}

	return s.customerRepo.GetByPhone(ctx, tenantID, strings.TrimSpace(phone))
}

func (s *customerService) GetByEmail(ctx context.Context, tenantID uint64, email string) (*domain.Customer, error) {
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("email cannot be empty")
	}

	return s.customerRepo.GetByEmail(ctx, tenantID, strings.TrimSpace(email))
}

func (s *customerService) GetAll(ctx context.Context, tenantID uint64, query domain.CustomerQuery) ([]*domain.Customer, int64, error) {
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

	return s.customerRepo.GetAll(ctx, tenantID, query)
}

// Helper functions

func (s *customerService) generateCustomerCode(name string) string {
	if name == "" {
		return fmt.Sprintf("CUST%d", time.Now().Unix())
	}

	// Take first 4 characters of name (uppercase) and append timestamp
	namePrefix := strings.ToUpper(name)
	if len(namePrefix) > 4 {
		namePrefix = namePrefix[:4]
	}
	
	return fmt.Sprintf("%s%d", namePrefix, time.Now().Unix()%10000)
}

func (s *customerService) generateCustomerCodeWithSuffix(name string) string {
	if name == "" {
		return fmt.Sprintf("CUST%d%d", time.Now().Unix(), time.Now().Nanosecond()%1000)
	}

	// Take first 4 characters of name (uppercase) and append timestamp with nanosecond
	namePrefix := strings.ToUpper(name)
	if len(namePrefix) > 4 {
		namePrefix = namePrefix[:4]
	}
	
	return fmt.Sprintf("%s%d%d", namePrefix, time.Now().Unix()%10000, time.Now().Nanosecond()%1000)
}