package handlers

import (
	"net/http"

	"github.com/exven/pos-system/modules/auth/domain"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) RegisterRoutes(e *echo.Group) {
	auth := e.Group("/auth")
	auth.POST("/login", h.Login)
	auth.POST("/register", h.Register)
	auth.POST("/refresh", h.RefreshToken)
	auth.POST("/logout", h.Logout)
	auth.POST("/change-password", h.ChangePassword)
	auth.POST("/reset-password", h.ResetPassword)
	auth.POST("/verify-email", h.VerifyEmail)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req domain.LoginRequest
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

	credentials := domain.LoginCredentials{
		Username: req.Username,
		Password: req.Password,
		TenantID: req.TenantID,
	}

	tokenPair, user, err := h.authService.Login(c.Request().Context(), credentials)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	response := domain.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: domain.UserResponse{
			ID:       user.ID,
			TenantID: user.TenantID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Phone:    user.Phone,
			IsActive: user.IsActive,
		},
	}

	if user.Role != nil {
		response.User.Role = domain.RoleResponse{
			ID:          user.Role.ID,
			Name:        user.Role.Name,
			DisplayName: user.Role.DisplayName,
			Permissions: user.Role.Permissions,
		}
	}

	return c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req domain.RegisterRequest
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

	user, err := h.authService.Register(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	response := domain.RegisterResponse{
		TenantID: user.TenantID,
		User: domain.UserResponse{
			ID:       user.ID,
			TenantID: user.TenantID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Phone:    user.Phone,
			IsActive: user.IsActive,
		},
		Message: "Registration successful. Please verify your email.",
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req domain.RefreshTokenRequest
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

	tokenPair, err := h.authService.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"expires_in":    tokenPair.ExpiresIn,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	if err := h.authService.Logout(c.Request().Context(), userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to logout",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logout successful",
	})
}

func (h *AuthHandler) ChangePassword(c echo.Context) error {
	var req domain.ChangePasswordRequest
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

	userID := c.Get("user_id").(uint64)

	if err := h.authService.ChangePassword(c.Request().Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Password changed successfully",
	})
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var req domain.ResetPasswordRequest
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

	if err := h.authService.ResetPassword(c.Request().Context(), req.Email, req.TenantID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Password reset instructions sent to your email",
	})
}

func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	var req domain.VerifyEmailRequest
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

	if err := h.authService.VerifyEmail(c.Request().Context(), req.Token); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Email verified successfully",
	})
}