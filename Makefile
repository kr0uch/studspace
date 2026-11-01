CMD_PATH=./cmd/main.go
DOCS_DIR=./docs
SWAG=swag
GO=go

.PHONY: swag
swag:
	$(SWAG) init -g $(CMD_PATH) -o $(DOCS_DIR) --parseInternal --outputTypes go,json,yaml
	$(GO) mod tidy

.PHONY: test
test:
	$(GO) test ./...
