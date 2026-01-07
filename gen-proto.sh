#!/bin/bash

protoc --proto_path=proto/auth \
  --go_out=services/auth-service/pkg/pb --go_opt=paths=source_relative \
  --go-grpc_out=services/auth-service/pkg/pb --go-grpc_opt=paths=source_relative \
  proto/auth/auth.proto

protoc --proto_path=proto/user \
  --go_out=services/user-service/pkg/pb --go_opt=paths=source_relative \
  --go-grpc_out=services/user-service/pkg/pb --go-grpc_opt=paths=source_relative \
  proto/user/user.proto

protoc --proto_path=proto/user \
  --go_out=services/auth-service/pkg/pb --go_opt=paths=source_relative \
  --go-grpc_out=services/auth-service/pkg/pb --go-grpc_opt=paths=source_relative \
  proto/user/user.proto