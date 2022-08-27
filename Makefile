#!make
include .env

.PHONY: run-consumer run-server run-client
run-consumer:
	go run cmd/consumer/main.go
run-server:
	go run cmd/server/main.go
run-client:
	go run cmd/client/main.go


.PHONY: docker-up docker-down
docker-up:
	docker-compose -f ".docker/docker-compose.yml" --env-file .env up -d
docker-down:
	docker-compose -f ".docker/docker-compose.yml" --env-file .env down


.PHONY: migration migrate migrate-dry
migration:
	goose -dir=${GOOSE_MIGRATION_DIR} create $(name) sql
migrate:
	goose -dir=${GOOSE_MIGRATION_DIR} ${GOOSE_DRIVER} ${GOOSE_DBSTRING} up
migrate-status:
	goose -dir=${GOOSE_MIGRATION_DIR} ${GOOSE_DRIVER} ${GOOSE_DBSTRING} status

.PHONY: .migrate-test
.migrate-qa:
	goose -dir=${QA_GOOSE_MIGRATION_DIR} ${QA_GOOSE_DRIVER} ${QA_GOOSE_DBSTRING} up

.PHONY: .dev-toolscd .protoc
.dev-tools:
	go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
.protoc:
	protoc -I ./api --go_out ./pkg/api/client --go_opt paths=source_relative --go-grpc_out ./pkg/api/client --go-grpc_opt paths=source_relative --grpc-gateway_out ./pkg/api/client --grpc-gateway_opt paths=source_relative ./api/client.proto && \
    protoc -I ./api --go_out ./pkg/api/server --go_opt paths=source_relative --go-grpc_out ./pkg/api/server --go-grpc_opt paths=source_relative --grpc-gateway_out ./pkg/api/server --grpc-gateway_opt paths=source_relative ./api/server.proto


.PHONY: .swagger .buf-generate
.swagger:
	protoc -I ./api --openapiv2_out ./third_party/swagger-ui/api --openapiv2_opt logtostderr=true ./api/client.proto
.buf-generate:
	buf mod update
	buf generate

.PHONY: .test .test-integration
.test:
	$(info Running tests...)
	go test ./...

.test-integration:
	$(info Running tests integration ...)
	go test -tags=integration ./tests -v

.PHONY: .generate-fixture
FILE_GENERATOR_FIXTURE = $(if $(filename),$(filename),"./third_party/generators/fixture/test/structures.go")
.generate-fixture:
	$(info Running generate structures ...)
	go run ./third_party/generators/fixture/generator.go --filename $(FILE_GENERATOR_FIXTURE)

