# Claude Code Agent Prompt: Multi-Tenant POS System

Create a complete Golang project skeleton for a multi-tenant Point of Sale (POS) system with the following requirements:

## Project Structure & Architecture

**Follow Hybrid DDD approach with modular monolith pattern:**
- Use pragmatic DDD with clear module boundaries
- Organize by business domains with shared infrastructure
- Single binary deployment with dependency injection
- Event-driven communication between modules

## Technical Stack Requirements

1. **Database**: PostgreSQL with GORM
2. **Cache**: Redis for session management
3. **Authentication**: JWT tokens
4. **Message Queue**: RabbitMQ for simple pub/sub (no event sourcing)
5. **RPC**: JSON-RPC for service communication
6. **HTTP Framework**: Echo
7. **DI Container**: Custom dependency injection container
8. **Testing**: Unit tests per module with mocking support

## Core Business Modules

Based on the database schema, create these main modules:
1. **auth** - Authentication, authorization, sessions
2. **tenant** - Multi-tenancy, subscriptions, plans
3. **products** - Product catalog, categories, inventory
4. **customers** - Customer management
5. **transactions** - Sales transactions, payments
6. **reports** - Analytics and reporting

## Project Structure

```
pos-system/
├── cmd/
│   ├── api/           # HTTP API server with DI setup
│   ├── worker/        # Background workers
│   └── migration/     # Database migrations
├── modules/           # Business domains
│   ├── auth/
│   ├── tenant/
│   ├── products/
│   ├── customers/
│   ├── transactions/
│   └── reports/
│       ├── domain/        # Entities, DTOs
│       ├── persistence/   # DB operations
│       ├── services/      # Handle business logic
│       ├── handlers/      # HTTP/RPC handlers
│       └── module.go      # Module registration & DI setup
├── shared/            # Cross-cutting concerns
│   ├── infrastructure/
│   │   ├── database/      # DB connections & Schema Models
│   │   ├── cache/         # Redis client
│   │   └── messaging/     # RabbitMQ pub/sub
│   ├── middleware/
│   │   ├── auth.go        # JWT validation
│   │   ├── tenant.go      # Tenant isolation
│   │   ├── ratelimit.go   # Rate limiting
│   │   └── logging.go     # Request logging
│   ├── types/
│   │   ├── errors/        # Custom error types
│   │   ├── events/        # Event definitions
│   │   └── pagination/    # Common pagination
│   ├── utils/
│   │   ├── crypto/        # Encryption utilities
│   │   ├── validator/     # Input validation
│   │   └── jsonrpc/       # RPC utilities
│   └── container/         # DI container implementation
├── internal/          # Application platform
│   ├── config/        # Configuration management
│   ├── server/        # HTTP server setup
│   ├── worker/        # Background worker setup
│   └── app/           # Application bootstrap
├── mocks/             # Generated mocks for testing
├── scripts/           # Build and deployment scripts
├── .env.example       # Environment variables
└── docker/            # Docker configurations
```

## Module Structure Template

Each business module should follow this clean, simplified structure:

```
modules/[module_name]/
├── domain/
│   ├── entities.go        # Domain entities (pure business objects)
│   ├── dtos.go           # Data Transfer Objects (Request/Response)
│   └── interfaces.go     # Repository & Service contracts
├── persistence/
│   ├── models.go         # Query-optimized models for this module
│   ├── repository.go     # Repository implementation
│   └── mappers.go        # Domain ↔ Model conversion
├── services/
│   ├── [entity]_service.go  # Core business logic
│   └── event_service.go     # Event handling
├── handlers/
│   ├── http.go           # REST endpoint handlers
│   ├── rpc.go            # JSON-RPC handlers
│   └── events.go         # Event subscribers
└── module.go             # Module registration with DI container
```

## Key Features to Implement

### 1. Multi-Tenancy
- Tenant isolation at database level using middleware
- Subscription plan enforcement
- Data retention policies (14-day limit for free tier)

### 2. Authentication & Authorization
- JWT-based authentication with Redis sessions
- Role-based access control (Owner, Manager, Cashier)
- Multi-outlet user assignment

### 3. Core POS Features
- Product management with multi-outlet inventory
- Customer management with denormalized transaction data
- Sales transactions with historical accuracy (snapshot fields)
- Multiple payment methods support
- Real-time stock updates via events

### 4. Dependency Injection System
- Custom DI container for service registration
- Module-based service registration
- Interface-based dependencies for testing

### 5. Event-Driven Architecture
- Simple pub/sub using RabbitMQ
- Cross-module communication via events
- Event handlers for business logic

## Database Schema Implementation

Implement the provided multi-tenant POS database schema with these key tables and denormalization:

### Core Tables:
- `tenants`, `subscription_plans`, `tenant_subscriptions`
- `users`, `roles`, `outlets`, `user_outlets`
- `products`, `product_categories`, `product_stocks`
- `customers`

### Transaction Tables (with denormalized snapshot fields):
```sql
sales_transactions:
- customer_name_snapshot, customer_phone_snapshot, customer_email_snapshot
- cashier_name_snapshot, outlet_name_snapshot, outlet_code_snapshot

sales_transaction_items:
- product_name_snapshot, product_sku_snapshot, product_category_snapshot
- product_unit_snapshot, cost_price_snapshot
```

### Archive Tables:
- `archived_sales_transactions`, `archived_sales_transaction_items`
- `data_retention_logs`

## Implementation Guidelines

### 1. Layer Responsibilities

#### Domain Layer
- Pure business entities without external dependencies
- Data Transfer Objects for API communication
- Interface definitions for repositories and services
- Business validation rules

#### Persistence Layer
- Query-optimized models for specific operations
- Repository implementations with database operations
- Database-specific query optimizations
- Data mapping between domain entities and persistence models

#### Services Layer
- Business logic orchestration
- Cross-entity operations and validations
- Event publishing and handling
- Integration with external services

#### Handlers Layer
- HTTP/REST endpoint handling
- JSON-RPC endpoint handling
- Event subscription and processing
- Request/response transformation

### 2. Data Flow Pattern

```
HTTP Request → Handler → Service → Repository → Database
                ↓
        Domain Entity ↔ Persistence Model ↔ Schema Model
                ↓
              DTO Response
```

### 3. Two-Model Architecture

#### Schema Models (Infrastructure)
- Located in `shared/infrastructure/database/schema_*.go`
- Pure table structure definition for migrations
- Focus on DDL concerns (CREATE, ALTER, DROP)
- Define relationships and constraints
- Used only for database schema management

#### Persistence Models (Module-specific)
- Located in `modules/[domain]/persistence/models.go`
- Query-optimized for specific operations
- Custom structs for complex queries and joins
- Performance-focused field selection
- Used for actual CRUD operations

```go
// Infrastructure schema model (DDL only)
// shared/infrastructure/database/schema_products.go
type Product struct {
    ID               uint64 `gorm:"primaryKey"`
    TenantID         uint64 `gorm:"index;not null"`
    SKU              string `gorm:"size:100;uniqueIndex:idx_tenant_sku"`
    Name             string `gorm:"size:255;not null"`
    // Complete schema definition...
}

// Persistence model (Query operations)
// modules/products/persistence/models.go
type ProductModel struct {
    ID           uint64
    TenantID     uint64
    SKU          string
    Name         string
    CategoryName string `gorm:"column:category_name"` // Join field
    StockQty     int    `gorm:"column:stock_qty"`     // Aggregated
}
```

### 1. Dependency Injection Container
```go
// shared/container/container.go
type Container interface {
    Register(name string, factory interface{})
    Get(name string) interface{}
    Resolve(target interface{}) error
}
```

### 2. Module Registration Pattern
```go
// modules/products/module.go
func (m *Module) Register() {
    m.container.Register("product.repository", func() domain.ProductRepository {
        return infrastructure.NewProductRepository(db)
    })
    m.container.Register("product.useCase", func() application.CreateProductUseCase {
        return application.NewCreateProductUseCase(repo, eventBus)
    })
}
```

### 3. Event Bus Implementation
```go
// shared/infrastructure/messaging/eventbus.go
type EventBus interface {
    Publish(ctx context.Context, event Event) error
    Subscribe(topic string, handler EventHandler) error
}
```

### 4. Inter-Module Communication
- Use dependency injection to inject other module services
- Communicate via events for loose coupling
- Use interfaces for module dependencies

### 5. Testing Strategy
- Generate mocks using mockery
- Unit tests per module with mocked dependencies
- Focus on business logic testing
- Integration tests for cross-module scenarios

## Configuration Management

```go
type Config struct {
    Database DatabaseConfig
    Redis    RedisConfig
    RabbitMQ RabbitMQConfig
    JWT      JWTConfig
    Server   ServerConfig
}
```

## Background Workers

Implement workers for:
- Data retention (archiving free tier transactions)
- Stock synchronization across outlets
- Report generation
- Notification sending

## API Design

### REST Endpoints Structure:
```
/api/v1/auth/*          # Authentication endpoints
/api/v1/tenants/*       # Tenant management
/api/v1/products/*      # Product catalog
/api/v1/customers/*     # Customer management
/api/v1/transactions/*  # Sales transactions
/api/v1/reports/*       # Analytics and reports
```

### Middleware Pipeline:
1. CORS and security headers
2. Request logging
3. Rate limiting
4. JWT authentication
5. Tenant isolation
6. Input validation

## Single Binary Deployment

### 1. Main Application Bootstrap
```go
// cmd/api/main.go
func main() {
    cfg := config.Load()
    container := di.NewContainer()
    
    registerSharedServices(container, cfg)
    
    // Register all modules
    auth.NewModule(container).Register()
    products.NewModule(container).Register()
    transactions.NewModule(container).Register()
    
    server := server.New(container)
    server.Start(cfg.Port)
}
```

### 2. Docker Configuration
- Single container deployment
- Multi-stage build for optimization
- Docker Compose for local development

## Code Quality Standards

1. **Error Handling**: Custom error types with proper context
2. **Logging**: Structured logging with tenant context
3. **Testing**: Unit tests with >80% coverage per module
4. **Documentation**: Comprehensive Go documentation
5. **Linting**: golangci-lint configuration
6. **Security**: Input validation, SQL injection prevention
7. **Mock Generation**: Use //go:generate mockery for interfaces

## Files to Generate

### Core Files:
- `go.mod` with all required dependencies
- Database migration files using schema models
- Schema models in `shared/infrastructure/database/schema_*.go`
- Persistence models in each module's persistence layer
- DI container implementation
- Event bus implementation
- Configuration management
- All module structures with sample implementations

### Infrastructure:
- Database connection and transaction manager
- Schema models for migration and DDL operations
- Redis client for session management
- RabbitMQ client for pub/sub
- Middleware implementations
- HTTP server setup with Echo

### Testing:
- Mock interfaces for all repositories and services
- Sample unit tests for each module
- Test configuration and utilities

### Deployment:
- Dockerfile for single binary
- Docker Compose for local development
- Environment configuration templates
- Build scripts

## Important Implementation Notes

1. **Tenant Isolation**: Every database query must include tenant_id filtering
2. **Historical Accuracy**: Always populate snapshot fields during transaction creation
3. **Module Independence**: Modules should not import each other directly
4. **Event-Driven**: Use events for cross-module communication
5. **Single Binary**: All modules compile into one executable
6. **DI Pattern**: Use dependency injection for all service dependencies
7. **Mock Testing**: Generate mocks for all interfaces for unit testing

## Sample Module Implementation

Include working examples of:
- Complete auth module with JWT and sessions
- Products module with inventory management and two-model architecture (schema + persistence)
- Basic transaction processing with denormalized data
- Event publishing and subscribing between modules
- Database operations with schema models for migrations and persistence models for queries
- HTTP handlers with proper error handling using Echo
- Unit tests with mocks for each layer
- Clear separation between domain entities, schema models, and persistence models

Start with a minimal working version that can be run locally with Docker Compose, then expand with additional features. Include comprehensive README with setup and usage instructions.

### Key Implementation Notes

1. **Clean Layer Separation**: Domain entities should be pure business objects, while persistence models handle database concerns
2. **Repository Pattern**: Each domain should have repository interfaces in domain layer and implementations in persistence layer
3. **Service Layer**: Contains business logic and orchestrates between repositories and external services
4. **Handler Layer**: Focuses only on HTTP/RPC request handling and response formatting
5. **Testing**: Each layer should be independently testable with proper mocking
6. **Module Independence**: Modules should communicate only through well-defined interfaces and events