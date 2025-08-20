package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/exven/pos-system/modules/products/domain"
)

type productService struct {
	productRepo  domain.ProductRepository
	categoryRepo domain.ProductCategoryRepository
}

func NewProductService(productRepo domain.ProductRepository, categoryRepo domain.ProductCategoryRepository) domain.ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *productService) Create(ctx context.Context, tenantID uint64, req domain.CreateProductRequest) (*domain.Product, error) {
	// Validate SKU uniqueness
	exists, err := s.productRepo.CheckSKUExists(ctx, tenantID, req.SKU, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("product with this SKU already exists")
	}

	// Validate category exists if provided
	if req.CategoryID != nil {
		_, err := s.categoryRepo.FindByID(ctx, tenantID, *req.CategoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
	}

	// Create product entity
	product := &domain.Product{
		TenantID:     tenantID,
		CategoryID:   req.CategoryID,
		SKU:          strings.TrimSpace(req.SKU),
		Barcode:      strings.TrimSpace(req.Barcode),
		Name:         strings.TrimSpace(req.Name),
		Description:  strings.TrimSpace(req.Description),
		Unit:         req.Unit,
		CostPrice:    req.CostPrice,
		SellingPrice: req.SellingPrice,
		MinStock:     req.MinStock,
		TrackStock:   req.TrackStock,
		IsActive:     true,
		Images:       req.Images,
		Variants:     req.Variants,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Set default unit if empty
	if product.Unit == "" {
		product.Unit = "pcs"
	}

	// Initialize slices/maps if nil
	if product.Images == nil {
		product.Images = []string{}
	}
	if product.Variants == nil {
		product.Variants = make(map[string]interface{})
	}

	err = s.productRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	// Return product with category information if exists
	return s.productRepo.FindByID(ctx, tenantID, product.ID)
}

func (s *productService) Update(ctx context.Context, tenantID, productID uint64, req domain.UpdateProductRequest) (*domain.Product, error) {
	// Check if product exists
	existingProduct, err := s.productRepo.FindByID(ctx, tenantID, productID)
	if err != nil {
		return nil, err
	}

	// Validate SKU uniqueness (excluding current product)
	exists, err := s.productRepo.CheckSKUExists(ctx, tenantID, req.SKU, &productID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("product with this SKU already exists")
	}

	// Validate category exists if provided
	if req.CategoryID != nil {
		_, err := s.categoryRepo.FindByID(ctx, tenantID, *req.CategoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
	}

	// Update product entity
	existingProduct.CategoryID = req.CategoryID
	existingProduct.SKU = strings.TrimSpace(req.SKU)
	existingProduct.Barcode = strings.TrimSpace(req.Barcode)
	existingProduct.Name = strings.TrimSpace(req.Name)
	existingProduct.Description = strings.TrimSpace(req.Description)
	existingProduct.Unit = req.Unit
	existingProduct.CostPrice = req.CostPrice
	existingProduct.SellingPrice = req.SellingPrice
	existingProduct.MinStock = req.MinStock
	existingProduct.TrackStock = req.TrackStock
	existingProduct.IsActive = req.IsActive
	existingProduct.Images = req.Images
	existingProduct.Variants = req.Variants
	existingProduct.UpdatedAt = time.Now()

	// Set default unit if empty
	if existingProduct.Unit == "" {
		existingProduct.Unit = "pcs"
	}

	// Initialize slices/maps if nil
	if existingProduct.Images == nil {
		existingProduct.Images = []string{}
	}
	if existingProduct.Variants == nil {
		existingProduct.Variants = make(map[string]interface{})
	}

	err = s.productRepo.Update(ctx, existingProduct)
	if err != nil {
		return nil, err
	}

	// Return updated product with category information
	return s.productRepo.FindByID(ctx, tenantID, productID)
}

func (s *productService) Delete(ctx context.Context, tenantID, productID uint64) error {
	// Check if product exists
	_, err := s.productRepo.FindByID(ctx, tenantID, productID)
	if err != nil {
		return err
	}

	return s.productRepo.Delete(ctx, tenantID, productID)
}

func (s *productService) GetByID(ctx context.Context, tenantID, productID uint64) (*domain.Product, error) {
	return s.productRepo.FindByID(ctx, tenantID, productID)
}

func (s *productService) GetAll(ctx context.Context, tenantID uint64, limit, offset int) ([]*domain.Product, int64, error) {
	// Set default pagination if not provided
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.productRepo.FindAll(ctx, tenantID, limit, offset)
}

func (s *productService) GetBySKU(ctx context.Context, tenantID uint64, sku string) (*domain.Product, error) {
	if strings.TrimSpace(sku) == "" {
		return nil, errors.New("SKU cannot be empty")
	}

	return s.productRepo.FindBySKU(ctx, tenantID, strings.TrimSpace(sku))
}

func (s *productService) GetByBarcode(ctx context.Context, tenantID uint64, barcode string) (*domain.Product, error) {
	if strings.TrimSpace(barcode) == "" {
		return nil, errors.New("barcode cannot be empty")
	}

	return s.productRepo.FindByBarcode(ctx, tenantID, strings.TrimSpace(barcode))
}

func (s *productService) GetByCategory(ctx context.Context, tenantID, categoryID uint64, limit, offset int) ([]*domain.Product, int64, error) {
	// Validate category exists
	_, err := s.categoryRepo.FindByID(ctx, tenantID, categoryID)
	if err != nil {
		return nil, 0, errors.New("category not found")
	}

	// Set default pagination if not provided
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return s.productRepo.FindByCategory(ctx, tenantID, categoryID, limit, offset)
}
