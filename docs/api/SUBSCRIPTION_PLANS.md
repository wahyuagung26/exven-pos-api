# Subscription Plans API Documentation

This document provides comprehensive API documentation for the Subscription Plans module of ExVen POS Lite system.

## Overview

The Subscription Plans API manages subscription plans that define the features and limitations available to tenants in the multi-tenant POS system. Subscription plans control access to features like the number of outlets, users, products, and transactions per month.

## Base URL

All subscription plans API endpoints are prefixed with `/api/v1/subscription-plans`

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

### 1. Get All Subscription Plans

Retrieves a paginated list of all available subscription plans with optional filtering and sorting.

**Endpoint:** `GET /api/v1/subscription-plans`

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 50, max: 100)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Example Request:**
```
GET /api/v1/subscription-plans?page=1&limit=20
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Subscription plans retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Free",
      "description": "Paket gratis dengan batasan fitur dan retensi data 14 hari",
      "price": 0.00,
      "max_outlets": 1,
      "max_users": 2,
      "max_products": null,
      "max_transactions_per_month": null,
      "features": [
        "basic_pos",
        "basic_reports",
        "data_retention_14_days"
      ],
      "is_active": true,
      "created_at": "2025-08-20T10:30:00Z",
      "updated_at": "2025-08-20T10:30:00Z"
    },
    {
      "id": 2,
      "name": "Starter",
      "description": "Paket untuk bisnis kecil",
      "price": 99000.00,
      "max_outlets": 2,
      "max_users": 5,
      "max_products": null,
      "max_transactions_per_month": null,
      "features": [
        "full_pos",
        "advanced_reports",
        "customer_management",
        "data_retention_unlimited"
      ],
      "is_active": true,
      "created_at": "2025-08-20T10:30:00Z",
      "updated_at": "2025-08-20T10:30:00Z"
    },
    {
      "id": 3,
      "name": "Business",
      "description": "Paket untuk bisnis menengah",
      "price": 299000.00,
      "max_outlets": 5,
      "max_users": 15,
      "max_products": null,
      "max_transactions_per_month": null,
      "features": [
        "full_pos",
        "advanced_reports",
        "customer_management",
        "inventory_management",
        "multi_payment",
        "data_retention_unlimited"
      ],
      "is_active": true,
      "created_at": "2025-08-20T10:30:00Z",
      "updated_at": "2025-08-20T10:30:00Z"
    },
    {
      "id": 4,
      "name": "Enterprise",
      "description": "Paket untuk bisnis besar",
      "price": 599000.00,
      "max_outlets": 999,
      "max_users": 999,
      "max_products": null,
      "max_transactions_per_month": null,
      "features": [
        "full_pos",
        "advanced_reports",
        "customer_management",
        "inventory_management",
        "multi_payment",
        "api_access",
        "custom_integration",
        "data_retention_unlimited"
      ],
      "is_active": true,
      "created_at": "2025-08-20T10:30:00Z",
      "updated_at": "2025-08-20T10:30:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 4
  }
}
```

---

### 2. Get Subscription Plan by ID

Retrieves a specific subscription plan by its ID.

**Endpoint:** `GET /api/v1/subscription-plans/{id}`

**Path Parameters:**
- `id`: Subscription plan ID (integer, required)

**Request Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**

*Success (200 OK):*
```json
{
  "message": "Subscription plan retrieved successfully",
  "data": {
    "id": 2,
    "name": "Starter",
    "description": "Paket untuk bisnis kecil",
    "price": 99000.00,
    "max_outlets": 2,
    "max_users": 5,
    "max_products": null,
    "max_transactions_per_month": null,
    "features": [
      "full_pos",
      "advanced_reports",
      "customer_management",
      "data_retention_unlimited"
    ],
    "is_active": true,
    "created_at": "2025-08-20T10:30:00Z",
    "updated_at": "2025-08-20T10:30:00Z"
  },
  "meta": null
}
```

*Error (404 Not Found):*
```json
{
  "message": "Subscription plan not found",
  "data": null,
  "errors": {}
}
```

*Error (400 Bad Request):*
```json
{
  "message": "Invalid subscription plan ID",
  "data": null,
  "errors": {}
}
```

---

## Data Models

### Subscription Plan Entity

Based on the database schema (`subscription_plans` table), the subscription plan entity includes:

```sql
CREATE TABLE subscription_plans (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    max_outlets INTEGER NOT NULL DEFAULT 1,
    max_users INTEGER NOT NULL DEFAULT 1,
    max_products INTEGER DEFAULT NULL, -- NULL = unlimited
    max_transactions_per_month INTEGER DEFAULT NULL,
    features JSONB, -- Array of available features
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Plan Features

The `features` field contains an array of feature strings that determine what functionality is available to tenants with this subscription plan:

#### Available Features:
- `basic_pos`: Basic point of sale functionality
- `full_pos`: Full point of sale functionality with advanced features
- `basic_reports`: Basic reporting and analytics
- `advanced_reports`: Advanced reporting with detailed analytics
- `customer_management`: Customer relationship management features
- `inventory_management`: Advanced inventory tracking and management
- `multi_payment`: Support for multiple payment methods
- `api_access`: Access to API endpoints for integrations
- `custom_integration`: Support for custom integrations
- `data_retention_14_days`: Data retention limited to 14 days (for free tier)
- `data_retention_unlimited`: Unlimited data retention

### Plan Limitations

Each subscription plan defines limits for:

1. **max_outlets**: Maximum number of outlets/stores the tenant can create
2. **max_users**: Maximum number of users the tenant can have
3. **max_products**: Maximum number of products (NULL = unlimited)
4. **max_transactions_per_month**: Maximum transactions per month (NULL = unlimited)

### Default Plans

The system includes 4 default subscription plans:

1. **Free Plan** - Free tier with basic features and 14-day data retention
2. **Starter Plan** - Entry-level paid plan for small businesses
3. **Business Plan** - Mid-tier plan for growing businesses with advanced features
4. **Enterprise Plan** - Full-featured plan for large businesses with unlimited resources

---

## Business Rules

1. **Public Access**: Subscription plans are publicly readable to allow potential customers to view available options
2. **Read-Only**: The API currently provides read-only access - plan creation and updates are handled by system administrators
3. **Active Plans**: Only active plans (`is_active = true`) are returned by default
4. **Feature-Based Access**: Plan features determine what functionality is available to tenants
5. **Limit Enforcement**: Plan limits are enforced during tenant operations (outlet creation, user creation, etc.)
6. **Currency**: All prices are in Indonesian Rupiah (IDR)

---

## Error Handling

### Common Error Codes

- `400 Bad Request`: Invalid request format or invalid plan ID
- `401 Unauthorized`: Missing or invalid JWT token
- `404 Not Found`: Subscription plan not found
- `500 Internal Server Error`: Server-side error

### Error Response Format

All error responses follow the standard error format with appropriate HTTP status codes and error messages.

---

## Usage Examples

### Getting All Plans for Plan Selection
```bash
curl -X GET "https://api.example.com/api/v1/subscription-plans" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Getting Specific Plan Details
```bash
curl -X GET "https://api.example.com/api/v1/subscription-plans/2" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Pagination Example
```bash
curl -X GET "https://api.example.com/api/v1/subscription-plans?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Integration Notes

1. **Tenant Subscription**: This API is used by the tenant management system to display available plans during registration
2. **Plan Enforcement**: Plan limits are checked by other modules when creating resources (outlets, users, products)
3. **Feature Checks**: Application features are enabled/disabled based on the tenant's current subscription plan features
4. **Billing Integration**: Plan pricing information is used for billing and payment processing