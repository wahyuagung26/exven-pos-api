package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/exven/pos-system/modules/outlets/domain"
	"github.com/labstack/echo/v4"
)

type OutletHandler struct {
	outletService domain.OutletService
}

func NewOutletHandler(outletService domain.OutletService) *OutletHandler {
	return &OutletHandler{
		outletService: outletService,
	}
}

func (h *OutletHandler) RegisterRoutes(e *echo.Group) {
	outlets := e.Group("/outlets")

	// Outlet routes
	outlets.POST("", h.CreateOutlet)
	outlets.GET("", h.GetOutlets)
	outlets.GET("/:id", h.GetOutlet)
	outlets.PUT("/:id", h.UpdateOutlet)
	outlets.DELETE("/:id", h.DeleteOutlet)
	outlets.GET("/code/:code", h.GetOutletByCode)
}

func (h *OutletHandler) CreateOutlet(c echo.Context) error {
	var req domain.CreateOutletRequest
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

	outlet, err := h.outletService.Create(c.Request().Context(), tenantID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	response := h.outletToResponse(outlet)
	return c.JSON(http.StatusCreated, response)
}

func (h *OutletHandler) GetOutlets(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	// Parse query parameters
	query := domain.OutletQuery{
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

	query.Name = c.QueryParam("name")
	query.Code = c.QueryParam("code")
	query.City = c.QueryParam("city")
	query.Province = c.QueryParam("province")

	if managerID := c.QueryParam("manager_id"); managerID != "" {
		if mid, err := strconv.ParseUint(managerID, 10, 64); err == nil {
			query.ManagerID = &mid
		}
	}

	if isActive := c.QueryParam("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			query.IsActive = &active
		}
	}

	query.Sort = c.QueryParam("sort")
	query.Order = c.QueryParam("order")

	outlets, total, err := h.outletService.GetAll(c.Request().Context(), tenantID, query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get outlets",
		})
	}

	outletResponses := make([]domain.OutletResponse, len(outlets))
	for i, outlet := range outlets {
		outletResponses[i] = h.outletToResponse(outlet)
	}

	response := domain.OutletListResponse{
		Outlets: outletResponses,
		Total:   total,
		Page:    query.Page,
		Limit:   query.Limit,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *OutletHandler) GetOutlet(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	outletID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid outlet ID",
		})
	}

	outlet, err := h.outletService.GetByID(c.Request().Context(), tenantID, outletID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Outlet not found",
		})
	}

	response := h.outletToResponse(outlet)
	return c.JSON(http.StatusOK, response)
}

func (h *OutletHandler) UpdateOutlet(c echo.Context) error {
	var req domain.UpdateOutletRequest
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

	outletID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid outlet ID",
		})
	}

	outlet, err := h.outletService.Update(c.Request().Context(), tenantID, outletID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	response := h.outletToResponse(outlet)
	return c.JSON(http.StatusOK, response)
}

func (h *OutletHandler) DeleteOutlet(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	outletID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid outlet ID",
		})
	}

	err = h.outletService.Delete(c.Request().Context(), tenantID, outletID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Outlet deleted successfully",
	})
}

func (h *OutletHandler) GetOutletByCode(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)
	code := c.Param("code")

	outlet, err := h.outletService.GetByCode(c.Request().Context(), tenantID, code)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Outlet not found",
		})
	}

	response := h.outletToResponse(outlet)
	return c.JSON(http.StatusOK, response)
}

// Helper functions

func (h *OutletHandler) outletToResponse(outlet *domain.Outlet) domain.OutletResponse {
	response := domain.OutletResponse{
		ID:          outlet.ID,
		TenantID:    outlet.TenantID,
		Name:        outlet.Name,
		Code:        outlet.Code,
		Description: outlet.Description,
		Address:     outlet.Address,
		City:        outlet.City,
		Province:    outlet.Province,
		PostalCode:  outlet.PostalCode,
		Phone:       outlet.Phone,
		Email:       outlet.Email,
		ManagerID:   outlet.ManagerID,
		IsActive:    outlet.IsActive,
		Settings:    outlet.Settings,
		CreatedAt:   outlet.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   outlet.UpdatedAt.Format(time.RFC3339),
	}

	if outlet.Manager != nil {
		response.Manager = &domain.ManagerResponse{
			ID:       outlet.Manager.ID,
			FullName: outlet.Manager.FullName,
			Email:    outlet.Manager.Email,
			Phone:    outlet.Manager.Phone,
		}
	}

	return response
}