# Outlets API Documentation

This document describes the CRUD endpoints for the outlets module.

## Endpoints

All endpoints require JWT authentication and proper tenant context.

### Base URL
```
/api/v1/outlets
```

## Outlet CRUD Operations

### 1. Create Outlet
**POST** `/api/v1/outlets`

**Request Body:**
```json
{
  "name": "Main Store",
  "code": "MAIN001",
  "description": "Main outlet in downtown",
  "address": "123 Main Street",
  "city": "Jakarta",
  "province": "DKI Jakarta",
  "postal_code": "12345",
  "phone": "+62123456789",
  "email": "main@example.com",
  "manager_id": 1,
  "settings": {
    "tax_rate": 10,
    "auto_print": true
  }
}
```

**Response (201):**
```json
{
  "id": 1,
  "tenant_id": 1,
  "name": "Main Store",
  "code": "MAIN001",
  "description": "Main outlet in downtown",
  "address": "123 Main Street",
  "city": "Jakarta",
  "province": "DKI Jakarta",
  "postal_code": "12345",
  "phone": "+62123456789",
  "email": "main@example.com",
  "manager_id": 1,
  "is_active": true,
  "settings": {
    "tax_rate": 10,
    "auto_print": true
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "manager": {
    "id": 1,
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "+62987654321"
  }
}
```

### 2. Get All Outlets
**GET** `/api/v1/outlets`

**Query Parameters:**
- `page` (int, default: 1) - Page number
- `limit` (int, default: 50, max: 100) - Records per page
- `name` (string) - Filter by name (partial match)
- `code` (string) - Filter by code (partial match)
- `city` (string) - Filter by city (partial match)
- `province` (string) - Filter by province (partial match)
- `manager_id` (int) - Filter by manager ID
- `is_active` (bool) - Filter by active status
- `sort` (string) - Sort field (name, code, city, created_at)
- `order` (string) - Sort order (ASC, DESC)

**Example:** `/api/v1/outlets?page=1&limit=10&city=Jakarta&sort=name&order=ASC`

**Response (200):**
```json
{
  "outlets": [
    {
      "id": 1,
      "tenant_id": 1,
      "name": "Main Store",
      "code": "MAIN001",
      "description": "Main outlet in downtown",
      "address": "123 Main Street",
      "city": "Jakarta",
      "province": "DKI Jakarta",
      "postal_code": "12345",
      "phone": "+62123456789",
      "email": "main@example.com",
      "manager_id": 1,
      "is_active": true,
      "settings": {
        "tax_rate": 10,
        "auto_print": true
      },
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "manager": {
        "id": 1,
        "full_name": "John Doe",
        "email": "john@example.com",
        "phone": "+62987654321"
      }
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 10
}
```

### 3. Get Outlet by ID
**GET** `/api/v1/outlets/{id}`

**Response (200):**
```json
{
  "id": 1,
  "tenant_id": 1,
  "name": "Main Store",
  "code": "MAIN001",
  "description": "Main outlet in downtown",
  "address": "123 Main Street",
  "city": "Jakarta",
  "province": "DKI Jakarta",
  "postal_code": "12345",
  "phone": "+62123456789",
  "email": "main@example.com",
  "manager_id": 1,
  "is_active": true,
  "settings": {
    "tax_rate": 10,
    "auto_print": true
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "manager": {
    "id": 1,
    "full_name": "John Doe",
    "email": "john@example.com",
    "phone": "+62987654321"
  }
}
```

### 4. Get Outlet by Code
**GET** `/api/v1/outlets/code/{code}`

**Response (200):**
Same as Get Outlet by ID

### 5. Update Outlet
**PUT** `/api/v1/outlets/{id}`

**Request Body:**
```json
{
  "name": "Updated Main Store",
  "code": "MAIN001",
  "description": "Updated description",
  "address": "123 Updated Street",
  "city": "Jakarta",
  "province": "DKI Jakarta",
  "postal_code": "12345",
  "phone": "+62123456789",
  "email": "updated@example.com",
  "manager_id": 2,
  "is_active": true,
  "settings": {
    "tax_rate": 15,
    "auto_print": false
  }
}
```

**Response (200):**
Same as Create Outlet response with updated values

### 6. Delete Outlet
**DELETE** `/api/v1/outlets/{id}`

**Response (200):**
```json
{
  "message": "Outlet deleted successfully"
}
```

## Error Responses

**400 Bad Request:**
```json
{
  "error": "outlet with this code already exists"
}
```

**404 Not Found:**
```json
{
  "error": "Outlet not found"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Failed to get outlets"
}
```

## Validation Rules

### Create/Update Outlet Request:
- `name`: Required, 1-255 characters
- `code`: Required, 1-50 characters, unique per tenant
- `description`: Optional
- `address`: Optional
- `city`: Optional, max 100 characters
- `province`: Optional, max 100 characters
- `postal_code`: Optional, max 10 characters
- `phone`: Optional, max 20 characters
- `email`: Optional, valid email format, max 255 characters
- `manager_id`: Optional, must reference existing user
- `settings`: Optional JSON object

## Features

1. **Multi-tenant Support**: All operations are tenant-isolated
2. **Manager Information**: Outlets can be associated with a manager (user)
3. **Flexible Settings**: JSON settings field for outlet-specific configurations
4. **Search and Filter**: Multiple filter options for listing outlets
5. **Pagination**: Configurable page size with limits
6. **Validation**: Comprehensive input validation
7. **Unique Codes**: Outlet codes are unique within each tenant

## Implementation Details

The outlets module follows the Domain-Driven Design (DDD) pattern with:

- **Domain Layer**: Entities, DTOs, and interfaces
- **Persistence Layer**: Models, repositories, and database operations
- **Service Layer**: Business logic and validation
- **Handler Layer**: HTTP endpoints and request/response handling

The implementation includes:
- GORM for ORM and database operations
- Echo for HTTP routing and middleware
- JWT authentication and tenant isolation
- JSON validation using struct tags
- Error handling with appropriate HTTP status codes