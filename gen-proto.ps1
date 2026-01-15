$ErrorActionPreference = "Stop"

Write-Host "Creating protobuf output directories..."

# =========================
# Create pb directories
# =========================
$dirs = @(
  "services/auth-service/pkg/pb",
  "services/user-service/pkg/pb",
  "services/asset-service/pkg/pb",
  "api-gateway/pkg/pb" 
)

foreach ($dir in $dirs) {
  if (-not (Test-Path -Path $dir)) {
      New-Item -ItemType Directory -Force -Path $dir | Out-Null
      Write-Host "Created: $dir"
  }
}

# =========================
# AUTH PROTO
# =========================
Write-Host "Generating Auth Proto..."
protoc `
  --proto_path=. `
  --go_out=services/auth-service/pkg/pb --go_opt=paths=source_relative `
  --go-grpc_out=services/auth-service/pkg/pb --go-grpc_opt=paths=source_relative `
  proto/auth/auth.proto

protoc `
  --proto_path=. `
  --go_out=api-gateway/pkg/pb --go_opt=paths=source_relative `
  --go-grpc_out=api-gateway/pkg/pb --go-grpc_opt=paths=source_relative `
  proto/auth/auth.proto

# =========================
# USER PROTO
# =========================
Write-Host "Generating User Proto..."
protoc `
  --proto_path=. `
  --go_out=services/user-service/pkg/pb --go_opt=paths=source_relative `
  --go-grpc_out=services/user-service/pkg/pb --go-grpc_opt=paths=source_relative `
  proto/user/user.proto

protoc `
  --proto_path=. `
  --go_out=services/auth-service/pkg/pb --go_opt=paths=source_relative `
  --go-grpc_out=services/auth-service/pkg/pb --go-grpc_opt=paths=source_relative `
  proto/user/user.proto

protoc `
  --proto_path=. `
  --go_out=api-gateway/pkg/pb --go_opt=paths=source_relative `
  --go-grpc_out=api-gateway/pkg/pb --go-grpc_opt=paths=source_relative `
  proto/user/user.proto

# =========================
# ASSET PROTO
# =========================
Write-Host "Generating Asset Proto..."
protoc `
  --proto_path=. `
  --go_out=services/asset-service/pkg/pb --go_opt=paths=source_relative `
  --go-grpc_out=services/asset-service/pkg/pb --go-grpc_opt=paths=source_relative `
  proto/asset/*.proto

protoc `
  --proto_path=. `
  --go_out=api-gateway/pkg/pb --go_opt=paths=source_relative `
  --go-grpc_out=api-gateway/pkg/pb --go-grpc_opt=paths=source_relative `
  proto/asset/*.proto
Write-Host "Protobuf generated successfully"