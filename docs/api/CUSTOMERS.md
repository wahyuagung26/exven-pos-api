# Customers API Documentation

This document provides comprehensive API documentation for the Customers module of ExVen POS Lite system.

## Overview

The Customers API manages customer information within a multi-tenant POS system. Each customer belongs to a specific tenant and includes features like loyalty points tracking, visit history, and transaction statistics. The system supports automatic customer code generation and maintains denormalized data for historical accuracy in transactions.

## Base URL

All customers API endpoints are prefixed with `/api/v1/customers`

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

### 1. Create Customer

Creates a new customer for the authenticated tenant.

**Endpoint:** `POST /api/v1/customers`

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "code": "CUST001",
  "name": "John Doe",
  "email": "john.doe@example.com",
  "phone": "+628123456789",
  "address": "123 Main Street",
  "city": "Jakarta",
  "province": "DKI Jakarta",
  "postal_code": "12345",
  "birth_date": "1990-01-15T00:00:00Z",
  "gender": "male",
  "notes": "VIP customer"
}
```

**Validation Rules:**
- `code`: Optional, 1-50 characters, unique within tenant (auto-generated if not provided)
- `name`: Required, 1-255 characters
- `email`: Optional, valid email format, max 255 characters, unique within tenant
- `phone`: Optional, max 20 characters, unique within tenant
- `address`: Optional
- `city`: Optional, max 100 characters
- `province`: Optional, max 100 characters
- `postal_code`: Optional, max 10 characters
- `birth_date`: Optional, valid date format
- `gender`: Optional, must be 'male' or 'female'
- `notes`: Optional

**Response:**

*Success (201 Created):*
```json
{
  "message": "Customer created successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "code": "CUST001",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+628123456789",
    "address": "123 Main Street",
    "city": "Jakarta",
    "province": "DKI Jakarta",
    "postal_code": "12345",
    "birth_date": "1990-01-15T00:00:00Z",
    "gender": "male",
    "loyalty_points": 0,
    "total_spent": 0.00,
    "visit_count": 0,
    "last_visit_at": null,
    "notes": "VIP customer",
    "is_active": true,
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T10:30:00Z"
  },
  "meta": null
}
```

*Error (409 Conflict):*
```json
{
  "message": "Customer with this email already exists",
  "data": null,
  "errors": {
    "email": ["Email already exists"]
  }
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

### 2. Get All Customers

Retrieves a paginated list of customers for the authenticated tenant with optional filtering and sorting.

**Endpoint:** `GET /api/v1/customers`

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 50, max: 100)
- `name`: Filter by customer name (partial match)
- `code`: Filter by customer code (partial match)
- `email`: Filter by email (partial match)
- `phone`: Filter by phone (partial match)
- `city`: Filter by city (partial match)
- `province`: Filter by province (partial match)
- `gender`: Filter by gender (exact match: male/female)
- `is_active`: Filter by active status (true/false)
- `sort`: Sort field (name, code, email, phone, city, total_spent, visit_count, last_visit_at, created_at)
- `order`: Sort order (asc, desc, default: asc)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Example Request:**
```
GET /api/v1/customers?page=1&limit=20&city=Jakarta&is_active=true&sort=name&order=asc
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Customers retrieved successfully",
  "data": [
    {
      "id": 1,
      "tenant_id": 100,
      "code": "CUST001",
      "name": "John Doe",
      "email": "john.doe@example.com",
      "phone": "+628123456789",
      "address": "123 Main Street",
      "city": "Jakarta",
      "province": "DKI Jakarta",
      "postal_code": "12345",
      "birth_date": "1990-01-15T00:00:00Z",
      "gender": "male",
      "loyalty_points": 150,
      "total_spent": 2750000.00,
      "visit_count": 12,
      "last_visit_at": "2025-08-19T14:30:00Z",
      "notes": "VIP customer",
      "is_active": true,
      "created_at": "2025-08-20T10:30:00Z",
      "updated_at": "2025-08-20T11:30:00Z"
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

### 3. Get Customer by ID

Retrieves a specific customer by its ID.

**Endpoint:** `GET /api/v1/customers/{id}`

**Path Parameters:**
- `id`: Customer ID (integer, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Customer retrieved successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "code": "CUST001",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+628123456789",
    "address": "123 Main Street",
    "city": "Jakarta",
    "province": "DKI Jakarta",
    "postal_code": "12345",
    "birth_date": "1990-01-15T00:00:00Z",
    "gender": "male",
    "loyalty_points": 150,
    "total_spent": 2750000.00,
    "visit_count": 12,
    "last_visit_at": "2025-08-19T14:30:00Z",
    "notes": "VIP customer",
    "is_active": true,
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T11:30:00Z"
  },
  "meta": null
}
```

*Error (404 Not Found):*
```json
{
  "message": "Customer not found",
  "data": null,
  "errors": {}
}
```

*Error (400 Bad Request):*
```json
{
  "message": "Invalid customer ID",
  "data": null,
  "errors": {}
}
```

---

### 4. Update Customer

Updates an existing customer.

**Endpoint:** `PUT /api/v1/customers/{id}`

**Path Parameters:**
- `id`: Customer ID (integer, required)

**Request Headers:**
```
Content-Type: application/json
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "code": "CUST001",
  "name": "John Doe Updated",
  "email": "john.doe.updated@example.com",
  "phone": "+628123456790",
  "address": "456 Main Street",
  "city": "Jakarta",
  "province": "DKI Jakarta",
  "postal_code": "12346",
  "birth_date": "1990-01-15T00:00:00Z",
  "gender": "male",
  "notes": "VIP customer updated",
  "is_active": true
}
```

**Validation Rules:**
- Same as Create Customer request
- `is_active`: Boolean flag to activate/deactivate customer

**Response:**

*Success (200 OK):*
```json
{
  "message": "Customer updated successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "code": "CUST001",
    "name": "John Doe Updated",
    "email": "john.doe.updated@example.com",
    "phone": "+628123456790",
    "address": "456 Main Street",
    "city": "Jakarta",
    "province": "DKI Jakarta",
    "postal_code": "12346",
    "birth_date": "1990-01-15T00:00:00Z",
    "gender": "male",
    "loyalty_points": 150,
    "total_spent": 2750000.00,
    "visit_count": 12,
    "last_visit_at": "2025-08-19T14:30:00Z",
    "notes": "VIP customer updated",
    "is_active": true,
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T12:30:00Z"
  },
  "meta": null
}
```

---

### 5. Delete Customer

Soft deletes a customer by setting `is_active` to false.

**Endpoint:** `DELETE /api/v1/customers/{id}`

**Path Parameters:**
- `id`: Customer ID (integer, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Customer deleted successfully",
  "data": null,
  "meta": null
}
```

*Error (404 Not Found):*
```json
{
  "message": "Customer not found",
  "data": null,
  "errors": {}
}
```

---

### 6. Get Customer by Code

Retrieves a specific customer by its unique code within the tenant.

**Endpoint:** `GET /api/v1/customers/code/{code}`

**Path Parameters:**
- `code`: Customer code (string, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Customer retrieved successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "code": "CUST001",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+628123456789",
    "address": "123 Main Street",
    "city": "Jakarta",
    "province": "DKI Jakarta",
    "postal_code": "12345",
    "birth_date": "1990-01-15T00:00:00Z",
    "gender": "male",
    "loyalty_points": 150,
    "total_spent": 2750000.00,
    "visit_count": 12,
    "last_visit_at": "2025-08-19T14:30:00Z",
    "notes": "VIP customer",
    "is_active": true,
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T11:30:00Z"
  },
  "meta": null
}
```

---

### 7. Get Customer by Phone

Retrieves a specific customer by their phone number.

**Endpoint:** `GET /api/v1/customers/phone/{phone}`

**Path Parameters:**
- `phone`: Customer phone number (string, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Customer retrieved successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "code": "CUST001",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+628123456789",
    "address": "123 Main Street",
    "city": "Jakarta",
    "province": "DKI Jakarta",
    "postal_code": "12345",
    "birth_date": "1990-01-15T00:00:00Z",
    "gender": "male",
    "loyalty_points": 150,
    "total_spent": 2750000.00,
    "visit_count": 12,
    "last_visit_at": "2025-08-19T14:30:00Z",
    "notes": "VIP customer",
    "is_active": true,
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T11:30:00Z"
  },
  "meta": null
}
```

---

### 8. Get Customer by Email

Retrieves a specific customer by their email address.

**Endpoint:** `GET /api/v1/customers/email/{email}`

**Path Parameters:**
- `email`: Customer email address (string, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Customer retrieved successfully",
  "data": {
    "id": 1,
    "tenant_id": 100,
    "code": "CUST001",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+628123456789",
    "address": "123 Main Street",
    "city": "Jakarta",
    "province": "DKI Jakarta",
    "postal_code": "12345",
    "birth_date": "1990-01-15T00:00:00Z",
    "gender": "male",
    "loyalty_points": 150,
    "total_spent": 2750000.00,
    "visit_count": 12,
    "last_visit_at": "2025-08-19T14:30:00Z",
    "notes": "VIP customer",
    "is_active": true,
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T11:30:00Z"
  },
  "meta": null
}
```

---

## Data Models

### Customer Entity

Based on the database schema (`customers` table), the customer entity includes:

```sql
CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    code VARCHAR(50), -- Unique customer code within tenant
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(10),
    birth_date DATE,
    gender gender_type, -- ENUM: 'male', 'female'
    loyalty_points INTEGER DEFAULT 0,
    total_spent DECIMAL(15,2) DEFAULT 0.00,
    visit_count INTEGER DEFAULT 0,
    last_visit_at TIMESTAMP WITH TIME ZONE NULL,
    notes TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    UNIQUE (tenant_id, code)
);
```

### Key Relationships

1. **Tenant**: Each customer belongs to exactly one tenant (multi-tenant isolation)
2. **Transactions**: Customer data is denormalized in transactions for historical accuracy
3. **Loyalty**: System tracks loyalty points, total spent, and visit statistics
4. **Code Generation**: Automatic customer code generation based on name if not provided

---

## Business Rules

1. **Tenant Isolation**: All operations are scoped to the authenticated user's tenant
2. **Unique Code**: Customer code must be unique within a tenant (auto-generated if not provided)
3. **Unique Phone**: Phone number must be unique within a tenant if provided
4. **Unique Email**: Email address must be unique within a tenant if provided
5. **Soft Delete**: Deleting a customer sets `is_active = false` instead of hard deletion
6. **Active Filter**: By default, only active customers (`is_active = true`) are returned
7. **Auto Code Generation**: System generates customer codes using name prefix + timestamp
8. **Stats Tracking**: System maintains denormalized statistics (total_spent, visit_count)
9. **Historical Accuracy**: Customer data is snapshot in transactions for historical reference

---

## Error Handling

### Common Error Codes

- `400 Bad Request`: Invalid request format or validation errors
- `401 Unauthorized`: Missing or invalid JWT token
- `403 Forbidden`: User doesn't have permission to access/modify customer
- `404 Not Found`: Customer not found or doesn't belong to user's tenant
- `409 Conflict`: Duplicate customer code, phone, or email within tenant
- `422 Unprocessable Entity`: Business logic validation failed
- `500 Internal Server Error`: Server-side error

### Validation Errors

The API performs validation on all input fields according to the defined rules and returns detailed field-level validation errors in the standard error response format.

---

## Multi-Tenant Considerations

1. **Automatic Tenant Filtering**: All queries automatically filter by the authenticated user's tenant ID
2. **Cross-Tenant Prevention**: Users cannot access customers from other tenants
3. **Code Uniqueness**: Customer codes are unique within each tenant, not globally
4. **Phone Uniqueness**: Phone numbers are unique within each tenant, not globally
5. **Email Uniqueness**: Email addresses are unique within each tenant, not globally

---

## Customer Statistics

The system automatically maintains customer statistics that are updated through transaction processing:

- **loyalty_points**: Accumulated loyalty points from purchases
- **total_spent**: Total amount spent by the customer across all transactions
- **visit_count**: Number of visits/transactions by the customer
- **last_visit_at**: Timestamp of the customer's last transaction

These statistics are updated automatically by the transaction system and should not be manually modified through the customer API.

---

## Auto-Generated Customer Codes

When creating a customer without providing a `code`, the system automatically generates one using the following pattern:

1. **With Name**: Takes first 4 characters of name (uppercase) + timestamp suffix
   - Example: "John Doe" → "JOHN1234"
2. **Without Name**: Uses "CUST" prefix + timestamp suffix
   - Example: "CUST1234"
3. **Collision Handling**: Adds nanosecond suffix if code already exists
   - Example: "JOHN1234567"

The generated code is guaranteed to be unique within the tenant.