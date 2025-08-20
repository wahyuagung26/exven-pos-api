package types

// SuccessResponse represents the standard success response format
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    *Meta       `json:"meta"`
}

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Message string                 `json:"message"`
	Data    interface{}            `json:"data"`
	Errors  map[string][]string    `json:"errors"`
}

// Meta contains pagination and metadata information
type Meta struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"per_page,omitempty"`
	Total   int `json:"total,omitempty"`
}