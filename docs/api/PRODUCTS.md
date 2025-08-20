# Products API Documentation

This document provides comprehensive API documentation for the Products module of ExVen POS Lite system.

## Overview

The Products API manages product catalog and categories within a multi-tenant POS system. It includes product management with SKU/barcode support, hierarchical categories, stock tracking, variants, and pricing. Each product belongs to a specific tenant and can be organized into categories for better management.

## Base URL

All products API endpoints are prefixed with `/api/v1/products`

## Authentication

All endpoints require JWT authentication. The JWT token must be included in the Authorization header:

```
Authorization: Bearer <jwt_token>
```

## Response Format

All API responses follow the standard response format:

```json
{
  "message": "Success message",
  "data": {},
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 150
  }
}
```

---

## Product Endpoints

### 1. Create Product

Creates a new product for the authenticated tenant.

**Endpoint:** `POST /api/v1/products`

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "category_id": 1,
  "sku": "PROD001",
  "barcode": "1234567890123",
  "name": "Premium Coffee Beans",
  "description": "High-quality arabica coffee beans",
  "unit": "kg",
  "cost_price": 50000.00,
  "selling_price": 75000.00,
  "min_stock": 10,
  "track_stock": true,
  "is_active": true,
  "images": [
    "https://example.com/product1.jpg",
    "https://example.com/product1-alt.jpg"
  ],
  "variants": {
    "size": ["250g", "500g", "1kg"],
    "roast_level": ["light", "medium", "dark"]
  }
}
```

**Validation Rules:**
- `category_id`: Optional, must be valid category ID within tenant
- `sku`: Required, 1-100 characters, unique within tenant
- `barcode`: Optional, 1-100 characters, unique within tenant if provided
- `name`: Required, 1-255 characters
- `description`: Optional
- `unit`: Optional, max 50 characters (default: "pcs")
- `cost_price`: Optional, decimal (default: 0.00)
- `selling_price`: Required, decimal, must be greater than 0
- `min_stock`: Optional, integer (default: 0)
- `track_stock`: Optional, boolean (default: true)
- `is_active`: Optional, boolean (default: true)
- `images`: Optional, array of image URLs
- `variants`: Optional, JSON object for product variants

**Response:**

*Success (201 Created):*
```json
{
  "message": "Product created successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "category_id": 1,
    "sku": "PROD001",
    "barcode": "1234567890123",
    "name": "Premium Coffee Beans",
    "description": "High-quality arabica coffee beans",
    "unit": "kg",
    "cost_price": 50000.00,
    "selling_price": 75000.00,
    "min_stock": 10,
    "track_stock": true,
    "is_active": true,
    "images": [
      "https://example.com/product1.jpg",
      "https://example.com/product1-alt.jpg"
    ],
    "variants": {
      "size": ["250g", "500g", "1kg"],
      "roast_level": ["light", "medium", "dark"]
    },
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T10:30:00Z",
    "category": {
      "id": 1,
      "name": "Beverages",
      "description": "Coffee and tea products"
    }
  },
  "meta": null
}
```

*Error (400 Bad Request):*
```json
{
  "message": "Validation failed",
  "data": null,
  "errors": {
    "sku": ["SKU already exists"],
    "selling_price": ["Selling price is required"]
  }
}
```

---

### 2. Get All Products

Retrieves a paginated list of products for the authenticated tenant with optional filtering and sorting.

**Endpoint:** `GET /api/v1/products`

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 50, max: 100)
- `name`: Filter by product name (partial match)
- `sku`: Filter by SKU (partial match)
- `category_id`: Filter by category ID
- `is_active`: Filter by active status (true/false)
- `track_stock`: Filter by stock tracking (true/false)
- `sort`: Sort field (name, sku, selling_price, created_at)
- `order`: Sort order (asc, desc, default: asc)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Example Request:**
```
GET /api/v1/products?page=1&limit=20&category_id=1&is_active=true&sort=name&order=asc
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Products retrieved successfully",
  "data": [
    {
      "id": 1,
      "tenant_id": 100,
      "category_id": 1,
      "sku": "PROD001",
      "barcode": "1234567890123",
      "name": "Premium Coffee Beans",
      "description": "High-quality arabica coffee beans",
      "unit": "kg",
      "cost_price": 50000.00,
      "selling_price": 75000.00,
      "min_stock": 10,
      "track_stock": true,
      "is_active": true,
      "images": [
        "https://example.com/product1.jpg"
      ],
      "variants": {
        "size": ["250g", "500g", "1kg"]
      },
      "created_at": "2025-08-20T10:30:00Z",
      "updated_at": "2025-08-20T10:30:00Z",
      "category": {
        "id": 1,
        "name": "Beverages",
        "description": "Coffee and tea products"
      }
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1
  }
}
```

---

### 3. Get Product by ID

Retrieves a specific product by its ID.

**Endpoint:** `GET /api/v1/products/{id}`

**Path Parameters:**
- `id`: Product ID (integer, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Product retrieved successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "category_id": 1,
    "sku": "PROD001",
    "barcode": "1234567890123",
    "name": "Premium Coffee Beans",
    "description": "High-quality arabica coffee beans",
    "unit": "kg",
    "cost_price": 50000.00,
    "selling_price": 75000.00,
    "min_stock": 10,
    "track_stock": true,
    "is_active": true,
    "images": [
      "https://example.com/product1.jpg",
      "https://example.com/product1-alt.jpg"
    ],
    "variants": {
      "size": ["250g", "500g", "1kg"],
      "roast_level": ["light", "medium", "dark"]
    },
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T10:30:00Z",
    "category": {
      "id": 1,
      "name": "Beverages",
      "description": "Coffee and tea products"
    }
  },
  "meta": null
}
```

*Error (404 Not Found):*
```json
{
  "message": "Product not found",
  "data": null,
  "errors": {}
}
```

---

### 4. Update Product

Updates an existing product.

**Endpoint:** `PUT /api/v1/products/{id}`

**Path Parameters:**
- `id`: Product ID (integer, required)

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "category_id": 1,
  "sku": "PROD001",
  "barcode": "1234567890123",
  "name": "Premium Coffee Beans Updated",
  "description": "Updated description for high-quality arabica coffee beans",
  "unit": "kg",
  "cost_price": 55000.00,
  "selling_price": 80000.00,
  "min_stock": 15,
  "track_stock": true,
  "is_active": true,
  "images": [
    "https://example.com/product1-new.jpg"
  ],
  "variants": {
    "size": ["250g", "500g", "1kg", "2kg"],
    "roast_level": ["light", "medium", "dark"]
  }
}
```

**Validation Rules:**
- Same as Create Product request
- All fields are optional in update request

**Response:**

*Success (200 OK):*
```json
{
  "message": "Product updated successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "category_id": 1,
    "sku": "PROD001",
    "barcode": "1234567890123",
    "name": "Premium Coffee Beans Updated",
    "description": "Updated description for high-quality arabica coffee beans",
    "unit": "kg",
    "cost_price": 55000.00,
    "selling_price": 80000.00,
    "min_stock": 15,
    "track_stock": true,
    "is_active": true,
    "images": [
      "https://example.com/product1-new.jpg"
    ],
    "variants": {
      "size": ["250g", "500g", "1kg", "2kg"],
      "roast_level": ["light", "medium", "dark"]
    },
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T11:30:00Z",
    "category": {
      "id": 1,
      "name": "Beverages",
      "description": "Coffee and tea products"
    }
  },
  "meta": null
}
```

---

### 5. Delete Product

Soft deletes a product by setting `is_active` to false.

**Endpoint:** `DELETE /api/v1/products/{id}`

**Path Parameters:**
- `id`: Product ID (integer, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Product deleted successfully",
  "data": null,
  "meta": null
}
```

*Error (400 Bad Request):*
```json
{
  "message": "Cannot delete product",
  "data": null,
  "errors": {
    "product": ["Cannot delete product with active stock or transactions"]
  }
}
```

---

### 6. Get Product by SKU

Retrieves a specific product by its SKU.

**Endpoint:** `GET /api/v1/products/sku/{sku}`

**Path Parameters:**
- `sku`: Product SKU (string, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Product retrieved successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "category_id": 1,
    "sku": "PROD001",
    "barcode": "1234567890123",
    "name": "Premium Coffee Beans",
    "description": "High-quality arabica coffee beans",
    "unit": "kg",
    "cost_price": 50000.00,
    "selling_price": 75000.00,
    "min_stock": 10,
    "track_stock": true,
    "is_active": true,
    "images": [
      "https://example.com/product1.jpg"
    ],
    "variants": {
      "size": ["250g", "500g", "1kg"]
    },
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T10:30:00Z",
    "category": {
      "id": 1,
      "name": "Beverages",
      "description": "Coffee and tea products"
    }
  },
  "meta": null
}
```

---

### 7. Get Product by Barcode

Retrieves a specific product by its barcode.

**Endpoint:** `GET /api/v1/products/barcode/{barcode}`

**Path Parameters:**
- `barcode`: Product barcode (string, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Product retrieved successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "category_id": 1,
    "sku": "PROD001",
    "barcode": "1234567890123",
    "name": "Premium Coffee Beans",
    "description": "High-quality arabica coffee beans",
    "unit": "kg",
    "cost_price": 50000.00,
    "selling_price": 75000.00,
    "min_stock": 10,
    "track_stock": true,
    "is_active": true,
    "images": [
      "https://example.com/product1.jpg"
    ],
    "variants": {
      "size": ["250g", "500g", "1kg"]
    },
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T10:30:00Z",
    "category": {
      "id": 1,
      "name": "Beverages",
      "description": "Coffee and tea products"
    }
  },
  "meta": null
}
```

---

## Product Category Endpoints

### 8. Create Product Category

Creates a new product category for the authenticated tenant.

**Endpoint:** `POST /api/v1/products/categories`

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "parent_id": null,
  "name": "Beverages",
  "description": "Coffee, tea, and other beverages",
  "image_url": "https://example.com/beverages.jpg",
  "sort_order": 1,
  "is_active": true
}
```

**Validation Rules:**
- `parent_id`: Optional, must be valid category ID within tenant for sub-categories
- `name`: Required, 1-255 characters
- `description`: Optional
- `image_url`: Optional, valid URL format, max 500 characters
- `sort_order`: Optional, integer (default: 0)
- `is_active`: Optional, boolean (default: true)

**Response:**

*Success (201 Created):*
```json
{
  "message": "Category created successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "parent_id": null,
    "name": "Beverages",
    "description": "Coffee, tea, and other beverages",
    "image_url": "https://example.com/beverages.jpg",
    "sort_order": 1,
    "is_active": true,
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T10:30:00Z",
    "parent": null,
    "sub_categories": []
  },
  "meta": null
}
```

---

### 9. Get All Product Categories

Retrieves a paginated list of product categories for the authenticated tenant with optional filtering and sorting.

**Endpoint:** `GET /api/v1/products/categories`

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 50, max: 100)
- `parent_id`: Filter by parent category ID (null for root categories)
- `is_active`: Filter by active status (true/false)
- `sort`: Sort field (name, sort_order, created_at)
- `order`: Sort order (asc, desc, default: asc)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Example Request:**
```
GET /api/v1/products/categories?page=1&limit=20&is_active=true&sort=sort_order&order=asc
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Categories retrieved successfully",
  "data": [
    {
      "id": 1,
      "tenant_id": 100,
      "parent_id": null,
      "name": "Beverages",
      "description": "Coffee, tea, and other beverages",
      "image_url": "https://example.com/beverages.jpg",
      "sort_order": 1,
      "is_active": true,
      "created_at": "2025-08-20T10:30:00Z",
      "updated_at": "2025-08-20T10:30:00Z",
      "parent": null
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1
  }
}
```

---

### 10. Get Product Category by ID

Retrieves a specific product category by its ID, including subcategories if they exist.

**Endpoint:** `GET /api/v1/products/categories/{id}`

**Path Parameters:**
- `id`: Category ID (integer, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Category retrieved successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "parent_id": null,
    "name": "Beverages",
    "description": "Coffee, tea, and other beverages",
    "image_url": "https://example.com/beverages.jpg",
    "sort_order": 1,
    "is_active": true,
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T10:30:00Z",
    "parent": null,
    "sub_categories": [
      {
        "id": 2,
        "tenant_id": 100,
        "parent_id": 1,
        "name": "Coffee",
        "description": "Various coffee products",
        "image_url": "https://example.com/coffee.jpg",
        "sort_order": 1,
        "is_active": true,
        "created_at": "2025-08-20T10:35:00Z",
        "updated_at": "2025-08-20T10:35:00Z"
      }
    ]
  },
  "meta": null
}
```

---

### 11. Update Product Category

Updates an existing product category.

**Endpoint:** `PUT /api/v1/products/categories/{id}`

**Path Parameters:**
- `id`: Category ID (integer, required)

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "parent_id": null,
  "name": "Beverages Updated",
  "description": "Updated description for beverages category",
  "image_url": "https://example.com/beverages-new.jpg",
  "sort_order": 2,
  "is_active": true
}
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Category updated successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "parent_id": null,
    "name": "Beverages Updated",
    "description": "Updated description for beverages category",
    "image_url": "https://example.com/beverages-new.jpg",
    "sort_order": 2,
    "is_active": true,
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T11:30:00Z",
    "parent": null
  },
  "meta": null
}
```

---

### 12. Delete Product Category

Soft deletes a product category by setting `is_active` to false.

**Endpoint:** `DELETE /api/v1/products/categories/{id}`

**Path Parameters:**
- `id`: Category ID (integer, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Category deleted successfully",
  "data": null,
  "meta": null
}
```

*Error (400 Bad Request):*
```json
{
  "message": "Cannot delete category",
  "data": null,
  "errors": {
    "category": ["Cannot delete category with active products or subcategories"]
  }
}
```

---

### 13. Get Category Hierarchy

Retrieves the complete category hierarchy for the authenticated tenant with nested subcategories.

**Endpoint:** `GET /api/v1/products/categories/hierarchy`

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Category hierarchy retrieved successfully",
  "data": [
    {
      "id": 1,
      "tenant_id": 100,
      "parent_id": null,
      "name": "Beverages",
      "description": "Coffee, tea, and other beverages",
      "image_url": "https://example.com/beverages.jpg",
      "sort_order": 1,
      "is_active": true,
      "created_at": "2025-08-20T10:30:00Z",
      "updated_at": "2025-08-20T10:30:00Z",
      "sub_categories": [
        {
          "id": 2,
          "tenant_id": 100,
          "parent_id": 1,
          "name": "Coffee",
          "description": "Various coffee products",
          "image_url": "https://example.com/coffee.jpg",
          "sort_order": 1,
          "is_active": true,
          "created_at": "2025-08-20T10:35:00Z",
          "updated_at": "2025-08-20T10:35:00Z",
          "sub_categories": [
            {
              "id": 3,
              "tenant_id": 100,
              "parent_id": 2,
              "name": "Espresso",
              "description": "Espresso-based products",
              "image_url": "https://example.com/espresso.jpg",
              "sort_order": 1,
              "is_active": true,
              "created_at": "2025-08-20T10:40:00Z",
              "updated_at": "2025-08-20T10:40:00Z",
              "sub_categories": []
            }
          ]
        }
      ]
    }
  ],
  "meta": null
}
```

---

### 14. Get Products by Category

Retrieves a paginated list of products belonging to a specific category.

**Endpoint:** `GET /api/v1/products/categories/{id}/products`

**Path Parameters:**
- `id`: Category ID (integer, required)

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 50, max: 100)
- `is_active`: Filter by active status (true/false)
- `sort`: Sort field (name, sku, selling_price, created_at)
- `order`: Sort order (asc, desc, default: asc)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Example Request:**
```
GET /api/v1/products/categories/1/products?page=1&limit=20&is_active=true&sort=name&order=asc
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Products retrieved successfully",
  "data": [
    {
      "id": 1,
      "tenant_id": 100,
      "category_id": 1,
      "sku": "PROD001",
      "barcode": "1234567890123",
      "name": "Premium Coffee Beans",
      "description": "High-quality arabica coffee beans",
      "unit": "kg",
      "cost_price": 50000.00,
      "selling_price": 75000.00,
      "min_stock": 10,
      "track_stock": true,
      "is_active": true,
      "images": [
        "https://example.com/product1.jpg"
      ],
      "variants": {
        "size": ["250g", "500g", "1kg"]
      },
      "created_at": "2025-08-20T10:30:00Z",
      "updated_at": "2025-08-20T10:30:00Z",
      "category": {
        "id": 1,
        "name": "Beverages",
        "description": "Coffee, tea, and other beverages"
      }
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1
  }
}
```

---

## Data Models

### Product Entity

Based on the database schema (`products` table), the product entity includes:

```sql
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    category_id BIGINT,
    sku VARCHAR(100) NOT NULL,
    barcode VARCHAR(100),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    unit VARCHAR(50) DEFAULT 'pcs',
    cost_price DECIMAL(12,2) DEFAULT 0.00,
    selling_price DECIMAL(12,2) NOT NULL,
    min_stock INTEGER DEFAULT 0,
    track_stock BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    images JSONB, -- Array of image URLs
    variants JSONB, -- Product variants (size, color, etc.)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES product_categories(id) ON DELETE SET NULL,
    UNIQUE (tenant_id, sku)
);
```

### Product Category Entity

Based on the database schema (`product_categories` table):

```sql
CREATE TABLE product_categories (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    parent_id BIGINT NULL, -- For sub-categories
    name VARCHAR(255) NOT NULL,
    description TEXT,
    image_url VARCHAR(500),
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES product_categories(id) ON DELETE SET NULL
);
```

### Key Relationships

1. **Tenant**: Each product/category belongs to exactly one tenant (multi-tenant isolation)
2. **Category**: Products can belong to a category (optional)
3. **Category Hierarchy**: Categories can have parent-child relationships for subcategories
4. **Stock**: Product stock is tracked per outlet through `product_stocks` table
5. **Transactions**: Products are referenced in sales transactions with snapshot data

---

## Business Rules

1. **Tenant Isolation**: All operations are scoped to the authenticated user's tenant
2. **Unique SKU**: Product SKU must be unique within a tenant
3. **Unique Barcode**: Product barcode must be unique within a tenant if provided
4. **Category Assignment**: Products can be assigned to categories within the same tenant
5. **Soft Delete**: Deleting products/categories sets `is_active = false` instead of hard deletion
6. **Active Filter**: By default, only active items (`is_active = true`) are returned
7. **Stock Tracking**: Products can have stock tracking enabled/disabled per product
8. **Category Hierarchy**: Categories support unlimited nesting levels
9. **Image Storage**: Product images are stored as JSON array of URLs
10. **Variants**: Product variants (size, color, etc.) stored as JSON for flexibility

---

## Error Handling

### Common Error Codes

- `400 Bad Request`: Invalid request format or validation errors
- `401 Unauthorized`: Missing or invalid JWT token
- `403 Forbidden`: User doesn't have permission to access/modify product
- `404 Not Found`: Product/Category not found or doesn't belong to user's tenant
- `409 Conflict`: Duplicate SKU or barcode within tenant
- `422 Unprocessable Entity`: Business logic validation failed
- `500 Internal Server Error`: Server-side error

### Validation Errors

The API performs validation on all input fields according to the defined rules and returns detailed field-level validation errors in the standard error response format.

---

## Multi-Tenant Considerations

1. **Automatic Tenant Filtering**: All queries automatically filter by the authenticated user's tenant ID
2. **Cross-Tenant Prevention**: Users cannot access products/categories from other tenants
3. **SKU/Barcode Uniqueness**: SKU and barcode are unique within each tenant, not globally
4. **Category Relationships**: Category parent-child relationships are restricted within the same tenant

---

## Stock Management Integration

Products integrate with the stock management system:

1. **Product Stocks**: Stock quantities are maintained per product per outlet in `product_stocks` table
2. **Stock Tracking**: Products with `track_stock = true` will have their stock automatically updated during transactions
3. **Minimum Stock**: Products can have minimum stock levels for low stock alerts
4. **Stock Movements**: All stock changes are logged in `stock_movements` table for audit trails

---

## Performance Considerations

1. **Indexes**: Database indexes are optimized for tenant-based queries and common search patterns
2. **Pagination**: All list endpoints support pagination to handle large datasets
3. **Selective Fields**: Consider implementing field selection to reduce payload size for mobile applications
4. **Caching**: Category hierarchies can be cached as they change infrequently
5. **Image Optimization**: Product images should be optimized and served via CDN for better performance