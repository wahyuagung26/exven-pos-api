# Outlets API Documentation

This document provides comprehensive API documentation for the Outlets module of ExVen POS Lite system.

## Overview

The Outlets API manages outlet/store locations within a multi-tenant POS system. Each outlet belongs to a specific tenant and can have multiple users assigned to it. Outlets are used to track sales transactions and inventory across different physical locations.

## Base URL

All outlets API endpoints are prefixed with `/api/v1/outlets`

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

### 1. Create Outlet

Creates a new outlet for the authenticated tenant.

**Endpoint:** `POST /api/v1/outlets`

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "name": "Main Store",
  "code": "MAIN001",
  "description": "Main outlet located downtown",
  "address": "123 Main Street",
  "city": "Jakarta",
  "province": "DKI Jakarta",
  "postal_code": "12345",
  "phone": "+628123456789",
  "email": "main@example.com",
  "manager_id": 123,
  "settings": {
    "printer_ip": "192.168.1.100",
    "tax_enabled": true,
    "tax_rate": 10.5
  }
}
```

**Validation Rules:**
- `name`: Required, 1-255 characters
- `code`: Required, 1-50 characters, unique within tenant
- `description`: Optional
- `address`: Optional
- `city`: Optional, max 100 characters
- `province`: Optional, max 100 characters  
- `postal_code`: Optional, max 10 characters
- `phone`: Optional, max 20 characters
- `email`: Optional, valid email format, max 255 characters
- `manager_id`: Optional, must be valid user ID within tenant
- `settings`: Optional, JSON object for outlet-specific configurations

**Response:**

*Success (201 Created):*
```json
{
  "message": "Outlet created successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "name": "Main Store",
    "code": "MAIN001",
    "description": "Main outlet located downtown",
    "address": "123 Main Street",
    "city": "Jakarta",
    "province": "DKI Jakarta",
    "postal_code": "12345",
    "phone": "+628123456789",
    "email": "main@example.com",
    "manager_id": 123,
    "is_active": true,
    "settings": {
      "printer_ip": "192.168.1.100",
      "tax_enabled": true,
      "tax_rate": 10.5
    },
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T10:30:00Z",
    "manager": {
      "id": 123,
      "full_name": "John Doe",
      "email": "john@example.com",
      "phone": "+628123456789"
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
    "name": ["Name is required"],
    "email": ["Invalid email format"]
  }
}
```

---

### 2. Get All Outlets

Retrieves a paginated list of outlets for the authenticated tenant with optional filtering and sorting.

**Endpoint:** `GET /api/v1/outlets`

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 50, max: 100)
- `name`: Filter by outlet name (partial match)
- `code`: Filter by outlet code (partial match)
- `city`: Filter by city (exact match)
- `province`: Filter by province (exact match)
- `manager_id`: Filter by manager user ID
- `is_active`: Filter by active status (true/false)
- `sort`: Sort field (name, code, city, province, created_at)
- `order`: Sort order (asc, desc, default: asc)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Example Request:**
```
GET /api/v1/outlets?page=1&limit=20&city=Jakarta&is_active=true&sort=name&order=asc
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Outlets retrieved successfully",
  "data": [
    {
      "id": 1,
      "tenant_id": 100,
      "name": "Main Store",
      "code": "MAIN001",
      "description": "Main outlet located downtown",
      "address": "123 Main Street",
      "city": "Jakarta",
      "province": "DKI Jakarta",
      "postal_code": "12345",
      "phone": "+628123456789",
      "email": "main@example.com",
      "manager_id": 123,
      "is_active": true,
      "settings": {
        "printer_ip": "192.168.1.100",
        "tax_enabled": true,
        "tax_rate": 10.5
      },
      "created_at": "2025-08-20T10:30:00Z",
      "updated_at": "2025-08-20T10:30:00Z",
      "manager": {
        "id": 123,
        "full_name": "John Doe",
        "email": "john@example.com",
        "phone": "+628123456789"
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

### 3. Get Outlet by ID

Retrieves a specific outlet by its ID.

**Endpoint:** `GET /api/v1/outlets/{id}`

**Path Parameters:**
- `id`: Outlet ID (integer, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Outlet retrieved successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "name": "Main Store",
    "code": "MAIN001",
    "description": "Main outlet located downtown",
    "address": "123 Main Street",
    "city": "Jakarta",
    "province": "DKI Jakarta",
    "postal_code": "12345",
    "phone": "+628123456789",
    "email": "main@example.com",
    "manager_id": 123,
    "is_active": true,
    "settings": {
      "printer_ip": "192.168.1.100",
      "tax_enabled": true,
      "tax_rate": 10.5
    },
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T10:30:00Z",
    "manager": {
      "id": 123,
      "full_name": "John Doe",
      "email": "john@example.com",
      "phone": "+628123456789"
    }
  },
  "meta": null
}
```

*Error (404 Not Found):*
```json
{
  "message": "Outlet not found",
  "data": null,
  "errors": {}
}
```

*Error (400 Bad Request):*
```json
{
  "message": "Invalid outlet ID",
  "data": null,
  "errors": {}
}
```

---

### 4. Update Outlet

Updates an existing outlet.

**Endpoint:** `PUT /api/v1/outlets/{id}`

**Path Parameters:**
- `id`: Outlet ID (integer, required)

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "name": "Main Store Updated",
  "code": "MAIN001",
  "description": "Updated main outlet description",
  "address": "456 Main Street",
  "city": "Jakarta",
  "province": "DKI Jakarta",
  "postal_code": "12346",
  "phone": "+628123456790",
  "email": "main.updated@example.com",
  "manager_id": 124,
  "is_active": true,
  "settings": {
    "printer_ip": "192.168.1.101",
    "tax_enabled": true,
    "tax_rate": 11.0
  }
}
```

**Validation Rules:**
- Same as Create Outlet request
- `is_active`: Boolean flag to activate/deactivate outlet

**Response:**

*Success (200 OK):*
```json
{
  "message": "Outlet updated successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "name": "Main Store Updated",
    "code": "MAIN001",
    "description": "Updated main outlet description",
    "address": "456 Main Street",
    "city": "Jakarta",
    "province": "DKI Jakarta",
    "postal_code": "12346",
    "phone": "+628123456790",
    "email": "main.updated@example.com",
    "manager_id": 124,
    "is_active": true,
    "settings": {
      "printer_ip": "192.168.1.101",
      "tax_enabled": true,
      "tax_rate": 11.0
    },
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T11:30:00Z",
    "manager": {
      "id": 124,
      "full_name": "Jane Smith",
      "email": "jane@example.com",
      "phone": "+628123456790"
    }
  },
  "meta": null
}
```

---

### 5. Delete Outlet

Soft deletes an outlet by setting `is_active` to false.

**Endpoint:** `DELETE /api/v1/outlets/{id}`

**Path Parameters:**
- `id`: Outlet ID (integer, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Outlet deleted successfully",
  "data": null,
  "meta": null
}
```

*Error (400 Bad Request):*
```json
{
  "message": "Cannot delete outlet",
  "data": null,
  "errors": {
    "outlet": ["Cannot delete outlet with active transactions"]
  }
}
```

---

### 6. Get Outlet by Code

Retrieves a specific outlet by its unique code within the tenant.

**Endpoint:** `GET /api/v1/outlets/code/{code}`

**Path Parameters:**
- `code`: Outlet code (string, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Outlet retrieved successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "name": "Main Store",
    "code": "MAIN001",
    "description": "Main outlet located downtown",
    "address": "123 Main Street",
    "city": "Jakarta",
    "province": "DKI Jakarta",
    "postal_code": "12345",
    "phone": "+628123456789",
    "email": "main@example.com",
    "manager_id": 123,
    "is_active": true,
    "settings": {
      "printer_ip": "192.168.1.100",
      "tax_enabled": true,
      "tax_rate": 10.5
    },
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T10:30:00Z",
    "manager": {
      "id": 123,
      "full_name": "John Doe",
      "email": "john@example.com",
      "phone": "+628123456789"
    }
  },
  "meta": null
}
```

---

## Data Models

### Outlet Entity

Based on the database schema (`outlets` table), the outlet entity includes:

```sql
CREATE TABLE outlets (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL, -- Unique within tenant
    description TEXT,
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),
    phone VARCHAR(20),
    email VARCHAR(255),
    manager_id BIGINT, -- References users table
    is_active BOOLEAN DEFAULT TRUE,
    settings JSONB, -- Outlet-specific configurations
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE (tenant_id, code)
);
```

### Key Relationships

1. **Tenant**: Each outlet belongs to exactly one tenant (multi-tenant isolation)
2. **Manager**: Each outlet can have one assigned manager (optional)
3. **Users**: Multiple users can be assigned to an outlet through `user_outlets` table
4. **Products**: Product stock is tracked per outlet through `product_stocks` table
5. **Transactions**: Sales transactions are recorded per outlet

---

## Business Rules

1. **Tenant Isolation**: All operations are scoped to the authenticated user's tenant
2. **Unique Code**: Outlet code must be unique within a tenant
3. **Manager Assignment**: Manager must be a valid user within the same tenant
4. **Soft Delete**: Deleting an outlet sets `is_active = false` instead of hard deletion
5. **Active Filter**: By default, only active outlets (`is_active = true`) are returned
6. **Transaction Dependency**: Outlets with active transactions cannot be deleted
7. **Settings**: JSON field for storing outlet-specific configurations (printer settings, tax rates, etc.)

---

## Error Handling

### Common Error Codes

- `400 Bad Request`: Invalid request format or validation errors
- `401 Unauthorized`: Missing or invalid JWT token
- `403 Forbidden`: User doesn't have permission to access/modify outlet
- `404 Not Found`: Outlet not found or doesn't belong to user's tenant
- `409 Conflict`: Duplicate outlet code within tenant
- `422 Unprocessable Entity`: Business logic validation failed
- `500 Internal Server Error`: Server-side error

### Validation Errors

The API performs validation on all input fields according to the defined rules and returns detailed field-level validation errors in the standard error response format.

---

## Multi-Tenant Considerations

1. **Automatic Tenant Filtering**: All queries automatically filter by the authenticated user's tenant ID
2. **Cross-Tenant Prevention**: Users cannot access outlets from other tenants
3. **Manager Validation**: Manager must belong to the same tenant as the outlet
4. **Code Uniqueness**: Outlet codes are unique within each tenant, not globally