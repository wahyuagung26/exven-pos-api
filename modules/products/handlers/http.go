package handlers

import (
	"strconv"
	"time"

	"github.com/exven/pos-system/modules/products/domain"
	"github.com/exven/pos-system/shared/utils/response"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	categoryService domain.ProductCategoryService
	productService  domain.ProductService
}

func NewProductHandler(categoryService domain.ProductCategoryService, productService domain.ProductService) *ProductHandler {
	return &ProductHandler{
		categoryService: categoryService,
		productService:  productService,
	}
}

func (h *ProductHandler) RegisterRoutes(e *echo.Group) {
	products := e.Group("/products")

	// Product routes
	products.POST("", h.CreateProduct)
	products.GET("", h.GetProducts)
	products.GET("/:id", h.GetProduct)
	products.PUT("/:id", h.UpdateProduct)
	products.DELETE("/:id", h.DeleteProduct)
	products.GET("/sku/:sku", h.GetProductBySKU)
	products.GET("/barcode/:barcode", h.GetProductByBarcode)

	// Product Categories routes
	categories := products.Group("/categories")
	categories.POST("", h.CreateCategory)
	categories.GET("", h.GetCategories)
	categories.GET("/:id", h.GetCategory)
	categories.PUT("/:id", h.UpdateCategory)
	categories.DELETE("/:id", h.DeleteCategory)
	categories.GET("/hierarchy", h.GetCategoryHierarchy)
	categories.GET("/:id/products", h.GetProductsByCategory)
}

// Product handlers

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var req domain.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationError(c, map[string][]string{
			"request": {err.Error()},
		})
	}

	tenantID := c.Get("tenant_id").(uint64)

	product, err := h.productService.Create(c.Request().Context(), tenantID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	productResponse := h.productToResponse(product)
	return response.Created(c, "Product created successfully", productResponse)
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	// Parse pagination parameters
	page := 1
	limit := 50

	if p := c.QueryParam("page"); p != "" {
		if pageInt, err := strconv.Atoi(p); err == nil && pageInt > 0 {
			page = pageInt
		}
	}

	if l := c.QueryParam("limit"); l != "" {
		if limitInt, err := strconv.Atoi(l); err == nil && limitInt > 0 && limitInt <= 100 {
			limit = limitInt
		}
	}

	offset := (page - 1) * limit

	products, total, err := h.productService.GetAll(c.Request().Context(), tenantID, limit, offset)
	if err != nil {
		return response.InternalError(c, "Failed to get products")
	}

	productResponses := make([]domain.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = h.productToResponse(product)
	}

	return response.SuccessWithPagination(c, "Products retrieved successfully", productResponses, page, limit, int(total))
}

func (h *ProductHandler) GetProduct(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID")
	}

	product, err := h.productService.GetByID(c.Request().Context(), tenantID, productID)
	if err != nil {
		return response.NotFound(c, "Product not found")
	}

	productResponse := h.productToResponse(product)
	return response.Success(c, "Product retrieved successfully", productResponse)
}

func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	var req domain.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationError(c, map[string][]string{
			"request": {err.Error()},
		})
	}

	tenantID := c.Get("tenant_id").(uint64)

	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID")
	}

	product, err := h.productService.Update(c.Request().Context(), tenantID, productID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	productResponse := h.productToResponse(product)
	return response.Success(c, "Product updated successfully", productResponse)
}

func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid product ID")
	}

	err = h.productService.Delete(c.Request().Context(), tenantID, productID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Product deleted successfully", nil)
}

func (h *ProductHandler) GetProductBySKU(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)
	sku := c.Param("sku")

	product, err := h.productService.GetBySKU(c.Request().Context(), tenantID, sku)
	if err != nil {
		return response.NotFound(c, "Product not found")
	}

	productResponse := h.productToResponse(product)
	return response.Success(c, "Product retrieved successfully", productResponse)
}

func (h *ProductHandler) GetProductByBarcode(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)
	barcode := c.Param("barcode")

	product, err := h.productService.GetByBarcode(c.Request().Context(), tenantID, barcode)
	if err != nil {
		return response.NotFound(c, "Product not found")
	}

	productResponse := h.productToResponse(product)
	return response.Success(c, "Product retrieved successfully", productResponse)
}

func (h *ProductHandler) GetProductsByCategory(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID")
	}

	// Parse pagination parameters
	page := 1
	limit := 50

	if p := c.QueryParam("page"); p != "" {
		if pageInt, err := strconv.Atoi(p); err == nil && pageInt > 0 {
			page = pageInt
		}
	}

	if l := c.QueryParam("limit"); l != "" {
		if limitInt, err := strconv.Atoi(l); err == nil && limitInt > 0 && limitInt <= 100 {
			limit = limitInt
		}
	}

	offset := (page - 1) * limit

	products, total, err := h.productService.GetByCategory(c.Request().Context(), tenantID, categoryID, limit, offset)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	productResponses := make([]domain.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = h.productToResponse(product)
	}

	return response.SuccessWithPagination(c, "Products retrieved successfully", productResponses, page, limit, int(total))
}

// Product Category handlers

func (h *ProductHandler) CreateCategory(c echo.Context) error {
	var req domain.CreateProductCategoryRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	tenantID := c.Get("tenant_id").(uint64)

	category, err := h.categoryService.Create(c.Request().Context(), tenantID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	categoryResponse := h.categoryToResponse(category)
	return response.Created(c, "Category created successfully", categoryResponse)
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
		return response.InternalError(c, "Failed to get categories")
	}

	categoryResponses := make([]domain.ProductCategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = h.categoryToResponse(category)
	}

	return response.SuccessWithPagination(c, "Categories retrieved successfully", categoryResponses, query.Page, query.Limit, int(total))
}

func (h *ProductHandler) GetCategory(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID")
	}

	category, err := h.categoryService.GetByID(c.Request().Context(), tenantID, categoryID)
	if err != nil {
		return response.NotFound(c, "Category not found")
	}

	categoryResponse := h.categoryToResponse(category)

	// Include subcategories if they exist
	if len(category.SubCategories) > 0 {
		subCategoryResponses := make([]domain.ProductCategoryResponse, len(category.SubCategories))
		for i, subCategory := range category.SubCategories {
			subCategoryResponses[i] = h.categoryToResponse(subCategory)
		}
		categoryResponse.SubCategories = subCategoryResponses
	}

	return response.Success(c, "Category retrieved successfully", categoryResponse)
}

func (h *ProductHandler) UpdateCategory(c echo.Context) error {
	var req domain.UpdateProductCategoryRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	tenantID := c.Get("tenant_id").(uint64)

	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID")
	}

	category, err := h.categoryService.Update(c.Request().Context(), tenantID, categoryID, req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	categoryResponse := h.categoryToResponse(category)
	return response.Success(c, "Category updated successfully", categoryResponse)
}

func (h *ProductHandler) DeleteCategory(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid category ID")
	}

	err = h.categoryService.Delete(c.Request().Context(), tenantID, categoryID)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Category deleted successfully", nil)
}

func (h *ProductHandler) GetCategoryHierarchy(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	categories, err := h.categoryService.GetHierarchy(c.Request().Context(), tenantID)
	if err != nil {
		return response.InternalError(c, "Failed to get category hierarchy")
	}

	categoryResponses := make([]domain.ProductCategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = h.categoryToResponseWithHierarchy(category)
	}

	return response.Success(c, "Category hierarchy retrieved successfully", categoryResponses)
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

func (h *ProductHandler) productToResponse(product *domain.Product) domain.ProductResponse {
	response := domain.ProductResponse{
		ID:           product.ID,
		TenantID:     product.TenantID,
		CategoryID:   product.CategoryID,
		SKU:          product.SKU,
		Barcode:      product.Barcode,
		Name:         product.Name,
		Description:  product.Description,
		Unit:         product.Unit,
		CostPrice:    product.CostPrice,
		SellingPrice: product.SellingPrice,
		MinStock:     product.MinStock,
		TrackStock:   product.TrackStock,
		IsActive:     product.IsActive,
		Images:       product.Images,
		Variants:     product.Variants,
		CreatedAt:    product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    product.UpdatedAt.Format(time.RFC3339),
	}

	if product.Category != nil {
		categoryResponse := h.categoryToResponse(product.Category)
		response.Category = &categoryResponse
	}

	return response
}
