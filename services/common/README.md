# Go gRPC Services Template System

This directory contains reusable Docker templates for Go gRPC services in the portfolio-finance project.

## 📁 Files

- `Dockerfile.go-grpc` - Generic Dockerfile template for Go gRPC services
- `docker-compose.examples.yml` - Example configurations for future services

## 🚀 Usage

### For Pricing API (Already configured)

```bash
cd deploy/compose
docker-compose up pricing-api
```

### For new Go gRPC services

1. **Create your Go service directory:**
   ```
   services/your-service/
   ├── main.go
   ├── go.mod
   └── go.sum
   ```

2. **Add to docker-compose.yml:**
   ```yaml
   your-service:
     build: 
       context: ../..
       dockerfile: services/common/Dockerfile.go-grpc
       args:
         SERVICE_NAME: your-service
         SERVICE_PORT: 50054  # Choose unique port
         PROTO_FILE: your-proto/v1/service.proto
     ports:
       - "50054:50054"
   ```

3. **Required Build Arguments:**
   - `SERVICE_NAME`: Name of the service directory
   - `SERVICE_PORT`: Port to expose (default: 50051)
   - `PROTO_FILE`: Path to protobuf file (relative to /proto/)

## 🔧 Template Features

- **Multi-stage build** for minimal production images
- **Protobuf compilation** at build time
- **Go dependency caching** for faster builds
- **Distroless runtime** for security
- **Configurable ports** and proto files

## 📋 Port Assignments

| Service       | Port  | Status |
|---------------|-------|---------|
| pricing-api   | 50051 | ✅ Active |
| event-gateway | 50052 | 📋 Planned |
| offers-worker | 50053 | 📋 Planned |

## 🛠️ Development

To test the template locally:
```bash
docker build --build-arg SERVICE_NAME=pricing-api --build-arg PROTO_FILE=pricing/v1/pricing.proto -f services/common/Dockerfile.go-grpc .
```

## 🎯 Benefits

- **Consistency**: All Go services use the same build pattern
- **Maintainability**: One template to update, all services benefit
- **Performance**: Optimized for Docker layer caching
- **Security**: Minimal attack surface with distroless images
- **Portfolio-ready**: Enterprise-grade patterns