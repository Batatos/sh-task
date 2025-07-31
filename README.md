# 🛡️ Security Microservice



## 🚀 Features

- **Go + PostgreSQL**: Clean, efficient backend with proper database integration
- **RESTful API**: Complete CRUD operations for security events
- **Docker & Docker Compose**: Containerized development and deployment
- **Clean Architecture**: Well-organized code structure with separation of concerns
- **Health Monitoring**: Built-in health checks and status endpoints
- **Middleware**: CORS, Request ID tracking, and error recovery

## 🏗️ Architecture

```
├── cmd/server/           # Application entry point
├── internal/
│   ├── database/         # Database connection management
│   ├── handler/          # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/           # Data models and structs
│   ├── repository/       # Database operations layer
│   ├── routes/           # Route definitions
│   └── server/           # HTTP server setup
├── database/             # Database schema and migrations
└── scripts/              # Utility scripts
```

## 🛠️ Tech Stack

- **Language**: Go 1.21
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL 15
- **Containerization**: Docker & Docker Compose
- **Testing**: Testify
- **Architecture**: Clean Architecture with Repository Pattern

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)

### Run with Docker
```bash
# Clone the repository
git clone https://github.com/Batatos/sh-task.git
cd sh-task

# Start the services
docker-compose up --build

# The API will be available at http://localhost:8080
```

### API Endpoints

#### Health & Status
- `GET /health` - Health check
- `GET /` - Root endpoint
- `GET /api/v1/status` - API status

#### Security Events (CRUD)
- `POST /api/v1/events/` - Create security event
- `GET /api/v1/events/` - List all events
- `GET /api/v1/events/:id` - Get specific event
- `PUT /api/v1/events/:id` - Update event
- `DELETE /api/v1/events/:id` - Delete event

### Example Usage

```bash
# Create a security event
curl -X POST http://localhost:8080/api/v1/events/ \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "login",
    "severity": "high",
    "source": "web-application",
    "description": "Multiple failed login attempts",
    "event_data": {
      "ip": "192.168.1.100",
      "user": "admin",
      "attempts": 5
    }
  }'

# Get all events
curl http://localhost:8080/api/v1/events/
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# Test API endpoints
./scripts/test-api.sh
```

## 📊 Database Schema

The service uses PostgreSQL with the following key tables:

- **security_events**: Stores security events with JSONB for flexible data
- **Indexes**: Optimized for common queries (event_type, severity, created_at)
- **Triggers**: Automatic updated_at timestamp management

## 🔧 Development

### Local Development
```bash
# Install dependencies
go mod download

# Run locally (requires PostgreSQL)
go run cmd/server/main.go
```

### Adding New Features
1. **Add models** in `internal/models/`
2. **Create repository** in `internal/repository/`
3. **Add handlers** in `internal/handler/`
4. **Define routes** in `internal/routes/routes.go`

## 🎯 Key Design Decisions

### **Clean Architecture**
- Separation of concerns between layers
- Dependency injection for testability
- Repository pattern for data access

### **Go Best Practices**
- Proper error handling
- Context usage for timeouts
- Graceful shutdown
- Structured logging

### **Database Design**
- UUID primary keys for scalability
- JSONB for flexible event data
- Proper indexing strategy
- Database constraints for data integrity

## 🔒 Security Features

- **Input Validation**: Request binding and validation
- **SQL Injection Protection**: Parameterized queries
- **CORS Configuration**: Proper cross-origin handling
- **Request Tracing**: Request ID middleware for debugging

## 📈 Scalability Considerations

- **Stateless Design**: Easy horizontal scaling
- **Database Connection Pooling**: Efficient resource usage
- **Containerization**: Consistent deployment across environments
- **Health Checks**: Monitoring and orchestration ready