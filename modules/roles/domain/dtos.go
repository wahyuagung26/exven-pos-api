package domain

type RoleResponse struct {
	ID          uint64   `json:"id"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	IsSystem    bool     `json:"is_system"`
	CreatedAt   string   `json:"created_at"`
}

type RoleListResponse struct {
	Roles []RoleResponse `json:"roles"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}