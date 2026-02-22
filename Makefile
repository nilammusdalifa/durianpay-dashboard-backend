go_bin ?= go
OPENAPI=../openapi.yaml
GEN_PKG=internal/openapigen

.DEFAULT_GOAL := help

help:
	@echo "Available targets:"
	@echo "  make dep         - Install dependencies"
	@echo "  make run         - Run the server locally"
	@echo "  make build       - Build binary into ./bin/server"
	@echo "  make test        - Run unit tests"
	@echo "  make openapi-gen - Regenerate OpenAPI code"
	@echo "  make gen-secret  - Generate a JWT secret"

dep:
	@$(go_bin) mod tidy

run:
	$(go_bin) run main.go

build:
	$(go_bin) build -o bin/server main.go

test:
	$(go_bin) test ./... -v

tool-openapi:
	@$(go_bin) install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

openapi-gen: tool-openapi
	@oapi-codegen -generate "types,chi-server,spec" -package openapigen -o $(GEN_PKG)/openapi.gen.go $(OPENAPI)

gen-secret:
	@$(go_bin) run ./script/gen-secret/main.go
