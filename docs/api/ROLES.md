# Roles API Documentation

This document provides comprehensive API documentation for the Roles module of ExVen POS Lite system.

## Overview

The Roles API manages user roles within the POS system. Roles define user permissions and access levels, including system-defined roles (Super Admin, Tenant Owner, Manager, Cashier) and custom roles created by tenant owners. The system uses role-based access control (RBAC) to secure API endpoints and features.

## Base URL

All roles API endpoints are prefixed with `/api/v1/roles`

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

### 1. Get All Roles

Retrieves a paginated list of all roles available in the system.

**Endpoint:** `GET /api/v1/roles`

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 50, max: 100)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Example Request:**
```
GET /api/v1/roles?page=1&limit=10
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Roles retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "super_admin",
      "display_name": "Super Admin",
      "description": "Full system access",
      "permissions": ["*"],
      "is_system": true,
      "created_at": "2025-08-20T10:30:00Z"
    },
    {
      "id": 2,
      "name": "tenant_owner",
      "display_name": "Pemilik Bisnis",
      "description": "Full access within tenant",
      "permissions": ["tenant.*"],
      "is_system": true,
      "created_at": "2025-08-20T10:30:00Z"
    },
    {
      "id": 3,
      "name": "manager",
      "display_name": "Manager",
      "description": "Mengelola outlet dan laporan",
      "permissions": ["outlet.*", "reports.*", "products.*", "customers.*"],
      "is_system": true,
      "created_at": "2025-08-20T10:30:00Z"
    },
    {
      "id": 4,
      "name": "cashier",
      "display_name": "Kasir",
      "description": "Melakukan penjualan",
      "permissions": ["sales.*", "customers.read", "products.read"],
      "is_system": true,
      "created_at": "2025-08-20T10:30:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 10,
    "total": 4
  }
}
```

*Error (500 Internal Server Error):*
```json
{
  "message": "Failed to get roles",
  "data": null,
  "errors": {}
}
```

---

### 2. Get Role by ID

Retrieves a specific role by its ID.

**Endpoint:** `GET /api/v1/roles/{id}`

**Path Parameters:**
- `id`: Role ID (integer, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Role retrieved successfully",
  "data": {
    "id": 2,
    "name": "tenant_owner",
    "display_name": "Pemilik Bisnis",
    "description": "Full access within tenant",
    "permissions": ["tenant.*"],
    "is_system": true,
    "created_at": "2025-08-20T10:30:00Z"
  },
  "meta": null
}
```

*Error (404 Not Found):*
```json
{
  "message": "Role not found",
  "data": null,
  "errors": {}
}
```

*Error (400 Bad Request):*
```json
{
  "message": "Invalid role ID",
  "data": null,
  "errors": {}
}
```

---

### 3. Get Role by Name

Retrieves a specific role by its unique name.

**Endpoint:** `GET /api/v1/roles/name/{name}`

**Path Parameters:**
- `name`: Role name (string, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Example Request:**
```
GET /api/v1/roles/name/manager
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Role retrieved successfully",
  "data": {
    "id": 3,
    "name": "manager",
    "display_name": "Manager",
    "description": "Mengelola outlet dan laporan",
    "permissions": ["outlet.*", "reports.*", "products.*", "customers.*"],
    "is_system": true,
    "created_at": "2025-08-20T10:30:00Z"
  },
  "meta": null
}
```

*Error (404 Not Found):*
```json
{
  "message": "Role not found",
  "data": null,
  "errors": {}
}
```

---

### 4. Get System Roles

Retrieves all system-defined roles (roles that cannot be modified or deleted).

**Endpoint:** `GET /api/v1/roles/system`

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "System roles retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "super_admin",
      "display_name": "Super Admin",
      "description": "Full system access",
      "permissions": ["*"],
      "is_system": true,
      "created_at": "2025-08-20T10:30:00Z"
    },
    {
      "id": 2,
      "name": "tenant_owner",
      "display_name": "Pemilik Bisnis",
      "description": "Full access within tenant",
      "permissions": ["tenant.*"],
      "is_system": true,
      "created_at": "2025-08-20T10:30:00Z"
    },
    {
      "id": 3,
      "name": "manager",
      "display_name": "Manager",
      "description": "Mengelola outlet dan laporan",
      "permissions": ["outlet.*", "reports.*", "products.*", "customers.*"],
      "is_system": true,
      "created_at": "2025-08-20T10:30:00Z"
    },
    {
      "id": 4,
      "name": "cashier",
      "display_name": "Kasir",
      "description": "Melakukan penjualan",
      "permissions": ["sales.*", "customers.read", "products.read"],
      "is_system": true,
      "created_at": "2025-08-20T10:30:00Z"
    }
  ],
  "meta": null
}
```

*Error (500 Internal Server Error):*
```json
{
  "message": "Failed to get system roles",
  "data": null,
  "errors": {}
}
```

---

## Data Models

### Role Entity

Based on the database schema (`roles` table), the role entity includes:

```sql
CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    permissions JSONB, -- Array permission strings
    is_system BOOLEAN DEFAULT FALSE, -- System roles cannot be deleted
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Role Response Object

```json
{
  "id": 1,
  "name": "role_name",
  "display_name": "Human Readable Name",
  "description": "Role description",
  "permissions": ["permission1", "permission2"],
  "is_system": true,
  "created_at": "2025-08-20T10:30:00Z"
}
```

### Field Descriptions

- `id`: Unique identifier for the role
- `name`: Unique system name for the role (lowercase, underscore-separated)
- `display_name`: Human-readable name for display purposes
- `description`: Optional description of the role's purpose
- `permissions`: Array of permission strings that define what the role can access
- `is_system`: Boolean flag indicating if this is a system-defined role (cannot be modified/deleted)
- `created_at`: Timestamp when the role was created (ISO 8601 format)

---

## System Roles

The system includes four predefined roles that cannot be modified or deleted:

### 1. Super Admin
- **Name:** `super_admin`
- **Permissions:** `["*"]` (full system access)
- **Use Case:** System administrators with complete access to all features and tenants

### 2. Tenant Owner
- **Name:** `tenant_owner`
- **Permissions:** `["tenant.*"]` (full access within tenant)
- **Use Case:** Business owners with complete control over their tenant's data and settings

### 3. Manager
- **Name:** `manager`
- **Permissions:** `["outlet.*", "reports.*", "products.*", "customers.*"]`
- **Use Case:** Store managers who can manage outlets, view reports, and manage inventory

### 4. Cashier
- **Name:** `cashier`
- **Permissions:** `["sales.*", "customers.read", "products.read"]`
- **Use Case:** Front-line staff who can process sales and view customer/product information

---

## Permission System

The permission system uses a hierarchical dot notation:

### Permission Patterns
- `*`: Full system access (super admin only)
- `tenant.*`: Full tenant access
- `outlet.*`: Full outlet management
- `sales.*`: Full sales operations
- `products.*`: Full product management
- `customers.*`: Full customer management
- `reports.*`: Full reporting access
- `[resource].read`: Read-only access to a resource
- `[resource].write`: Write access to a resource
- `[resource].delete`: Delete access to a resource

### Permission Hierarchy
- Wildcard permissions (`*`) grant access to all sub-permissions
- Module permissions (`tenant.*`) grant access to all operations within that module
- Specific permissions (`customers.read`) grant access to specific operations only

---

## Business Rules

1. **System Roles**: System roles (is_system = true) cannot be modified or deleted
2. **Unique Names**: Role names must be unique across the entire system
3. **Permission Validation**: All permissions must follow the defined permission pattern
4. **Role Assignment**: Users can only be assigned roles that exist in the system
5. **Hierarchical Permissions**: Higher-level permissions automatically include lower-level permissions

---

## Error Handling

### Common Error Codes

- `400 Bad Request`: Invalid request format or invalid role ID
- `401 Unauthorized`: Missing or invalid JWT token
- `403 Forbidden`: User doesn't have permission to access roles
- `404 Not Found`: Role not found
- `500 Internal Server Error`: Server-side error

### Error Response Format

All errors follow the standard error response format with appropriate HTTP status codes and descriptive messages.

---

## Usage Examples

### Get All Roles with Pagination
```bash
curl -X GET "https://api.example.com/api/v1/roles?page=1&limit=10" \
  -H "Authorization: Bearer <jwt_token>"
```

### Get Specific Role by ID
```bash
curl -X GET "https://api.example.com/api/v1/roles/3" \
  -H "Authorization: Bearer <jwt_token>"
```

### Get Role by Name
```bash
curl -X GET "https://api.example.com/api/v1/roles/name/manager" \
  -H "Authorization: Bearer <jwt_token>"
```

### Get System Roles Only
```bash
curl -X GET "https://api.example.com/api/v1/roles/system" \
  -H "Authorization: Bearer <jwt_token>"
```

---

## Notes

1. **Read-Only API**: The current implementation only provides read operations for roles
2. **System Roles**: The four system roles are pre-seeded in the database during system initialization
3. **Permission Enforcement**: Role permissions are enforced by middleware in other API endpoints
4. **Future Extensions**: Create, Update, and Delete operations may be added for custom roles in future versions