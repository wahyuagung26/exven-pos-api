# Authentication API Documentation

This document describes all authentication and authorization endpoints for the ExVen POS Lite API.

## Base URL

All authentication endpoints are prefixed with `/api/v1/auth`

## Standard Response Format

All responses follow the standard format defined in the architecture:

### Success Response
```json
{
  "message": "Success message",
  "data": {}, 
  "meta": null
}
```

### Error Response
```json
{
  "message": "Error message",
  "data": null,
  "errors": {
    "field_name": ["Error description"]
  }
}
```

## Authentication Endpoints

### 1. Login

Authenticate user and get access tokens.

- **URL**: `POST /api/v1/auth/login`
- **Authentication**: Not required

#### Request Body
```json
{
  "email": "string",     // required, valid email
  "password": "string"   // required
}
```

#### Request Example
```json
{
  "email": "admin@company.com",
  "password": "password123"
}
```

#### Success Response (200 OK)
```json
{
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600,
    "user": {
      "id": 1,
      "tenant_id": 1,
      "email": "admin@company.com",
      "full_name": "Admin User",
      "phone": "+6281234567890",
      "is_active": true,
      "role": {
        "id": 1,
        "name": "tenant_owner",
        "display_name": "Pemilik Bisnis",
        "permissions": ["tenant.*"]
      }
    }
  },
  "meta": null
}
```

#### Error Response (401 Unauthorized)
```json
{
  "message": "Invalid credentials",
  "data": null,
  "errors": {}
}
```

#### Error Response (400 Bad Request - Validation)
```json
{
  "message": "Validation failed",
  "data": null,
  "errors": {
    "email": ["Email is required"],
    "password": ["Password is required"]
  }
}
```

---

### 2. Register

Create new tenant and user account.

- **URL**: `POST /api/v1/auth/register`
- **Authentication**: Not required

#### Request Body
```json
{
  "tenant_name": "string",    // required, min: 3, max: 255
  "business_type": "string",  // optional
  "email": "string",          // required, valid email
  "phone": "string",          // optional
  "password": "string",       // required, min: 8
  "full_name": "string"       // required, max: 255
}
```

#### Request Example
```json
{
  "tenant_name": "My Coffee Shop",
  "business_type": "Food & Beverage",
  "email": "owner@mycoffeeshop.com",
  "phone": "+6281234567890",
  "password": "password123",
  "full_name": "John Doe"
}
```

#### Success Response (201 Created)
```json
{
  "message": "Registration successful. Please verify your email.",
  "data": {
    "tenant_id": 1,
    "user": {
      "id": 1,
      "tenant_id": 1,
      "email": "owner@mycoffeeshop.com",
      "full_name": "John Doe",
      "phone": "+6281234567890",
      "is_active": true
    },
    "message": "Registration successful. Please verify your email."
  },
  "meta": null
}
```

#### Error Response (400 Bad Request - Validation)
```json
{
  "message": "Validation failed",
  "data": null,
  "errors": {
    "email": ["Invalid email format"],
    "password": ["Password must be at least 8 characters"],
    "tenant_name": ["Tenant name is required"]
  }
}
```

#### Error Response (400 Bad Request - Business Logic)
```json
{
  "message": "Email already exists",
  "data": null,
  "errors": {}
}
```

---

### 3. Refresh Token

Generate new access token using refresh token.

- **URL**: `POST /api/v1/auth/refresh`
- **Authentication**: Not required

#### Request Body
```json
{
  "refresh_token": "string"  // required
}
```

#### Request Example
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Success Response (200 OK)
```json
{
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600
  },
  "meta": null
}
```

#### Error Response (401 Unauthorized)
```json
{
  "message": "Invalid or expired refresh token",
  "data": null,
  "errors": {}
}
```

---

### 4. Logout

Invalidate current session and tokens.

- **URL**: `POST /api/v1/auth/logout`
- **Authentication**: Required (Bearer Token)

#### Request Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Success Response (200 OK)
```json
{
  "message": "Logout successful",
  "data": null,
  "meta": null
}
```

#### Error Response (500 Internal Server Error)
```json
{
  "message": "Failed to logout",
  "data": null,
  "errors": {}
}
```

---

### 5. Change Password

Change user password (requires current password).

- **URL**: `POST /api/v1/auth/change-password`
- **Authentication**: Required (Bearer Token)

#### Request Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Request Body
```json
{
  "old_password": "string",  // required
  "new_password": "string"   // required, min: 8
}
```

#### Request Example
```json
{
  "old_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

#### Success Response (200 OK)
```json
{
  "message": "Password changed successfully",
  "data": null,
  "meta": null
}
```

#### Error Response (400 Bad Request - Wrong Password)
```json
{
  "message": "Current password is incorrect",
  "data": null,
  "errors": {}
}
```

#### Error Response (400 Bad Request - Validation)
```json
{
  "message": "Validation failed",
  "data": null,
  "errors": {
    "old_password": ["Old password is required"],
    "new_password": ["New password must be at least 8 characters"]
  }
}
```

---

### 6. Reset Password

Send password reset instructions to email.

- **URL**: `POST /api/v1/auth/reset-password`
- **Authentication**: Not required

#### Request Body
```json
{
  "email": "string"  // required, valid email
}
```

#### Request Example
```json
{
  "email": "user@company.com"
}
```

#### Success Response (200 OK)
```json
{
  "message": "Password reset instructions sent to your email",
  "data": null,
  "meta": null
}
```

#### Error Response (400 Bad Request - Email Not Found)
```json
{
  "message": "Email address not found",
  "data": null,
  "errors": {}
}
```

#### Error Response (400 Bad Request - Validation)
```json
{
  "message": "Validation failed",
  "data": null,
  "errors": {
    "email": ["Invalid email format"]
  }
}
```

---

### 7. Verify Email

Verify email address using verification token.

- **URL**: `POST /api/v1/auth/verify-email`
- **Authentication**: Not required

#### Request Body
```json
{
  "token": "string"  // required
}
```

#### Request Example
```json
{
  "token": "abc123def456ghi789jkl012mno345pqr678stu901vwx234yz"
}
```

#### Success Response (200 OK)
```json
{
  "message": "Email verified successfully",
  "data": null,
  "meta": null
}
```

#### Error Response (400 Bad Request - Invalid Token)
```json
{
  "message": "Invalid or expired verification token",
  "data": null,
  "errors": {}
}
```

#### Error Response (400 Bad Request - Validation)
```json
{
  "message": "Validation failed",
  "data": null,
  "errors": {
    "token": ["Token is required"]
  }
}
```

---

## Data Models

### User Response Model
```typescript
interface UserResponse {
  id: number;
  tenant_id: number;
  email: string;
  full_name: string;
  phone?: string;
  is_active: boolean;
  role: RoleResponse;
}
```

### Role Response Model
```typescript
interface RoleResponse {
  id: number;
  name: string;
  display_name: string;
  permissions: string[];
}
```

### Token Pair Model
```typescript
interface TokenPair {
  access_token: string;
  refresh_token: string;
  expires_in: number;  // seconds until access_token expires
}
```

---

## Authentication Flow

### Standard Authentication Flow
1. **Registration**: User registers with tenant information
2. **Email Verification**: User verifies email address (optional but recommended)
3. **Login**: User authenticates and receives access/refresh tokens
4. **API Access**: User includes access token in Authorization header for protected endpoints
5. **Token Refresh**: When access token expires, use refresh token to get new tokens
6. **Logout**: User invalidates session and tokens

### Token Usage
- **Access Token**: Include in Authorization header as `Bearer {token}`
- **Refresh Token**: Used only for `/auth/refresh` endpoint
- **Token Expiry**: Access tokens expire in 1 hour, refresh tokens in 30 days

### Multi-Tenant Considerations
- Each user belongs to a single tenant
- Tenant isolation is enforced at the database level
- User can only access data within their tenant

---

## Error Codes

| HTTP Status | Code | Description |
|-------------|------|-------------|
| 400 | Bad Request | Invalid request format or validation errors |
| 401 | Unauthorized | Invalid credentials or expired/invalid tokens |
| 404 | Not Found | Resource not found |
| 500 | Internal Server Error | Server-side errors |

---

## Security Notes

1. **Password Requirements**: Minimum 8 characters
2. **Token Security**: 
   - Access tokens expire in 1 hour
   - Refresh tokens expire in 30 days
   - Tokens are invalidated on logout
3. **Rate Limiting**: Authentication endpoints are rate-limited
4. **HTTPS Only**: All authentication endpoints must use HTTPS in production
5. **Password Hashing**: Passwords are hashed using bcrypt
6. **Session Management**: User sessions are stored in Redis for fast invalidation

---

## Testing Examples

### cURL Examples

#### Login
```bash
curl -X POST https://api.example.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@company.com",
    "password": "password123"
  }'
```

#### Register
```bash
curl -X POST https://api.example.com/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_name": "My Coffee Shop",
    "business_type": "Food & Beverage",
    "email": "owner@mycoffeeshop.com",
    "phone": "+6281234567890",
    "password": "password123",
    "full_name": "John Doe"
  }'
```

#### Access Protected Endpoint
```bash
curl -X GET https://api.example.com/api/v1/protected-endpoint \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```