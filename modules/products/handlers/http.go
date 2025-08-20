package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/exven/pos-system/modules/products/domain"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	categoryService domain.ProductCategoryService
}

func NewProductHandler(categoryService domain.ProductCategoryService) *ProductHandler {
	return &ProductHandler{
		categoryService: categoryService,
	}
}

func (h *ProductHandler) RegisterRoutes(e *echo.Group) {
	products := e.Group("/products")

	// Product Categories routes
	categories := products.Group("/categories")
	categories.POST("", h.CreateCategory)
	categories.GET("", h.GetCategories)
	categories.GET("/:id", h.GetCategory)
	categories.PUT("/:id", h.UpdateCategory)
	categories.DELETE("/:id", h.DeleteCategory)
	categories.GET("/hierarchy", h.GetCategoryHierarchy)
}

// Product Category handlers

func (h *ProductHandler) CreateCategory(c echo.Context) error {
	var req domain.CreateProductCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	tenantID := c.Get("tenant_id").(uint64)

	category, err := h.categoryService.Create(c.Request().Context(), tenantID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	response := h.categoryToResponse(category)
	return c.JSON(http.StatusCreated, response)
}

func (h *ProductHandler) GetCategories(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	// Parse query parameters
	query := domain.ProductCategoryQuery{
		Page:  1,
		Limit: 50,
	}

	if page := c.QueryParam("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			query.Page = p
		}
	}

	if limit := c.QueryParam("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			query.Limit = l
		}
	}

	if parentID := c.QueryParam("parent_id"); parentID != "" {
		if pid, err := strconv.ParseUint(parentID, 10, 64); err == nil {
			query.ParentID = &pid
		}
	}

	if isActive := c.QueryParam("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			query.IsActive = &active
		}
	}

	query.Sort = c.QueryParam("sort")
	query.Order = c.QueryParam("order")

	categories, total, err := h.categoryService.GetAll(c.Request().Context(), tenantID, query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get categories",
		})
	}

	categoryResponses := make([]domain.ProductCategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = h.categoryToResponse(category)
	}

	response := domain.ProductCategoryListResponse{
		Categories: categoryResponses,
		Total:      total,
		Page:       query.Page,
		Limit:      query.Limit,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetCategory(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid category ID",
		})
	}

	category, err := h.categoryService.GetByID(c.Request().Context(), tenantID, categoryID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Category not found",
		})
	}

	response := h.categoryToResponse(category)

	// Include subcategories if they exist
	if len(category.SubCategories) > 0 {
		subCategoryResponses := make([]domain.ProductCategoryResponse, len(category.SubCategories))
		for i, subCategory := range category.SubCategories {
			subCategoryResponses[i] = h.categoryToResponse(subCategory)
		}
		response.SubCategories = subCategoryResponses
	}

	return c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) UpdateCategory(c echo.Context) error {
	var req domain.UpdateProductCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	tenantID := c.Get("tenant_id").(uint64)

	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid category ID",
		})
	}

	category, err := h.categoryService.Update(c.Request().Context(), tenantID, categoryID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	response := h.categoryToResponse(category)
	return c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) DeleteCategory(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid category ID",
		})
	}

	err = h.categoryService.Delete(c.Request().Context(), tenantID, categoryID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Category deleted successfully",
	})
}

func (h *ProductHandler) GetCategoryHierarchy(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	categories, err := h.categoryService.GetHierarchy(c.Request().Context(), tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get category hierarchy",
		})
	}

	categoryResponses := make([]domain.ProductCategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = h.categoryToResponseWithHierarchy(category)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"categories": categoryResponses,
	})
}

// Helper functions

func (h *ProductHandler) categoryToResponse(category *domain.ProductCategory) domain.ProductCategoryResponse {
	response := domain.ProductCategoryResponse{
		ID:          category.ID,
		TenantID:    category.TenantID,
		ParentID:    category.ParentID,
		Name:        category.Name,
		Description: category.Description,
		ImageURL:    category.ImageURL,
		SortOrder:   category.SortOrder,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   category.UpdatedAt.Format(time.RFC3339),
	}

	if category.Parent != nil {
		parentResponse := h.categoryToResponse(category.Parent)
		response.Parent = &parentResponse
	}

	return response
}

func (h *ProductHandler) categoryToResponseWithHierarchy(category *domain.ProductCategory) domain.ProductCategoryResponse {
	response := h.categoryToResponse(category)

	if len(category.SubCategories) > 0 {
		subCategoryResponses := make([]domain.ProductCategoryResponse, len(category.SubCategories))
		for i, subCategory := range category.SubCategories {
			subCategoryResponses[i] = h.categoryToResponseWithHierarchy(subCategory)
		}
		response.SubCategories = subCategoryResponses
	}

	return response
}
