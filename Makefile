.PHONY: run-bot run-server

run-bot:
	go run cmd/bot/main.go

run-server:
	go run cmd/server/main.go

LOCAL_BIN:=$(CURDIR)/bin
.PHONY: .dev-tools
.dev-tools:
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: .dev-buf
.dev-buf:
	buf generate