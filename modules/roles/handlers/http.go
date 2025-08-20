package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/exven/pos-system/modules/roles/domain"
	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	service domain.RoleService
}

func NewRoleHandler(service domain.RoleService) *RoleHandler {
	return &RoleHandler{
		service: service,
	}
}

func (h *RoleHandler) RegisterRoutes(e *echo.Group) {
	roles := e.Group("/roles")

	roles.GET("", h.GetRoles)
	roles.GET("/:id", h.GetRole)
	roles.GET("/name/:name", h.GetRoleByName)
	roles.GET("/system", h.GetSystemRoles)
}

func (h *RoleHandler) GetRoles(c echo.Context) error {
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

	roles, total, err := h.service.GetAll(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get roles",
		})
	}

	roleResponses := make([]domain.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = h.roleToResponse(role)
	}

	response := domain.RoleListResponse{
		Roles: roleResponses,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *RoleHandler) GetRole(c echo.Context) error {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid role ID",
		})
	}

	role, err := h.service.GetByID(c.Request().Context(), roleID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Role not found",
		})
	}

	response := h.roleToResponse(role)
	return c.JSON(http.StatusOK, response)
}

func (h *RoleHandler) GetRoleByName(c echo.Context) error {
	roleName := c.Param("name")

	role, err := h.service.GetByName(c.Request().Context(), roleName)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Role not found",
		})
	}

	response := h.roleToResponse(role)
	return c.JSON(http.StatusOK, response)
}

func (h *RoleHandler) GetSystemRoles(c echo.Context) error {
	roles, err := h.service.GetSystemRoles(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get system roles",
		})
	}

	roleResponses := make([]domain.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = h.roleToResponse(role)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"roles": roleResponses,
		"total": len(roleResponses),
	})
}

func (h *RoleHandler) roleToResponse(role *domain.Role) domain.RoleResponse {
	return domain.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: role.Description,
		Permissions: role.Permissions,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt.Format(time.RFC3339),
	}
}