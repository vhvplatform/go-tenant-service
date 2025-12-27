# go-tenant-service

> Part of the SaaS Framework - Extracted from monorepo

## Description

[Add service description here]

## Features

- **Multi-Tenant Isolation** - Advanced isolation mechanisms with database isolation, namespace support, and data encryption options
- **Tenant Management** - Complete CRUD operations for tenant lifecycle management
- **Custom Configuration Management** - Flexible per-tenant configuration system with key-value storage
- **Tenant Usage Dashboards** - Real-time usage metrics tracking API calls, storage, and bandwidth consumption
- **Advanced API Security** - Tier-based rate limiting, API key authentication, and tenant context isolation
- **User Management** - Associate users with tenants with role-based access
- **RESTful API** - Clean HTTP/JSON API with comprehensive endpoints
- **gRPC Support** - High-performance gRPC endpoints for service-to-service communication

## Prerequisites

- Go 1.25.5+
- MongoDB 4.4+ (if applicable)
- Redis 6.0+ (if applicable)
- RabbitMQ 3.9+ (if applicable)

## Installation

```bash
# Clone the repository
git clone https://github.com/vhvplatform/go-tenant-service.git
cd go-tenant-service

# Install dependencies
go mod download
```

## Configuration

Copy the example environment file and update with your values:

```bash
cp .env.example .env
```

See [DEPENDENCIES.md](docs/DEPENDENCIES.md) for a complete list of environment variables.

## Development

### Running Locally

```bash
# Run the service
make run

# Or with go run
go run cmd/main.go
```

### Running with Docker

```bash
# Build and run
make docker-build
make docker-run
```

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run benchmarks
make bench

# Run performance tests
make perf-test
```

### Linting

```bash
# Run linters
make lint

# Format code
make fmt

# Run security scan
make security-scan
```

## API Documentation

### Tenant Endpoints

- `POST /api/v1/tenants` - Create a new tenant
- `GET /api/v1/tenants` - List all tenants (paginated)
- `GET /api/v1/tenants/:id` - Get tenant details
- `PUT /api/v1/tenants/:id` - Update tenant
- `DELETE /api/v1/tenants/:id` - Delete tenant (soft delete)

### User Management Endpoints

- `POST /api/v1/tenants/:id/users` - Add user to tenant
- `DELETE /api/v1/tenants/:id/users/:user_id` - Remove user from tenant

### Configuration Endpoints

- `GET /api/v1/tenants/:id/configuration` - Get tenant configuration
- `PUT /api/v1/tenants/:id/configuration` - Update tenant configuration

### Usage Metrics Endpoints

- `GET /api/v1/tenants/:id/metrics` - Get current usage metrics
- `GET /api/v1/tenants/:id/metrics/history` - Get usage history

### Security Features

- **API Key Authentication** - Pass `X-API-Key` header or `Authorization: Bearer <token>`
- **Tenant Context** - Include `X-Tenant-ID` and `X-Tenant-Tier` headers
- **Rate Limiting** - Tier-based limits (Free: 60/min, Basic: 300/min, Professional: 1000/min, Enterprise: 5000/min)

See [docs/API.md](docs/API.md) for detailed API documentation.

## Deployment

See [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md) for deployment instructions.

## Architecture

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for architecture details.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.

## Related Repositories

- [go-shared](https://github.com/vhvplatform/go-shared) - Shared Go libraries

## License

MIT License - see [LICENSE](LICENSE) for details

## Support

- Documentation: [Wiki](https://github.com/vhvplatform/go-tenant-service/wiki)
- Issues: [GitHub Issues](https://github.com/vhvplatform/go-tenant-service/issues)
- Discussions: [GitHub Discussions](https://github.com/vhvplatform/go-tenant-service/discussions)
