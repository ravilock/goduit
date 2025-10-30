ifeq (, $(shell command -v docker-compose))
DOCKER_COMPOSE=docker compose
else
DOCKER_COMPOSE=docker-compose
endif

SVC_API := goduit-api
SVC_DB := mongo
SVC_QUEUE := goduit-queue
SVC_REDIS := goduit-redis
SVC_FEED_WORKER := goduit-feed-worker
LOGS_CMD := $(DOCKER_COMPOSE) logs --follow --tail=5

DB_EXEC_CMD := $(DOCKER_COMPOSE) exec mongo bash -c

# Queue configuration (can be overridden: make run QUEUE=redis)
QUEUE ?= rabbitmq

.PHONY: setup
setup:
	echo 'Installing Golang CI Lint'
	# binary will be $(go env GOPATH)/bin/golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s v2.1.6

run: run-all

run-all:
ifeq ($(QUEUE),redis)
	@echo "Starting services with Redis..."
	@QUEUE_TYPE=redis QUEUE_URL=redis://goduit-redis:6379 $(DOCKER_COMPOSE) --profile redis up -d
else
	@echo "Starting services with RabbitMQ..."
	@QUEUE_TYPE=rabbitmq QUEUE_URL=amqp://guest:guest@goduit-queue:5672/ $(DOCKER_COMPOSE) --profile rabbitmq up -d
endif

run-api:
	@$(DOCKER_COMPOSE) up -d $(SVC_API)

run-article-feed-worker:
	@$(DOCKER_COMPOSE) up -d $(SVC_FEED_WORKER)

run-db:
	@$(DOCKER_COMPOSE) up -d $(SVC_DB)

stop: stop-all

stop-all:
	@$(DOCKER_COMPOSE) --profile redis --profile rabbitmq down

logs-api:
	@$(LOGS_CMD) $(SVC_API)

logs-db:
	@$(LOGS_CMD) $(SVC_DB)

logs-worker:
	@$(LOGS_CMD) $(SVC_FEED_WORKER)

logs-queue:
ifeq ($(QUEUE),redis)
	@$(LOGS_CMD) $(SVC_REDIS)
else
	@$(LOGS_CMD) $(SVC_QUEUE)
endif

logs-all:
	@$(LOGS_CMD)

connect-db:
	@$(DB_EXEC_CMD) 'mongosh "mongodb://goduit:goduit-password@mongo:27017/"'

.PHONY: test
test:
	@$(DOCKER_COMPOSE) exec $(SVC_API) go test -count=1 `go list ./... | grep -v integrationTests`

.PHONY: test-verbose
test-verbose:
	@$(DOCKER_COMPOSE) exec $(SVC_API) go test -v -count=1 `go list ./... | grep -v integrationTests`

.PHONY: test-integration
test-integration:
	@$(DOCKER_COMPOSE) exec $(SVC_API) sh -c 'JWT_PRIVATE_KEY_BASE64=$$(base64 -w 0 /app/jwtRS256.key) JWT_PUBLIC_KEY_BASE64=$$(base64 -w 0 /app/jwtRS256.key.pub) go test ./integrationTests/... -count=1 -p 1'

.PHONY: test-integration-verbose
test-integration-verbose:
	@$(DOCKER_COMPOSE) exec $(SVC_API) sh -c 'JWT_PRIVATE_KEY_BASE64=$$(base64 -w 0 /app/jwtRS256.key) JWT_PUBLIC_KEY_BASE64=$$(base64 -w 0 /app/jwtRS256.key.pub) go test ./integrationTests/... -v -count=1 -p 1'

.PHONY: test-both-queues
test-both-queues:
	@echo "=========================================="
	@echo "Testing with RabbitMQ..."
	@echo "=========================================="
	@$(MAKE) stop-all
	@$(MAKE) run QUEUE=rabbitmq
	@sleep 10
	@$(MAKE) test-integration
	@echo ""
	@echo "=========================================="
	@echo "Testing with Redis..."
	@echo "=========================================="
	@$(MAKE) stop-all
	@$(MAKE) run QUEUE=redis
	@sleep 10
	@$(MAKE) test-integration
	@$(MAKE) stop-all
	@echo ""
	@echo "=========================================="
	@echo "✅ All tests passed with both queues!"
	@echo "=========================================="

.PHONY: lint-check
lint-check:
	./bin/golangci-lint run --timeout 5m

bash:
	@$(DOCKER_COMPOSE) exec $(SVC_API) bash

.PHONY: generate-mocks
generate-mocks:
	mockery --all
