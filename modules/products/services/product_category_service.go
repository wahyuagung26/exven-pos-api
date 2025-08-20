package services

import (
	"context"
	"fmt"

	"github.com/exven/pos-system/modules/products/domain"
)

type ProductCategoryService struct {
	repo domain.ProductCategoryRepository
}

func NewProductCategoryService(repo domain.ProductCategoryRepository) *ProductCategoryService {
	return &ProductCategoryService{
		repo: repo,
	}
}

func (s *ProductCategoryService) Create(ctx context.Context, tenantID uint64, req domain.CreateProductCategoryRequest) (*domain.ProductCategory, error) {
	// Check if category name already exists
	exists, err := s.repo.CheckIfExists(ctx, tenantID, req.Name, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check category existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("category with name '%s' already exists", req.Name)
	}

	// If parent_id is specified, validate parent exists
	if req.ParentID != nil {
		parent, err := s.repo.FindByID(ctx, tenantID, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent category not found: %w", err)
		}
		if !parent.IsActive {
			return nil, fmt.Errorf("parent category is not active")
		}
	}

	category := &domain.ProductCategory{
		TenantID:    tenantID,
		ParentID:    req.ParentID,
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		SortOrder:   req.SortOrder,
		IsActive:    true,
	}

	err = s.repo.Create(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (s *ProductCategoryService) Update(ctx context.Context, tenantID, categoryID uint64, req domain.UpdateProductCategoryRequest) (*domain.ProductCategory, error) {
	// Check if category exists
	category, err := s.repo.FindByID(ctx, tenantID, categoryID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	// Check if new name conflicts with existing categories (exclude current)
	exists, err := s.repo.CheckIfExists(ctx, tenantID, req.Name, &categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to check category existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("category with name '%s' already exists", req.Name)
	}

	// If parent_id is specified and changed, validate parent exists
	if req.ParentID != nil {
		// Prevent self-reference
		if *req.ParentID == categoryID {
			return nil, fmt.Errorf("category cannot be its own parent")
		}

		// Check parent exists
		parent, err := s.repo.FindByID(ctx, tenantID, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent category not found: %w", err)
		}
		if !parent.IsActive {
			return nil, fmt.Errorf("parent category is not active")
		}

		// Prevent circular references by checking if new parent is a descendant
		if err := s.validateNoCircularReference(ctx, tenantID, categoryID, *req.ParentID); err != nil {
			return nil, err
		}
	}

	// Update category
	category.Name = req.Name
	category.Description = req.Description
	category.ImageURL = req.ImageURL
	category.ParentID = req.ParentID
	category.SortOrder = req.SortOrder
	category.IsActive = req.IsActive

	err = s.repo.Update(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

func (s *ProductCategoryService) Delete(ctx context.Context, tenantID, categoryID uint64) error {
	// Check if category exists
	_, err := s.repo.FindByID(ctx, tenantID, categoryID)
	if err != nil {
		return fmt.Errorf("category not found: %w", err)
	}

	// Check if category has subcategories
	hasChildren, err := s.repo.HasChildren(ctx, tenantID, categoryID)
	if err != nil {
		return fmt.Errorf("failed to check for subcategories: %w", err)
	}
	if hasChildren {
		return fmt.Errorf("cannot delete category that has subcategories")
	}

	// TODO: Check if category has products
	// This would require implementing product repository first

	err = s.repo.Delete(ctx, tenantID, categoryID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

func (s *ProductCategoryService) GetByID(ctx context.Context, tenantID, categoryID uint64) (*domain.ProductCategory, error) {
	category, err := s.repo.FindByID(ctx, tenantID, categoryID)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	// Load subcategories
	subcategories, err := s.repo.FindByParentID(ctx, tenantID, &categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to load subcategories: %w", err)
	}
	category.SubCategories = subcategories

	return category, nil
}

func (s *ProductCategoryService) GetAll(ctx context.Context, tenantID uint64, query domain.ProductCategoryQuery) ([]*domain.ProductCategory, int64, error) {
	categories, total, err := s.repo.FindAll(ctx, tenantID, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, total, nil
}

func (s *ProductCategoryService) GetHierarchy(ctx context.Context, tenantID uint64) ([]*domain.ProductCategory, error) {
	// Get all root categories (no parent)
	rootCategories, err := s.repo.FindByParentID(ctx, tenantID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get root categories: %w", err)
	}

	// For each root category, load its hierarchy
	for _, category := range rootCategories {
		if err := s.loadCategoryHierarchy(ctx, tenantID, category); err != nil {
			return nil, err
		}
	}

	return rootCategories, nil
}

// Helper function to validate no circular reference
func (s *ProductCategoryService) validateNoCircularReference(ctx context.Context, tenantID, categoryID, newParentID uint64) error {
	// Get all descendants of the category being updated
	descendants, err := s.getAllDescendants(ctx, tenantID, categoryID)
	if err != nil {
		return err
	}

	// Check if newParentID is in the descendants
	for _, descendant := range descendants {
		if descendant.ID == newParentID {
			return fmt.Errorf("circular reference detected: new parent is a descendant of this category")
		}
	}

	return nil
}

// Helper function to get all descendants recursively
func (s *ProductCategoryService) getAllDescendants(ctx context.Context, tenantID, categoryID uint64) ([]*domain.ProductCategory, error) {
	var descendants []*domain.ProductCategory

	children, err := s.repo.FindByParentID(ctx, tenantID, &categoryID)
	if err != nil {
		return nil, err
	}

	for _, child := range children {
		descendants = append(descendants, child)

		// Recursively get grandchildren
		grandchildren, err := s.getAllDescendants(ctx, tenantID, child.ID)
		if err != nil {
			return nil, err
		}
		descendants = append(descendants, grandchildren...)
	}

	return descendants, nil
}

// Helper function to recursively load category hierarchy
func (s *ProductCategoryService) loadCategoryHierarchy(ctx context.Context, tenantID uint64, category *domain.ProductCategory) error {
	subcategories, err := s.repo.FindByParentID(ctx, tenantID, &category.ID)
	if err != nil {
		return err
	}

	category.SubCategories = subcategories

	// Recursively load subcategories
	for _, subcategory := range subcategories {
		if err := s.loadCategoryHierarchy(ctx, tenantID, subcategory); err != nil {
			return err
		}
	}

	return nil
}
