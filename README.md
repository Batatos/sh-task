# Skyhawk Security Microservice

A minimal HTTP microservice built with Go, designed to be a foundation for security event processing. This is a simple starting point that can be extended incrementally.

## 🚀 Current Features

### Core Functionality
- **Basic HTTP API**: Simple REST endpoints for security events
- **Health Monitoring**: Health check endpoint for monitoring
- **CORS Support**: Cross-origin resource sharing enabled
- **Graceful Shutdown**: Proper server shutdown handling

### API Endpoints
- `GET /health` - Health check
- `GET /` - Service information
- `GET /api/v1/status` - Service status
- `POST /api/v1/events` - Create security event
- `GET /api/v1/events` - List security events

### Technical Stack
- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Containerization**: Docker + Docker Compose

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Cloud Events  │───▶│  Security API   │───▶│  Detection      │
│   (AWS/GCP/Azure│    │   (Gin/HTTP)    │    │  Engine         │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                       │
                                ▼                       ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Elasticsearch │    │      Redis      │
                       │   (Event Store) │    │   (Cache/PubSub)│
                       └─────────────────┘    └─────────────────┘
                                │                       │
                                ▼                       ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   AWS Services  │    │   Prometheus    │
                       │ (S3/Kinesis/SNS)│    │   + Grafana     │
                       └─────────────────┘    └─────────────────┘
```

## 📋 Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Elasticsearch 8.11+
- Redis 7+
- AWS CLI (for AWS integration)

## 🛠️ Quick Start

### 1. Clone and Setup
```bash
git clone <repository-url>
cd skyhawk-security-microservice
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Start the Service
```bash
# Using Docker Compose (recommended)
docker-compose up --build

# Or run locally
go run cmd/server/main.go
```

### 4. Test the API
```bash
# Make the test script executable
chmod +x scripts/test-api.sh

# Run tests
./scripts/test-api.sh

# Or test manually
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/status
```

### 5. Create a Test Event
```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "login",
    "source": "test",
    "user_id": "user123",
    "ip_address": "192.168.1.100"
  }'
```

## 🚀 Quick Start

### 1. Create a Security Event
```bash
curl -X POST http://localhost:8080/api/v1/security/events \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "login",
    "source": "aws",
    "account_id": "123456789012",
    "region": "us-east-1",
    "user_id": "user123",
    "ip_address": "192.168.1.100",
    "action": "authentication",
    "status": "failed",
    "severity": "medium"
  }'
```

### 2. Create a Detection Rule
```bash
curl -X POST http://localhost:8080/api/v1/security/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "failed_login_detection",
    "description": "Detect multiple failed login attempts",
    "enabled": true,
    "conditions": [
      {
        "field": "event_type",
        "operator": "equals",
        "value": "login"
      },
      {
        "field": "status",
        "operator": "equals",
        "value": "failed"
      }
    ],
    "threshold": 5,
    "time_window": "5m",
    "severity": "high"
  }'
```

### 3. Check for Incidents
```bash
curl http://localhost:8080/api/v1/security/incidents
```

## 📊 API Endpoints

### Security Events
- `POST /api/v1/security/events` - Create a security event
- `GET /api/v1/security/events` - List security events
- `GET /api/v1/security/events/:id` - Get specific event

### Security Incidents
- `GET /api/v1/security/incidents` - List security incidents
- `GET /api/v1/security/incidents/:id` - Get specific incident
- `PUT /api/v1/security/incidents/:id/status` - Update incident status
- `GET /api/v1/security/incidents/statistics` - Get incident statistics

### Detection Rules
- `POST /api/v1/security/rules` - Create detection rule
- `GET /api/v1/security/rules` - List detection rules
- `PUT /api/v1/security/rules/:id` - Update detection rule
- `DELETE /api/v1/security/rules/:id` - Delete detection rule

### Monitoring
- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics

## 🔧 Configuration

### Main Configuration (`config/config.yaml`)
```yaml
server:
  port: 8080
  host: "0.0.0.0"

database:
  redis:
    url: "redis://localhost:6379"
  elasticsearch:
    url: "http://localhost:9200"

aws:
  region: "us-east-1"
  s3_bucket: "security-logs"
  kinesis_stream: "security-events"
  sns_topic: "security-alerts"

detection:
  rules:
    - name: "suspicious_login"
      threshold: 5
      time_window: "5m"
```

## 🧪 Testing

### Run Unit Tests
```bash
go test ./...
```

### Run Integration Tests
```bash
go test ./tests/...
```

### Run with Coverage
```bash
go test -cover ./...
```

## 📈 Monitoring

### Prometheus Metrics
Access metrics at: `http://localhost:9090`

### Grafana Dashboard
Access dashboard at: `http://localhost:3000` (admin/admin)

### Health Checks
```bash
curl http://localhost:8080/health
```

## 🚀 Deployment

### Docker Deployment
```bash
# Build image
docker build -t skyhawk-security-microservice .

# Run container
docker run -p 8080:8080 skyhawk-security-microservice
```

### Kubernetes Deployment
```bash
# Apply Kubernetes manifests
kubectl apply -f k8s/
```

### Production Considerations
- Use proper secrets management
- Configure TLS/SSL certificates
- Set up proper logging and monitoring
- Configure backup strategies
- Implement proper authentication/authorization

## 🔒 Security Features

- **Input Validation**: All inputs are validated and sanitized
- **Rate Limiting**: Built-in rate limiting for API endpoints
- **CORS Support**: Configurable CORS policies
- **Request ID Tracking**: All requests are tracked with unique IDs
- **Secure Headers**: Security headers are automatically added

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Contact the development team
- Check the documentation

## 🔄 Roadmap

- [ ] Support for additional cloud providers (GCP, Azure)
- [ ] Machine learning-based anomaly detection
- [ ] Advanced threat intelligence integration
- [ ] Real-time dashboard improvements
- [ ] API rate limiting and throttling
- [ ] Enhanced logging and audit trails
- [ ] Multi-tenant support
- [ ] Webhook integrations

---

**Built with ❤️ for cloud security**