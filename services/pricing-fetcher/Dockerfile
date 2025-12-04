FROM golang:1.22-alpine AS build

WORKDIR /app

# Copy the entire services module (for shared dependencies if needed in future)
COPY services/ .

# Download Go modules
RUN go mod download

# Build the pricing-fetcher binary
WORKDIR /app/pricing-fetcher
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pricing-fetcher ./cmd/server

# Runtime
FROM gcr.io/distroless/static:debug-nonroot
COPY --from=build /app/pricing-fetcher/pricing-fetcher /usr/local/bin/pricing-fetcher
ENTRYPOINT ["/usr/local/bin/pricing-fetcher"]
