# 👥 User Service (Contact Management)

> **High-performance contact management API built with Go and MongoDB**

## 📋 Overview

The User Service is a Go-based microservice that provides comprehensive contact management functionality. Built for high performance and scalability, it offers full CRUD operations for personal contacts with MongoDB persistence and enterprise-grade reliability.

## 🏗️ Architecture

### Technology Stack
- **Language**: Go 1.21+
- **Database**: MongoDB with connection pooling
- **Framework**: Native Go HTTP server with custom routing
- **Container**: Docker with multi-stage build
- **Orchestration**: Kubernetes with StatefulSet database

### Service Details
- **Port**: 5000
- **Health Check**: `/healthz`
- **API Prefix**: `/contacts`
- **Docker Image**: `yaswanthmitta/multiapp-user-management-go`
- **Database**: MongoDB (user-db:27017)

## 🚀 API Documentation

### Base URL
```
http://localhost:5000
```

### Endpoints

#### Create Contact
**POST** `/contacts`

Create a new contact entry.

**Request Body:**
```json
{
  "name": "John Doe",
  "phone": "+1-234-567-8900"
}
```

**Response:**
```json
{
  "message": "Contact created successfully",
  "contact": {
    "id": "507f1f77bcf86cd799439011",
    "name": "John Doe",
    "phone": "+1-234-567-8900"
  }
}
```

#### Get All Contacts
**GET** `/contacts`

Retrieve all contacts.

**Response:**
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "name": "John Doe",
    "phone": "+1-234-567-8900"
  },
  {
    "id": "507f1f77bcf86cd799439012",
    "name": "Jane Smith",
    "phone": "+1-234-567-8901"
  }
]
```

#### Get Contact by ID
**GET** `/contacts/{id}`

Retrieve a specific contact by ID.

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "name": "John Doe",
  "phone": "+1-234-567-8900"
}
```

#### Update Contact
**PUT** `/contacts/{id}`

Update an existing contact.

**Request Body:**
```json
{
  "name": "John Updated",
  "phone": "+1-234-567-9999"
}
```

**Response:**
```json
{
  "message": "Contact updated successfully"
}
```

#### Delete Contact
**DELETE** `/contacts/{id}`

Delete a contact by ID.

**Response:**
```json
{
  "message": "Contact deleted successfully"
}
```

#### Health Check
**GET** `/healthz`

Service health status.

**Response:**
```json
{
  "status": "ok"
}
```

### Error Responses
```json
{
  "error": "Contact not found"
}
```

## 🔧 Configuration

### Environment Variables
```bash
MONGO_URI=mongodb://user-db:27017
PORT=5000
```

### MongoDB Configuration
- **Database**: `contacts_db`
- **Collection**: `contacts`
- **Connection Timeout**: 10 seconds
- **Connection Pooling**: Enabled

## 🏗️ Code Structure

### Main Components

#### Contact Model
```go
type Contact struct {
    ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name  string             `bson:"name" json:"name"`
    Phone string             `bson:"phone" json:"phone"`
}
```

#### Database Connection
```go
func init() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    
    contactsCollection = client.Database("contacts_db").Collection("contacts")
}
```

#### CORS Middleware
```go
func EnableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

## 🐳 Docker

**Multi-stage build**: Go 1.21-alpine for build, alpine for runtime

**Security**: Non-root user, minimal Alpine base image, ca-certificates

**Optimization**: CGO disabled, static binary compilation

**Best practices**: Layer caching, .dockerignore, minimal dependencies

## ☸️ Kubernetes

**Deployment**: Go app with rolling updates and replica management

**StatefulSet**: MongoDB with persistent storage (1Gi PVC)

**Services**: ClusterIP for internal communication

**HPA**: Auto-scaling based on CPU/memory thresholds

**Probes**: Liveness, readiness, and startup health checks

**Resources**: CPU/memory limits and requests configured

## 🔄 CI/CD Pipeline

**CI Tests**: Go linting (go fmt/vet), Code analysis (staticcheck), SCA (dependency check), SAST (gosec), Build validation, Unit tests

**CI Flow**: GitHub runners → Docker build with Git SHA → Push to registry → Update K8s manifests

**CD Flow**: Self-hosted runners → Rolling updates → DAST health checks → Auto rollback

**Image Tagging**: `yaswanthmitta/multiapp-user-management-go:<git-sha>`

## 🔒 Security Implementation

### Application Security
- **Input Validation**: Strict request body validation
- **Error Handling**: Secure error messages without data leakage
- **CORS Configuration**: Controlled cross-origin access
- **Connection Security**: MongoDB connection with authentication

### Container Security
- **Non-root User**: Application runs as non-privileged user
- **Minimal Base Image**: Alpine Linux for reduced attack surface
- **Read-only Filesystem**: Immutable container layers
- **Resource Limits**: CPU and memory constraints

### Database Security
- **Network Isolation**: MongoDB accessible only within cluster
- **Persistent Storage**: Data encryption at rest
- **Access Control**: Service-specific database access
- **Backup Strategy**: Regular data backups

## 📊 Monitoring & Observability

### Metrics Collection
- **Prometheus Integration**: Custom metrics for Go applications
- **Request Metrics**: HTTP request count, duration, status codes
- **Database Metrics**: MongoDB connection pool, query performance
- **System Metrics**: CPU, memory, goroutine count

### Key Performance Indicators
```go
// Custom metrics (example)
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
        },
        []string{"method", "endpoint"},
    )
)
```

### Health Monitoring
- **Liveness Probe**: `/healthz` endpoint for container health
- **Readiness Probe**: Database connectivity check
- **Startup Probe**: Application initialization verification

## 🧪 Testing

### Unit Testing
```bash
# Run Go tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Testing
```bash
# Test with local MongoDB
docker run -d --name test-mongo -p 27017:27017 mongo:7.0

# Set test environment
export MONGO_URI=mongodb://localhost:27017
export PORT=5000

# Run application
go run main.go

# Test endpoints
curl -X POST http://localhost:5000/contacts \
  -H "Content-Type: application/json" \
  -d '{"name": "Test User", "phone": "+1-234-567-8900"}'

curl http://localhost:5000/contacts
```

### Load Testing
```bash
# Using Apache Bench
ab -n 1000 -c 10 http://localhost:5000/contacts

# Using wrk
wrk -t12 -c400 -d30s http://localhost:5000/contacts
```

## 🚨 Troubleshooting

### Common Issues

#### MongoDB Connection Failed
```bash
# Check MongoDB pod status
kubectl get pods -l app=user-db
kubectl logs -l app=user-db

# Check service connectivity
kubectl exec -it <user-service-pod> -- nslookup user-db

# Verify MongoDB is running
kubectl exec -it <mongodb-pod> -- mongo --eval "db.adminCommand('ismaster')"
```

#### High Memory Usage
```bash
# Check Go memory stats
kubectl top pods -l app=user-service

# Check for memory leaks
kubectl exec -it <pod> -- go tool pprof http://localhost:6060/debug/pprof/heap
```

#### Slow Database Queries
```bash
# Check MongoDB performance
kubectl exec -it <mongodb-pod> -- mongo contacts_db --eval "db.contacts.getIndexes()"

# Add database indexes if needed
kubectl exec -it <mongodb-pod> -- mongo contacts_db --eval "db.contacts.createIndex({name: 1})"
```

### Debug Commands
```bash
# Check service status
kubectl get svc user-service

# View application logs
kubectl logs -f deployment/user-service

# Check database connectivity
kubectl exec -it <user-service-pod> -- nc -zv user-db 27017

# Monitor resource usage
kubectl top pods -l app=user-service
```

## 📈 Performance Optimization

### Go Optimization
```go
// Connection pooling
clientOptions := options.Client().
    ApplyURI(mongoURI).
    SetMaxPoolSize(100).
    SetMinPoolSize(10).
    SetMaxConnIdleTime(30 * time.Second)

// Context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

### Database Optimization
```javascript
// MongoDB indexes
db.contacts.createIndex({ "name": 1 })
db.contacts.createIndex({ "phone": 1 })
db.contacts.createIndex({ "_id": 1, "name": 1 })
```

### Resource Configuration
```yaml
resources:
  requests:
    memory: "128Mi"
    cpu: "100m"
  limits:
    memory: "256Mi"
    cpu: "200m"
```

## 🔗 Integration Points

### Frontend Integration
```javascript
// Frontend API calls
const response = await fetch('/api/users/contacts', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name: 'John', phone: '+1234567890' })
});
```

### Service Mesh
- **Nginx Routing**: `/api/users/*` routes to this service
- **Service Discovery**: Kubernetes DNS resolution
- **Load Balancing**: Kubernetes service load balancing
- **Health Checks**: Integrated with Kubernetes probes

## 📚 Dependencies

### Go Modules
```go
module user-service

go 1.21

require (
    go.mongodb.org/mongo-driver v1.12.1
)
```

### External Services
- **MongoDB**: Document database for contact storage
- **Kubernetes**: Container orchestration platform
- **Docker Hub**: Container image registry

## 🏷️ Tags
`golang` `mongodb` `microservices` `rest-api` `kubernetes` `docker` `contacts` `crud`

---

**🌟 High-performance Go microservice with enterprise-grade MongoDB integration!**