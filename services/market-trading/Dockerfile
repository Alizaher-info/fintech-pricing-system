# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git protobuf-dev protoc

# Copy go mod files from market-trading directory
COPY services/market-trading/go.mod services/market-trading/go.sum* ./services/market-trading/
WORKDIR /app/services/market-trading
RUN go mod download

# Install protoc-gen-go and protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# Copy proto files
COPY protos/ /app/protos/

# Copy source code
COPY services/market-trading/ /app/services/market-trading/

# Generate protobuf code
RUN cd /app && \
    protoc --proto_path=protos \
    --go_out=. --go_opt=module=services/market-trading \
    --go-grpc_out=. --go-grpc_opt=module=services/market-trading \
    market/v1/market.proto && \
    echo "=== Proto files generated ===" && \
    echo "=== Looking for pb.go files ===" && \
    find /app -name "*.pb.go" && \
    echo "=== Moving to correct location ===" && \
    rm -rf /app/services/market-trading/gen && \
    mv /app/gen /app/services/market-trading/gen && \
    ls -laR /app/services/market-trading/gen/

# Tidy dependencies and build
RUN go mod tidy && \
    echo "=== go.mod after tidy ===" && \
    cat go.mod
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o /app/market-trading ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/market-trading .

# Note: .env files should be provided at runtime via:
# - Docker environment variables (-e flag)
# - Docker Compose env_file
# - Kubernetes ConfigMaps/Secrets

EXPOSE 50052

CMD ["./market-trading"]
