module github.com/aidan-neel/shulker/apps/server

go 1.25.0

require (
	connectrpc.com/connect v1.19.1
	github.com/aidan-neel/shulker/apps/proto v0.0.0
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.8.0
	github.com/joho/godotenv v1.5.1
	github.com/rs/cors v1.11.1
	github.com/zeebo/blake3 v0.2.4
	golang.org/x/crypto v0.48.0
	golang.org/x/net v0.51.0
)

replace github.com/aidan-neel/shulker/apps/proto => ../proto

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.28.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260226221140-a57be14db171 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
