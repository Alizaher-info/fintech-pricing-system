
FROM golang:1.22-alpine AS build

# protoc + plugins installieren - konsolidiert für bessere Layer-Effizienz (Alpine)
RUN apk update && apk upgrade --available && \
    apk add --no-cache \
        protobuf-dev=24.4-r4 \
        build-base=0.5-r3 && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

WORKDIR /app

# Komplettes services Modul kopieren (für shared dependencies)
COPY services/ .
COPY protos/ ./protos/

# Go Modul herunterladen
RUN go mod download

# Protobuf-Dateien generieren (falls noch nicht vorhanden)
RUN mkdir -p pricing-api/gen/pricing/v1 && \
    protoc --proto_path=protos \
           --go_out=. --go_opt=module=services \
           --go-grpc_out=. --go-grpc_opt=module=services \
           pricing/v1/pricing.proto

# Binary aus dem pricing-api Unterverzeichnis erstellen
WORKDIR /app/pricing-api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pricing-api

# Runtime
FROM gcr.io/distroless/static:debug-nonroot
COPY --from=build /app/pricing-api/pricing-api /usr/local/bin/pricing-api
EXPOSE 50051
ENTRYPOINT ["/usr/local/bin/pricing-api"]
