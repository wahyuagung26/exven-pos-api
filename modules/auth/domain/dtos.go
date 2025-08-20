package domain

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	TenantID uint64 `json:"tenant_id" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	User         UserResponse `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RegisterRequest struct {
	TenantName   string `json:"tenant_name" validate:"required,min=3,max=255"`
	BusinessType string `json:"business_type"`
	Email        string `json:"email" validate:"required,email"`
	Phone        string `json:"phone"`
	Username     string `json:"username" validate:"required,min=3,max=100"`
	Password     string `json:"password" validate:"required,min=8"`
	FullName     string `json:"full_name" validate:"required,max=255"`
}

type RegisterResponse struct {
	TenantID uint64       `json:"tenant_id"`
	User     UserResponse `json:"user"`
	Message  string       `json:"message"`
}

type UserResponse struct {
	ID       uint64       `json:"id"`
	TenantID uint64       `json:"tenant_id"`
	Username string       `json:"username"`
	Email    string       `json:"email"`
	FullName string       `json:"full_name"`
	Phone    string       `json:"phone,omitempty"`
	Role     RoleResponse `json:"role"`
	IsActive bool         `json:"is_active"`
}

type RoleResponse struct {
	ID          uint64   `json:"id"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Permissions []string `json:"permissions"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email" validate:"required,email"`
	TenantID uint64 `json:"tenant_id" validate:"required"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}