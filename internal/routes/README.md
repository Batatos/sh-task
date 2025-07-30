# Routes Package

This package contains all route definitions for the Skyhawk Security Microservice.

## Structure

```
internal/routes/
├── routes.go      # Main route setup and configuration
└── routes_test.go # Route testing
```

## Best Practices

### 1. Route Organization
- **Group related endpoints** under meaningful prefixes (e.g., `/api/v1/events`)
- **Use versioning** for API endpoints (`/api/v1/`, `/api/v2/`)
- **Keep route definitions separate** from business logic

### 2. Middleware Application
- **Global middleware** is applied at the router level
- **Route-specific middleware** can be applied to groups
- **Handler-specific middleware** should be in the handler package

### 3. Adding New Routes

To add new route groups:

1. **Add the handler** in `internal/handler/`
2. **Update the Handler struct** in `internal/handler/handler.go`
3. **Add route definitions** in `internal/routes/routes.go`

Example:
```go
// In internal/handler/handler.go
type Handler struct {
    HealthHandler *HealthHandler
    EventHandler  *EventHandler
    UserHandler   *UserHandler  // New handler
}

// In internal/routes/routes.go
users := apiV1.Group("/users")
{
    users.POST("/", handlers.UserHandler.CreateUser)
    users.GET("/:id", handlers.UserHandler.GetUser)
}
```

### 4. Testing Routes
- **Unit test route setup** in `routes_test.go`
- **Integration test endpoints** in handler tests
- **Use httptest** for HTTP testing

## Current Routes

### Health Endpoints
- `GET /health` - Health check
- `GET /` - Root endpoint
- `GET /api/v1/status` - API status

### Event Endpoints
- `POST /api/v1/events/` - Create event
- `GET /api/v1/events/` - List all events
- `GET /api/v1/events/:id` - Get specific event
- `PUT /api/v1/events/:id` - Update event
- `DELETE /api/v1/events/:id` - Delete event 