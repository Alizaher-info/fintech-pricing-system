FROM golang:1.22.6 AS build

# protoc + plugins installieren - konsolidiert für bessere Layer-Effizienz
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        protobuf-compiler=3.21.12-3 && \
    rm -rf /var/lib/apt/lists/* && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0

WORKDIR /app
# Modul und main
COPY services/pricing-api/go.mod ./go.mod
RUN go mod download

# Protos kopieren und Go-Stubs generieren
COPY protos/ /proto/
RUN mkdir -p gen && \
    protoc -I /proto \
      --go_out=. --go_opt=module=pricingapi \
      --go-grpc_out=. --go-grpc_opt=module=pricingapi \
      /proto/pricing/v1/pricing.proto

# Source kopieren und bauen
COPY services/pricing-api/. .
RUN go build -o /pricing-api

# Runtime
FROM gcr.io/distroless/base-debian12:debug-nonroot
COPY --from=build /pricing-api /pricing-api
EXPOSE 50051
ENTRYPOINT ["/pricing-api"]
