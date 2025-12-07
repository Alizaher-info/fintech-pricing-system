module services

go 1.22

require google.golang.org/grpc v1.65.0 // gRPC framework

require google.golang.org/protobuf v1.34.1 // Protocol Buffers support

require github.com/golang-jwt/jwt/v5 v5.1.0 // JWT authentication

require github.com/lib/pq v1.10.9 // to connect to Postgres DB

require github.com/joho/godotenv v1.5.1 // to load .env files

require github.com/segmentio/kafka-go v0.4.47 // Kafka client library

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
)
