package handlers

import (
	"strconv"
	"time"

	"github.com/exven/pos-system/modules/roles/domain"
	"github.com/exven/pos-system/shared/utils/response"
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
		return response.InternalError(c, "Failed to get roles")
	}

	roleResponses := make([]domain.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = h.roleToResponse(role)
	}

	return response.SuccessWithPagination(c, "Roles retrieved successfully", roleResponses, page, limit, int(total))
}

func (h *RoleHandler) GetRole(c echo.Context) error {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid role ID")
	}

	role, err := h.service.GetByID(c.Request().Context(), roleID)
	if err != nil {
		return response.NotFound(c, "Role not found")
	}

	roleResponse := h.roleToResponse(role)
	return response.Success(c, "Role retrieved successfully", roleResponse)
}

func (h *RoleHandler) GetRoleByName(c echo.Context) error {
	roleName := c.Param("name")

	role, err := h.service.GetByName(c.Request().Context(), roleName)
	if err != nil {
		return response.NotFound(c, "Role not found")
	}

	roleResponse := h.roleToResponse(role)
	return response.Success(c, "Role retrieved successfully", roleResponse)
}

func (h *RoleHandler) GetSystemRoles(c echo.Context) error {
	roles, err := h.service.GetSystemRoles(c.Request().Context())
	if err != nil {
		return response.InternalError(c, "Failed to get system roles")
	}

	roleResponses := make([]domain.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = h.roleToResponse(role)
	}

	return response.Success(c, "System roles retrieved successfully", roleResponses)
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