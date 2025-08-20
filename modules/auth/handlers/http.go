package handlers

import (
	"github.com/exven/pos-system/modules/auth/domain"
	"github.com/exven/pos-system/shared/utils/response"
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
		return response.BadRequest(c, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationErrorFromErr(c, err)
	}

	credentials := domain.LoginCredentials{
		Email:    req.Email,
		Password: req.Password,
	}

	tokenPair, user, err := h.authService.Login(c.Request().Context(), credentials)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	loginResponse := domain.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: domain.UserResponse{
			ID:       user.ID,
			TenantID: user.TenantID,
			Email:    user.Email,
			FullName: user.FullName,
			Phone:    user.Phone,
			IsActive: user.IsActive,
		},
	}

	if user.Role != nil {
		loginResponse.User.Role = domain.RoleResponse{
			ID:          user.Role.ID,
			Name:        user.Role.Name,
			DisplayName: user.Role.DisplayName,
			Permissions: user.Role.Permissions,
		}
	}

	return response.Success(c, "Login successful", loginResponse)
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req domain.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationErrorFromErr(c, err)
	}

	user, err := h.authService.Register(c.Request().Context(), req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	registerResponse := domain.RegisterResponse{
		TenantID: user.TenantID,
		User: domain.UserResponse{
			ID:       user.ID,
			TenantID: user.TenantID,
			Email:    user.Email,
			FullName: user.FullName,
			Phone:    user.Phone,
			IsActive: user.IsActive,
		},
		Message: "Registration successful. Please verify your email.",
	}

	return response.Created(c, "Registration successful. Please verify your email.", registerResponse)
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req domain.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationErrorFromErr(c, err)
	}

	tokenPair, err := h.authService.RefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	refreshResponse := map[string]interface{}{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"expires_in":    tokenPair.ExpiresIn,
	}

	return response.Success(c, "Token refreshed successfully", refreshResponse)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	userID := c.Get("user_id").(uint64)

	if err := h.authService.Logout(c.Request().Context(), userID); err != nil {
		return response.InternalError(c, "Failed to logout")
	}

	return response.Success(c, "Logout successful", nil)
}

func (h *AuthHandler) ChangePassword(c echo.Context) error {
	var req domain.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationErrorFromErr(c, err)
	}

	userID := c.Get("user_id").(uint64)

	if err := h.authService.ChangePassword(c.Request().Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Password changed successfully", nil)
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var req domain.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationErrorFromErr(c, err)
	}

	if err := h.authService.ResetPassword(c.Request().Context(), req.Email); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Password reset instructions sent to your email", nil)
}

func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	var req domain.VerifyEmailRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return response.ValidationErrorFromErr(c, err)
	}

	if err := h.authService.VerifyEmail(c.Request().Context(), req.Token); err != nil {
		return response.BadRequest(c, err.Error())
	}

	return response.Success(c, "Email verified successfully", nil)
}
