FROM golang:1.22-alpine AS build
# protoc + plugins installieren - konsolidiert für bessere Layer-Effizienz (Alpine)
RUN apk update && apk upgrade --available && \
    apk add --no-cache \
        protobuf-dev \
        build-base && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

WORKDIR /app
# Modul und main
COPY services/pricing-api/go.mod ./go.mod
RUN go mod download

# Source kopieren und Protobuf generieren
COPY services/pricing-api/. .
COPY protos/ ./protos/

# Protobuf-Dateien generieren
RUN mkdir -p gen/pricing/v1 && \
    protoc --go_out=. --go_opt=module=pricingapi \
           --go-grpc_out=. --go-grpc_opt=module=pricingapi \
           protos/pricing/v1/pricing.proto

# Binary erstellen
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pricing-api

# Runtime
FROM gcr.io/distroless/static:debug-nonroot
COPY --from=build /app/pricing-api /usr/local/bin/pricing-api
EXPOSE 50051
ENTRYPOINT ["/usr/local/bin/pricing-api"]
