.PHONY: run-bot run-server run-client

run-bot:
	go run cmd/bot/main.go

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go

docker-up:
	docker-compose -f ".docker/docker-compose.yml" --env-file .env up -d

docker-down:
	docker-compose -f ".docker/docker-compose.yml" --env-file .env down


.PHONY: .dev-toolscd
.dev-tools:
	go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: .protoc
.protoc:
	protoc -I ./api --go_out ./pkg/api --go_opt paths=source_relative --go-grpc_out ./pkg/api --go-grpc_opt paths=source_relative --grpc-gateway_out ./pkg/api --grpc-gateway_opt paths=source_relative ./api/api.proto

.PHONY: .swagger
.swagger:
	protoc -I ./api --openapiv2_out ./swagger/api --openapiv2_opt logtostderr=true ./api/api.proto

.PHONY: .buf-generate
.buf-generate:
	buf mod update
	buf generate