package domain

type CreateOutletRequest struct {
	Name        string                 `json:"name" validate:"required,min=1,max=255"`
	Code        string                 `json:"code" validate:"required,min=1,max=50"`
	Description string                 `json:"description"`
	Address     string                 `json:"address"`
	City        string                 `json:"city" validate:"max=100"`
	Province    string                 `json:"province" validate:"max=100"`
	PostalCode  string                 `json:"postal_code" validate:"max=10"`
	Phone       string                 `json:"phone" validate:"max=20"`
	Email       string                 `json:"email" validate:"omitempty,email,max=255"`
	ManagerID   *uint64                `json:"manager_id"`
	Settings    map[string]interface{} `json:"settings"`
}

type UpdateOutletRequest struct {
	Name        string                 `json:"name" validate:"required,min=1,max=255"`
	Code        string                 `json:"code" validate:"required,min=1,max=50"`
	Description string                 `json:"description"`
	Address     string                 `json:"address"`
	City        string                 `json:"city" validate:"max=100"`
	Province    string                 `json:"province" validate:"max=100"`
	PostalCode  string                 `json:"postal_code" validate:"max=10"`
	Phone       string                 `json:"phone" validate:"max=20"`
	Email       string                 `json:"email" validate:"omitempty,email,max=255"`
	ManagerID   *uint64                `json:"manager_id"`
	IsActive    bool                   `json:"is_active"`
	Settings    map[string]interface{} `json:"settings"`
}

type OutletResponse struct {
	ID          uint64                 `json:"id"`
	TenantID    uint64                 `json:"tenant_id"`
	Name        string                 `json:"name"`
	Code        string                 `json:"code"`
	Description string                 `json:"description"`
	Address     string                 `json:"address"`
	City        string                 `json:"city"`
	Province    string                 `json:"province"`
	PostalCode  string                 `json:"postal_code"`
	Phone       string                 `json:"phone"`
	Email       string                 `json:"email"`
	ManagerID   *uint64                `json:"manager_id"`
	IsActive    bool                   `json:"is_active"`
	Settings    map[string]interface{} `json:"settings"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
	Manager     *ManagerResponse       `json:"manager,omitempty"`
}

type ManagerResponse struct {
	ID       uint64 `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type OutletListResponse struct {
	Outlets []OutletResponse `json:"outlets"`
	Total   int64            `json:"total"`
	Page    int              `json:"page"`
	Limit   int              `json:"limit"`
}

type OutletQuery struct {
	Name       string  `query:"name"`
	Code       string  `query:"code"`
	City       string  `query:"city"`
	Province   string  `query:"province"`
	ManagerID  *uint64 `query:"manager_id"`
	IsActive   *bool   `query:"is_active"`
	Page       int     `query:"page"`
	Limit      int     `query:"limit"`
	Sort       string  `query:"sort"`
	Order      string  `query:"order"`
}