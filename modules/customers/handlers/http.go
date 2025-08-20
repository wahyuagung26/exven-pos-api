package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/exven/pos-system/modules/customers/domain"
	"github.com/exven/pos-system/shared/utils/response"
	"github.com/labstack/echo/v4"
)

type CustomerHandler struct {
	customerService domain.CustomerService
}

func NewCustomerHandler(customerService domain.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

func (h *CustomerHandler) RegisterRoutes(e *echo.Group) {
	customers := e.Group("/customers")

	// Customer routes
	customers.POST("", h.CreateCustomer)
	customers.GET("", h.GetCustomers)
	customers.GET("/:id", h.GetCustomer)
	customers.PUT("/:id", h.UpdateCustomer)
	customers.DELETE("/:id", h.DeleteCustomer)
	customers.GET("/code/:code", h.GetCustomerByCode)
	customers.GET("/phone/:phone", h.GetCustomerByPhone)
	customers.GET("/email/:email", h.GetCustomerByEmail)
}

func (h *CustomerHandler) CreateCustomer(c echo.Context) error {
	var req domain.CreateCustomerRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	tenantID := c.Get("tenant_id").(uint64)

	customer, err := h.customerService.Create(c.Request().Context(), tenantID, req)
	if err != nil {
		if err.Error() == "customer with this code already exists" ||
			err.Error() == "customer with this phone already exists" ||
			err.Error() == "customer with this email already exists" {
			errors := map[string][]string{}
			if err.Error() == "customer with this code already exists" {
				errors["code"] = []string{"Code already exists"}
			} else if err.Error() == "customer with this phone already exists" {
				errors["phone"] = []string{"Phone already exists"}
			} else if err.Error() == "customer with this email already exists" {
				errors["email"] = []string{"Email already exists"}
			}
			return response.Error(c, http.StatusConflict, err.Error(), errors)
		}
		return response.BadRequest(c, err.Error())
	}

	responseData := h.customerToResponse(customer)
	return response.Created(c, "Customer created successfully", responseData)
}

func (h *CustomerHandler) GetCustomers(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	// Parse query parameters
	query := domain.CustomerQuery{
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
	query.Email = c.QueryParam("email")
	query.Phone = c.QueryParam("phone")
	query.City = c.QueryParam("city")
	query.Province = c.QueryParam("province")
	query.Gender = c.QueryParam("gender")

	if isActive := c.QueryParam("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			query.IsActive = &active
		}
	}

	query.Sort = c.QueryParam("sort")
	query.Order = c.QueryParam("order")

	customers, total, err := h.customerService.GetAll(c.Request().Context(), tenantID, query)
	if err != nil {
		return response.InternalError(c, "Failed to get customers")
	}

	customerResponses := make([]domain.CustomerResponse, len(customers))
	for i, customer := range customers {
		customerResponses[i] = h.customerToResponse(customer)
	}

	return response.SuccessWithPagination(c, "Customers retrieved successfully", customerResponses, query.Page, query.Limit, int(total))
}

func (h *CustomerHandler) GetCustomer(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	customerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid customer ID")
	}

	customer, err := h.customerService.GetByID(c.Request().Context(), tenantID, customerID)
	if err != nil {
		if err.Error() == "customer not found" {
			return response.NotFound(c, "Customer not found")
		}
		return response.InternalError(c, "Failed to get customer")
	}

	responseData := h.customerToResponse(customer)
	return response.Success(c, "Customer retrieved successfully", responseData)
}

func (h *CustomerHandler) UpdateCustomer(c echo.Context) error {
	var req domain.UpdateCustomerRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	tenantID := c.Get("tenant_id").(uint64)

	customerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid customer ID")
	}

	customer, err := h.customerService.Update(c.Request().Context(), tenantID, customerID, req)
	if err != nil {
		if err.Error() == "customer not found" {
			return response.NotFound(c, "Customer not found")
		}
		if err.Error() == "customer with this code already exists" ||
			err.Error() == "customer with this phone already exists" ||
			err.Error() == "customer with this email already exists" {
			errors := map[string][]string{}
			if err.Error() == "customer with this code already exists" {
				errors["code"] = []string{"Code already exists"}
			} else if err.Error() == "customer with this phone already exists" {
				errors["phone"] = []string{"Phone already exists"}
			} else if err.Error() == "customer with this email already exists" {
				errors["email"] = []string{"Email already exists"}
			}
			return response.Error(c, http.StatusConflict, err.Error(), errors)
		}
		return response.BadRequest(c, err.Error())
	}

	responseData := h.customerToResponse(customer)
	return response.Success(c, "Customer updated successfully", responseData)
}

func (h *CustomerHandler) DeleteCustomer(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)

	customerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid customer ID")
	}

	err = h.customerService.Delete(c.Request().Context(), tenantID, customerID)
	if err != nil {
		if err.Error() == "customer not found" {
			return response.NotFound(c, "Customer not found")
		}
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Customer deleted successfully", nil)
}

func (h *CustomerHandler) GetCustomerByCode(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)
	code := c.Param("code")

	customer, err := h.customerService.GetByCode(c.Request().Context(), tenantID, code)
	if err != nil {
		if err.Error() == "customer not found" {
			return response.NotFound(c, "Customer not found")
		}
		if err.Error() == "code cannot be empty" {
			return response.BadRequest(c, "Code cannot be empty")
		}
		return response.InternalError(c, "Failed to get customer")
	}

	responseData := h.customerToResponse(customer)
	return response.Success(c, "Customer retrieved successfully", responseData)
}

func (h *CustomerHandler) GetCustomerByPhone(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)
	phone := c.Param("phone")

	customer, err := h.customerService.GetByPhone(c.Request().Context(), tenantID, phone)
	if err != nil {
		if err.Error() == "customer not found" {
			return response.NotFound(c, "Customer not found")
		}
		if err.Error() == "phone cannot be empty" {
			return response.BadRequest(c, "Phone cannot be empty")
		}
		return response.InternalError(c, "Failed to get customer")
	}

	responseData := h.customerToResponse(customer)
	return response.Success(c, "Customer retrieved successfully", responseData)
}

func (h *CustomerHandler) GetCustomerByEmail(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint64)
	email := c.Param("email")

	customer, err := h.customerService.GetByEmail(c.Request().Context(), tenantID, email)
	if err != nil {
		if err.Error() == "customer not found" {
			return response.NotFound(c, "Customer not found")
		}
		if err.Error() == "email cannot be empty" {
			return response.BadRequest(c, "Email cannot be empty")
		}
		return response.InternalError(c, "Failed to get customer")
	}

	responseData := h.customerToResponse(customer)
	return response.Success(c, "Customer retrieved successfully", responseData)
}

// Helper functions

func (h *CustomerHandler) customerToResponse(customer *domain.Customer) domain.CustomerResponse {
	return domain.CustomerResponse{
		ID:            customer.ID,
		TenantID:      customer.TenantID,
		Code:          customer.Code,
		Name:          customer.Name,
		Email:         customer.Email,
		Phone:         customer.Phone,
		Address:       customer.Address,
		City:          customer.City,
		Province:      customer.Province,
		PostalCode:    customer.PostalCode,
		BirthDate:     customer.BirthDate,
		Gender:        customer.Gender,
		LoyaltyPoints: customer.LoyaltyPoints,
		TotalSpent:    customer.TotalSpent,
		VisitCount:    customer.VisitCount,
		LastVisitAt:   customer.LastVisitAt,
		Notes:         customer.Notes,
		IsActive:      customer.IsActive,
		CreatedAt:     customer.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     customer.UpdatedAt.Format(time.RFC3339),
	}
}