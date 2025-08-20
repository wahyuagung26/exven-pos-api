package handlers

import (
	"strconv"
	"time"

	"github.com/exven/pos-system/modules/subscription_plans/domain"
	"github.com/exven/pos-system/shared/utils/response"
	"github.com/labstack/echo/v4"
)

type SubscriptionPlanHandler struct {
	service domain.SubscriptionPlanService
}

func NewSubscriptionPlanHandler(service domain.SubscriptionPlanService) *SubscriptionPlanHandler {
	return &SubscriptionPlanHandler{
		service: service,
	}
}

func (h *SubscriptionPlanHandler) RegisterRoutes(e *echo.Group) {
	plans := e.Group("/subscription-plans")

	plans.GET("", h.GetSubscriptionPlans)
	plans.GET("/:id", h.GetSubscriptionPlan)
}

func (h *SubscriptionPlanHandler) GetSubscriptionPlans(c echo.Context) error {
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

	plans, total, err := h.service.GetAll(c.Request().Context(), limit, offset)
	if err != nil {
		return response.InternalError(c, "Failed to get subscription plans")
	}

	planResponses := make([]domain.SubscriptionPlanResponse, len(plans))
	for i, plan := range plans {
		planResponses[i] = h.planToResponse(plan)
	}

	return response.SuccessWithPagination(c, "Subscription plans retrieved successfully", planResponses, page, limit, int(total))
}

func (h *SubscriptionPlanHandler) GetSubscriptionPlan(c echo.Context) error {
	planID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return response.BadRequest(c, "Invalid subscription plan ID")
	}

	plan, err := h.service.GetByID(c.Request().Context(), planID)
	if err != nil {
		return response.NotFound(c, "Subscription plan not found")
	}

	planResponse := h.planToResponse(plan)
	return response.Success(c, "Subscription plan retrieved successfully", planResponse)
}


func (h *SubscriptionPlanHandler) planToResponse(plan *domain.SubscriptionPlan) domain.SubscriptionPlanResponse {
	return domain.SubscriptionPlanResponse{
		ID:                      plan.ID,
		Name:                    plan.Name,
		Description:             plan.Description,
		Price:                   plan.Price,
		MaxOutlets:              plan.MaxOutlets,
		MaxUsers:                plan.MaxUsers,
		MaxProducts:             plan.MaxProducts,
		MaxTransactionsPerMonth: plan.MaxTransactionsPerMonth,
		Features:                plan.Features,
		IsActive:                plan.IsActive,
		CreatedAt:               plan.CreatedAt.Format(time.RFC3339),
		UpdatedAt:               plan.UpdatedAt.Format(time.RFC3339),
	}
}