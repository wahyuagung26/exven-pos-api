package domain

type SubscriptionPlanResponse struct {
	ID                      uint64   `json:"id"`
	Name                    string   `json:"name"`
	Description             string   `json:"description"`
	Price                   float64  `json:"price"`
	MaxOutlets              int      `json:"max_outlets"`
	MaxUsers                int      `json:"max_users"`
	MaxProducts             *int     `json:"max_products"`
	MaxTransactionsPerMonth *int     `json:"max_transactions_per_month"`
	Features                []string `json:"features"`
	IsActive                bool     `json:"is_active"`
	CreatedAt               string   `json:"created_at"`
	UpdatedAt               string   `json:"updated_at"`
}

type SubscriptionPlanListResponse struct {
	Plans []SubscriptionPlanResponse `json:"plans"`
	Total int64                      `json:"total"`
	Page  int                        `json:"page"`
	Limit int                        `json:"limit"`
}