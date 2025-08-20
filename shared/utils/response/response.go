package response

import (
	"net/http"

	"github.com/exven/pos-system/shared/types"
	"github.com/labstack/echo/v4"
)

// Success response for single item
func Success(c echo.Context, message string, data interface{}) error {
	response := types.SuccessResponse{
		Message: message,
		Data:    data,
		Meta:    nil,
	}
	return c.JSON(http.StatusOK, response)
}

// Created response for new resources
func Created(c echo.Context, message string, data interface{}) error {
	response := types.SuccessResponse{
		Message: message,
		Data:    data,
		Meta:    nil,
	}
	return c.JSON(http.StatusCreated, response)
}

// SuccessWithPagination response for list with pagination
func SuccessWithPagination(c echo.Context, message string, data interface{}, page, perPage, total int) error {
	response := types.SuccessResponse{
		Message: message,
		Data:    data,
		Meta: &types.Meta{
			Page:    page,
			PerPage: perPage,
			Total:   total,
		},
	}
	return c.JSON(http.StatusOK, response)
}

// ValidationError response with validation errors
func ValidationError(c echo.Context, fieldErrors map[string][]string) error {
	response := types.ErrorResponse{
		Message: "Validation failed",
		Data:    nil,
		Errors:  fieldErrors,
	}
	return c.JSON(http.StatusBadRequest, response)
}

// BadRequest error response
func BadRequest(c echo.Context, message string) error {
	response := types.ErrorResponse{
		Message: message,
		Data:    nil,
		Errors:  make(map[string][]string),
	}
	return c.JSON(http.StatusBadRequest, response)
}

// Unauthorized error response
func Unauthorized(c echo.Context, message string) error {
	response := types.ErrorResponse{
		Message: message,
		Data:    nil,
		Errors:  make(map[string][]string),
	}
	return c.JSON(http.StatusUnauthorized, response)
}

// NotFound error response
func NotFound(c echo.Context, message string) error {
	response := types.ErrorResponse{
		Message: message,
		Data:    nil,
		Errors:  make(map[string][]string),
	}
	return c.JSON(http.StatusNotFound, response)
}

// InternalError response
func InternalError(c echo.Context, message string) error {
	response := types.ErrorResponse{
		Message: message,
		Data:    nil,
		Errors:  make(map[string][]string),
	}
	return c.JSON(http.StatusInternalServerError, response)
}

// Error response with custom status code, message and errors
func Error(c echo.Context, statusCode int, message string, errors map[string][]string) error {
	if errors == nil {
		errors = make(map[string][]string)
	}

	response := types.ErrorResponse{
		Message: message,
		Data:    nil,
		Errors:  errors,
	}
	return c.JSON(statusCode, response)
}