CMD_PATH=./cmd/main.go
DOCS_DIR=./docs
SWAG=swag
GO=go
DOCKER=docker
DOCKER_COMPOSE=$(DOCKER) compose
DOCKER_SYSTEM=$(DOCKER) system

.PHONY: swag
swag:
	$(SWAG) init -g $(CMD_PATH) -o $(DOCS_DIR) --parseInternal --outputTypes go,json,yaml
	$(GO) mod tidy

.PHONY: test
test:
	$(GO) test ./...

.PHONY: docker-build
docker-build:
	$(DOCKER_COMPOSE) build

.PHONY: docker-up
docker-up:
	$(DOCKER_COMPOSE) up -d

.PHONY: docker-down
docker-down:
	$(DOCKER_COMPOSE) down

.PHONY: docker-restart
docker-restart: docker-down docker-up

.PHONY: docker-logs
docker-logs:
	$(DOCKER_COMPOSE) logs -f

.PHONY: docker-clean
docker-clean:
	$(DOCKER_COMPOSE) down -v
	$(DOCKER_SYSTEM) prune -f
