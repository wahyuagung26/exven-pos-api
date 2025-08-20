# Multi-Tenant POS System API

A comprehensive Point of Sale (POS) system built with Go, following Domain-Driven Design (DDD) principles with a modular monolith architecture.

## Architecture Overview

- **Modular Monolith**: Single deployable unit with clear module boundaries
- **Hybrid DDD**: Pragmatic approach with domain-driven design principles
- **Event-Driven**: RabbitMQ for cross-module communication
- **Multi-Tenant**: Complete tenant isolation with subscription management
- **JWT Authentication**: Secure token-based authentication with Redis sessions

## Tech Stack

- **Language**: Go 1.21+
- **HTTP Framework**: Echo v4
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **Message Queue**: RabbitMQ
- **Authentication**: JWT
- **Container**: Docker & Docker Compose

## Project Structure

```
api_pos_system/
├── cmd/
│   ├── api/           # HTTP API server entry point
│   ├── worker/        # Background workers
│   └── migration/     # Database migrations
├── modules/           # Business domains
│   ├── auth/          # Authentication & authorization
│   ├── tenant/        # Multi-tenancy management
│   ├── products/      # Product catalog
│   ├── customers/     # Customer management
│   ├── transactions/  # Sales transactions
│   └── reports/       # Analytics and reporting
├── shared/            # Cross-cutting concerns
│   ├── infrastructure/# Database, cache, messaging
│   ├── middleware/    # HTTP middleware
│   ├── container/     # Dependency injection
│   └── types/         # Shared types and utilities
├── internal/          # Application configuration
│   ├── config/        # Configuration management
│   └── server/        # HTTP server setup
├── docker/            # Docker configurations
└── docs/              # Documentation and SQL schema
```

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)
- Make (optional)

### 1. Clone and Setup

```bash
# Clone the repository
git clone <repository-url>
cd api_pos_system

# Copy environment file
cp .env.example .env

# Edit .env with your configuration
nano .env
```

### 2. Start with Docker Compose

```bash
# Start all services (PostgreSQL, Redis, RabbitMQ, API)
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f api
```

### 3. Database Setup

The database schema will be automatically initialized from `docs/schema.sql` when PostgreSQL starts.

For manual migration:
```bash
docker-compose run --rm migrate
```

### 4. Access the API

The API will be available at: `http://localhost:8080`

Health check endpoint: `http://localhost:8080/health`

## Development

### Local Development

```bash
# Install dependencies
go mod download

# Run database migrations
go run cmd/migration/main.go up

# Run the API server
go run cmd/api/main.go

# Run tests
go test ./...

# Run with hot reload (install air first)
air
```

### Building

```bash
# Build binary
go build -o bin/api cmd/api/main.go

# Build Docker image
docker build -t pos-api .
```

## API Endpoints

### Authentication

- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - Register new tenant and user
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout user
- `POST /api/v1/auth/change-password` - Change password
- `POST /api/v1/auth/reset-password` - Request password reset

### Protected Endpoints (require JWT)

All other endpoints require a valid JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

## Configuration

Key configuration options in `.env`:

```env
# Application
APP_ENV=development
APP_PORT=8080

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=posuser
DB_PASSWORD=pospassword
DB_NAME=posdb

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# RabbitMQ
RABBITMQ_URL=amqp://admin:admin@rabbitmq:5672/

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRY_HOURS=24
```

## Multi-Tenancy

The system supports complete tenant isolation with:
- Separate data for each tenant
- Subscription plans with feature limits
- Data retention policies (14 days for free tier)
- Multiple outlets per tenant
- Role-based access control (Owner, Manager, Cashier)

## Testing

```bash
# Unit tests
go test ./modules/...

# Integration tests
go test ./tests/integration/...

# With coverage
go test -cover ./...

# Generate mocks
go generate ./...
```

## Monitoring

### Service Health

- PostgreSQL: `http://localhost:5432`
- Redis: `http://localhost:6379`
- RabbitMQ Management: `http://localhost:15672` (admin/admin)
- API Health: `http://localhost:8080/health`

### Logs

```bash
# All services
docker-compose logs

# Specific service
docker-compose logs api
docker-compose logs postgres

# Follow logs
docker-compose logs -f api
```

## Troubleshooting

### Database Connection Issues

```bash
# Check PostgreSQL status
docker-compose ps postgres

# Test connection
docker-compose exec postgres psql -U posuser -d posdb
```

### Redis Connection Issues

```bash
# Check Redis status
docker-compose ps redis

# Test connection
docker-compose exec redis redis-cli ping
```

### RabbitMQ Issues

```bash
# Check RabbitMQ status
docker-compose ps rabbitmq

# Access management console
open http://localhost:15672
```

## Production Deployment

### Environment Variables

Ensure all production secrets are properly set:
- Generate strong JWT_SECRET
- Use secure database passwords
- Configure proper CORS origins
- Set appropriate rate limits

### Security Checklist

- [ ] Change default passwords
- [ ] Enable SSL/TLS
- [ ] Configure firewall rules
- [ ] Set up monitoring and alerting
- [ ] Enable audit logging
- [ ] Regular security updates

## License

[Your License Here]

## Support

For issues and questions, please open an issue in the repository.