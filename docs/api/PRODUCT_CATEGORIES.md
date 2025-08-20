# Product Categories API Documentation

This document provides comprehensive API documentation for the Product Categories module of ExVen POS Lite system.

## Overview

The Product Categories API manages product categorization within a multi-tenant POS system. Each category belongs to a specific tenant and supports hierarchical structures with parent-child relationships. Categories are used to organize products and enable better inventory management and reporting.

## Base URL

All product categories API endpoints are prefixed with `/api/v1/products/categories`

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

## Endpoints

### 1. Create Product Category

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
  "name": "Electronics",
  "description": "Electronic devices and accessories",
  "image_url": "https://example.com/images/electronics.jpg",
  "parent_id": null,
  "sort_order": 1
}
```

**Validation Rules:**
- `name`: Required, 1-255 characters
- `description`: Optional
- `image_url`: Optional, must be valid URL format
- `parent_id`: Optional, must be valid category ID within tenant
- `sort_order`: Optional, integer for display ordering

**Response:**

*Success (201 Created):*
```json
{
  "message": "Category created successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "parent_id": null,
    "name": "Electronics",
    "description": "Electronic devices and accessories",
    "image_url": "https://example.com/images/electronics.jpg",
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

*Error (400 Bad Request):*
```json
{
  "message": "Validation failed",
  "data": null,
  "errors": {
    "name": ["Name is required"],
    "image_url": ["Invalid URL format"]
  }
}
```

---

### 2. Get All Product Categories

Retrieves a paginated list of product categories for the authenticated tenant with optional filtering and sorting.

**Endpoint:** `GET /api/v1/products/categories`

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 50, max: 100)
- `parent_id`: Filter by parent category ID (use "null" or empty for root categories)
- `is_active`: Filter by active status (true/false)
- `sort`: Sort field (name, sort_order, created_at)
- `order`: Sort order (asc, desc, default: asc)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Example Request:**
```
GET /api/v1/products/categories?page=1&limit=20&parent_id=null&is_active=true&sort=sort_order&order=asc
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
      "name": "Electronics",
      "description": "Electronic devices and accessories",
      "image_url": "https://example.com/images/electronics.jpg",
      "sort_order": 1,
      "is_active": true,
      "created_at": "2025-08-20T10:30:00Z",
      "updated_at": "2025-08-20T10:30:00Z",
      "parent": null,
      "sub_categories": []
    },
    {
      "id": 2,
      "tenant_id": 100,
      "parent_id": 1,
      "name": "Smartphones",
      "description": "Mobile phones and accessories",
      "image_url": "https://example.com/images/smartphones.jpg",
      "sort_order": 1,
      "is_active": true,
      "created_at": "2025-08-20T10:31:00Z",
      "updated_at": "2025-08-20T10:31:00Z",
      "parent": {
        "id": 1,
        "tenant_id": 100,
        "parent_id": null,
        "name": "Electronics",
        "description": "Electronic devices and accessories",
        "image_url": "https://example.com/images/electronics.jpg",
        "sort_order": 1,
        "is_active": true,
        "created_at": "2025-08-20T10:30:00Z",
        "updated_at": "2025-08-20T10:30:00Z"
      }
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 2
  }
}
```

---

### 3. Get Product Category by ID

Retrieves a specific product category by its ID, including parent and subcategories if they exist.

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
    "name": "Electronics",
    "description": "Electronic devices and accessories",
    "image_url": "https://example.com/images/electronics.jpg",
    "sort_order": 1,
    "is_active": true,
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T10:30:00Z",
    "sub_categories": [
      {
        "id": 2,
        "tenant_id": 100,
        "parent_id": 1,
        "name": "Smartphones",
        "description": "Mobile phones and accessories",
        "image_url": "https://example.com/images/smartphones.jpg",
        "sort_order": 1,
        "is_active": true,
        "created_at": "2025-08-20T10:31:00Z",
        "updated_at": "2025-08-20T10:31:00Z"
      }
    ]
  },
  "meta": null
}
```

*Error (404 Not Found):*
```json
{
  "message": "Category not found",
  "data": null,
  "errors": {}
}
```

*Error (400 Bad Request):*
```json
{
  "message": "Invalid category ID",
  "data": null,
  "errors": {}
}
```

---

### 4. Update Product Category

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
  "name": "Electronics Updated",
  "description": "Updated description for electronic devices",
  "image_url": "https://example.com/images/electronics-updated.jpg",
  "parent_id": null,
  "sort_order": 2,
  "is_active": true
}
```

**Validation Rules:**
- Same as Create Category request
- `is_active`: Boolean flag to activate/deactivate category

**Response:**

*Success (200 OK):*
```json
{
  "message": "Category updated successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "parent_id": null,
    "name": "Electronics Updated",
    "description": "Updated description for electronic devices",
    "image_url": "https://example.com/images/electronics-updated.jpg",
    "sort_order": 2,
    "is_active": true,
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T11:30:00Z",
    "parent": null,
    "sub_categories": []
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
    "name": ["Name is required"],
    "parent_id": ["Cannot set category as its own parent"]
  }
}
```

---

### 5. Delete Product Category

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
    "category": ["Cannot delete category with active subcategories or products"]
  }
}
```

*Error (404 Not Found):*
```json
{
  "message": "Category not found",
  "data": null,
  "errors": {}
}
```

---

### 6. Get Category Hierarchy

Retrieves the complete category hierarchy as a nested tree structure for the authenticated tenant.

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
      "name": "Electronics",
      "description": "Electronic devices and accessories",
      "image_url": "https://example.com/images/electronics.jpg",
      "sort_order": 1,
      "is_active": true,
      "created_at": "2025-08-20T10:30:00Z",
      "updated_at": "2025-08-20T10:30:00Z",
      "sub_categories": [
        {
          "id": 2,
          "tenant_id": 100,
          "parent_id": 1,
          "name": "Smartphones",
          "description": "Mobile phones and accessories",
          "image_url": "https://example.com/images/smartphones.jpg",
          "sort_order": 1,
          "is_active": true,
          "created_at": "2025-08-20T10:31:00Z",
          "updated_at": "2025-08-20T10:31:00Z",
          "sub_categories": [
            {
              "id": 3,
              "tenant_id": 100,
              "parent_id": 2,
              "name": "iPhone",
              "description": "Apple iPhone devices",
              "image_url": "https://example.com/images/iphone.jpg",
              "sort_order": 1,
              "is_active": true,
              "created_at": "2025-08-20T10:32:00Z",
              "updated_at": "2025-08-20T10:32:00Z",
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

### 7. Get Products by Category

Retrieves a paginated list of products that belong to a specific category.

**Endpoint:** `GET /api/v1/products/categories/{id}/products`

**Path Parameters:**
- `id`: Category ID (integer, required)

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 50, max: 100)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Example Request:**
```
GET /api/v1/products/categories/1/products?page=1&limit=20
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
      "sku": "ELEC001",
      "barcode": "1234567890123",
      "name": "Smartphone XYZ",
      "description": "Latest smartphone with advanced features",
      "unit": "pcs",
      "cost_price": 500.00,
      "selling_price": 750.00,
      "min_stock": 10,
      "track_stock": true,
      "is_active": true,
      "images": [
        "https://example.com/images/smartphone-1.jpg",
        "https://example.com/images/smartphone-2.jpg"
      ],
      "variants": {
        "color": ["black", "white", "blue"],
        "storage": ["64GB", "128GB", "256GB"]
      },
      "created_at": "2025-08-20T10:33:00Z",
      "updated_at": "2025-08-20T10:33:00Z",
      "category": {
        "id": 1,
        "tenant_id": 100,
        "parent_id": null,
        "name": "Electronics",
        "description": "Electronic devices and accessories",
        "image_url": "https://example.com/images/electronics.jpg",
        "sort_order": 1,
        "is_active": true,
        "created_at": "2025-08-20T10:30:00Z",
        "updated_at": "2025-08-20T10:30:00Z"
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

*Error (400 Bad Request):*
```json
{
  "message": "Invalid category ID",
  "data": null,
  "errors": {}
}
```

---

## Data Models

### Product Category Entity

Based on the database schema (`product_categories` table), the category entity includes:

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

1. **Tenant**: Each category belongs to exactly one tenant (multi-tenant isolation)
2. **Hierarchical Structure**: Categories can have parent-child relationships for organization
3. **Products**: Products can be assigned to categories for organization and filtering
4. **Self-Referencing**: Categories can reference other categories as parents

---

## Business Rules

1. **Tenant Isolation**: All operations are scoped to the authenticated user's tenant
2. **Hierarchical Organization**: Categories support unlimited nesting levels
3. **Sort Order**: Categories can be ordered using the `sort_order` field
4. **Soft Delete**: Deleting a category sets `is_active = false` instead of hard deletion
5. **Active Filter**: By default, only active categories (`is_active = true`) are returned
6. **Dependency Prevention**: Categories with subcategories or products cannot be deleted
7. **Parent Validation**: Parent category must belong to the same tenant
8. **Self-Reference Prevention**: Category cannot be set as its own parent

---

## Error Handling

### Common Error Codes

- `400 Bad Request`: Invalid request format or validation errors
- `401 Unauthorized`: Missing or invalid JWT token
- `403 Forbidden`: User doesn't have permission to access/modify category
- `404 Not Found`: Category not found or doesn't belong to user's tenant
- `409 Conflict`: Circular reference in parent-child relationship
- `422 Unprocessable Entity`: Business logic validation failed
- `500 Internal Server Error`: Server-side error

### Validation Errors

The API performs validation on all input fields according to the defined rules and returns detailed field-level validation errors in the standard error response format.

---

## Multi-Tenant Considerations

1. **Automatic Tenant Filtering**: All queries automatically filter by the authenticated user's tenant ID
2. **Cross-Tenant Prevention**: Users cannot access categories from other tenants
3. **Parent Validation**: Parent category must belong to the same tenant as the child
4. **Hierarchy Isolation**: Category hierarchies are completely isolated per tenant

---

## Usage Examples

### Creating a Category Hierarchy

1. **Create Root Category:**
   ```
   POST /api/v1/products/categories
   {
     "name": "Electronics",
     "parent_id": null,
     "sort_order": 1
   }
   ```

2. **Create Subcategory:**
   ```
   POST /api/v1/products/categories
   {
     "name": "Smartphones",
     "parent_id": 1,
     "sort_order": 1
   }
   ```

3. **Create Sub-subcategory:**
   ```
   POST /api/v1/products/categories
   {
     "name": "iPhone",
     "parent_id": 2,
     "sort_order": 1
   }
   ```

### Filtering Categories

- **Get only root categories:**
  ```
  GET /api/v1/products/categories?parent_id=null
  ```

- **Get subcategories of a specific category:**
  ```
  GET /api/v1/products/categories?parent_id=1
  ```

- **Get only active categories ordered by sort_order:**
  ```
  GET /api/v1/products/categories?is_active=true&sort=sort_order&order=asc
  ```