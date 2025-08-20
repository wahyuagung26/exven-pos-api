# Products API Endpoints

This document describes the CRUD endpoints for the products table based on the implementation.

## Base URL
```
/api/v1/products
```

## Product Endpoints

### 1. Create Product
```
POST /api/v1/products
```

**Request Body:**
```json
{
  "category_id": 1,
  "sku": "PROD-001",
  "barcode": "1234567890",
  "name": "Sample Product",
  "description": "Product description",
  "unit": "pcs",
  "cost_price": 10.00,
  "selling_price": 15.00,
  "min_stock": 5,
  "track_stock": true,
  "images": ["https://example.com/image1.jpg"],
  "variants": {"size": "M", "color": "red"}
}
```

### 2. Get All Products
```
GET /api/v1/products?page=1&limit=50
```

### 3. Get Product by ID
```
GET /api/v1/products/:id
```

### 4. Update Product
```
PUT /api/v1/products/:id
```

**Request Body:** (same as create, plus `is_active`)
```json
{
  "category_id": 1,
  "sku": "PROD-001-UPDATED",
  "barcode": "1234567890",
  "name": "Updated Product",
  "description": "Updated description",
  "unit": "pcs",
  "cost_price": 12.00,
  "selling_price": 18.00,
  "min_stock": 10,
  "track_stock": true,
  "is_active": true,
  "images": ["https://example.com/image1.jpg"],
  "variants": {"size": "L", "color": "blue"}
}
```

### 5. Delete Product
```
DELETE /api/v1/products/:id
```

### 6. Get Product by SKU
```
GET /api/v1/products/sku/:sku
```

### 7. Get Product by Barcode
```
GET /api/v1/products/barcode/:barcode
```

## Product Category Endpoints

### 8. Create Category
```
POST /api/v1/products/categories
```

### 9. Get All Categories
```
GET /api/v1/products/categories?page=1&limit=50
```

### 10. Get Category by ID
```
GET /api/v1/products/categories/:id
```

### 11. Update Category
```
PUT /api/v1/products/categories/:id
```

### 12. Delete Category
```
DELETE /api/v1/products/categories/:id
```

### 13. Get Category Hierarchy
```
GET /api/v1/products/categories/hierarchy
```

### 14. Get Products by Category
```
GET /api/v1/products/categories/:id/products?page=1&limit=50
```

## Features Implemented

✅ **Complete CRUD operations** for products table
✅ **Multi-tenant support** with tenant_id filtering
✅ **Input validation** using struct tags
✅ **Pagination support** for list endpoints
✅ **Category relationships** with product-category associations
✅ **SKU and barcode lookup** for quick product search
✅ **JSON fields support** for images and variants
✅ **Domain-Driven Design** architecture with proper separation of concerns
✅ **Repository pattern** with GORM integration
✅ **Service layer** with business logic validation
✅ **HTTP handlers** with proper error handling
✅ **Dependency injection** through module registration

## Database Schema Compliance

The implementation follows the `products` table schema from `schema.sql`:
- ✅ All table columns mapped correctly
- ✅ Tenant isolation implemented
- ✅ Foreign key relationships with categories
- ✅ JSON fields for images and variants
- ✅ Proper indexing through GORM tags
- ✅ Data validation and constraints