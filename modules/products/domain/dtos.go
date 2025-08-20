package domain

// ProductCategory DTOs

type CreateProductCategoryRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url" validate:"omitempty,url"`
	ParentID    *uint64 `json:"parent_id"`
	SortOrder   int     `json:"sort_order"`
}

type UpdateProductCategoryRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url" validate:"omitempty,url"`
	ParentID    *uint64 `json:"parent_id"`
	SortOrder   int     `json:"sort_order"`
	IsActive    bool    `json:"is_active"`
}

type ProductCategoryResponse struct {
	ID            uint64                    `json:"id"`
	TenantID      uint64                    `json:"tenant_id"`
	ParentID      *uint64                   `json:"parent_id"`
	Name          string                    `json:"name"`
	Description   string                    `json:"description"`
	ImageURL      string                    `json:"image_url"`
	SortOrder     int                       `json:"sort_order"`
	IsActive      bool                      `json:"is_active"`
	CreatedAt     string                    `json:"created_at"`
	UpdatedAt     string                    `json:"updated_at"`
	Parent        *ProductCategoryResponse  `json:"parent,omitempty"`
	SubCategories []ProductCategoryResponse `json:"sub_categories,omitempty"`
}

type ProductCategoryListResponse struct {
	Categories []ProductCategoryResponse `json:"categories"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
}

type ProductCategoryQuery struct {
	ParentID *uint64 `query:"parent_id"`
	IsActive *bool   `query:"is_active"`
	Page     int     `query:"page"`
	Limit    int     `query:"limit"`
	Sort     string  `query:"sort"`
	Order    string  `query:"order"`
}

// Product DTOs

type CreateProductRequest struct {
	CategoryID   *uint64                `json:"category_id"`
	SKU          string                 `json:"sku" validate:"required,max=100"`
	Barcode      string                 `json:"barcode" validate:"max=100"`
	Name         string                 `json:"name" validate:"required,min=1,max=255"`
	Description  string                 `json:"description"`
	Unit         string                 `json:"unit" validate:"max=50"`
	CostPrice    float64                `json:"cost_price" validate:"min=0"`
	SellingPrice float64                `json:"selling_price" validate:"required,min=0"`
	MinStock     int                    `json:"min_stock" validate:"min=0"`
	TrackStock   bool                   `json:"track_stock"`
	Images       []string               `json:"images"`
	Variants     map[string]interface{} `json:"variants"`
}

type UpdateProductRequest struct {
	CategoryID   *uint64                `json:"category_id"`
	SKU          string                 `json:"sku" validate:"required,max=100"`
	Barcode      string                 `json:"barcode" validate:"max=100"`
	Name         string                 `json:"name" validate:"required,min=1,max=255"`
	Description  string                 `json:"description"`
	Unit         string                 `json:"unit" validate:"max=50"`
	CostPrice    float64                `json:"cost_price" validate:"min=0"`
	SellingPrice float64                `json:"selling_price" validate:"required,min=0"`
	MinStock     int                    `json:"min_stock" validate:"min=0"`
	TrackStock   bool                   `json:"track_stock"`
	IsActive     bool                   `json:"is_active"`
	Images       []string               `json:"images"`
	Variants     map[string]interface{} `json:"variants"`
}

type ProductResponse struct {
	ID           uint64                   `json:"id"`
	TenantID     uint64                   `json:"tenant_id"`
	CategoryID   *uint64                  `json:"category_id"`
	SKU          string                   `json:"sku"`
	Barcode      string                   `json:"barcode"`
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	Unit         string                   `json:"unit"`
	CostPrice    float64                  `json:"cost_price"`
	SellingPrice float64                  `json:"selling_price"`
	MinStock     int                      `json:"min_stock"`
	TrackStock   bool                     `json:"track_stock"`
	IsActive     bool                     `json:"is_active"`
	Images       []string                 `json:"images"`
	Variants     map[string]interface{}   `json:"variants"`
	CreatedAt    string                   `json:"created_at"`
	UpdatedAt    string                   `json:"updated_at"`
	Category     *ProductCategoryResponse `json:"category,omitempty"`
}

type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

type ProductQuery struct {
	CategoryID *uint64 `query:"category_id"`
	SKU        string  `query:"sku"`
	Barcode    string  `query:"barcode"`
	Name       string  `query:"name"`
	IsActive   *bool   `query:"is_active"`
	Page       int     `query:"page"`
	Limit      int     `query:"limit"`
	Sort       string  `query:"sort"`
	Order      string  `query:"order"`
}
