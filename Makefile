ifeq (, $(shell command -v docker-compose))
DOCKER_COMPOSE=docker compose
else
DOCKER_COMPOSE=docker-compose
endif

SVC_API := goduit-api
SVC_DB := mongo
SVC_QUEUE := goduit-queue
SVC_FEED_WORKER := goduit-feed-worker
LOGS_CMD := $(DOCKER_COMPOSE) logs --follow --tail=5

DB_EXEC_CMD := $(DOCKER_COMPOSE) exec mongo bash -c

run: run-all

run-all: run-db run-queue run-api run-article-feed-worker

run-api:
	@$(DOCKER_COMPOSE) up -d $(SVC_API)

run-article-feed-worker:
	@$(DOCKER_COMPOSE) up -d $(SVC_FEED_WORKER)

run-db:
	@$(DOCKER_COMPOSE) up -d $(SVC_DB)

run-queue:
	@$(DOCKER_COMPOSE) up -d $(SVC_QUEUE)

stop: stop-all

stop-all:
	@$(DOCKER_COMPOSE) stop

stop-api:
	@$(DOCKER_COMPOSE) stop $(SVC_API)

stop-db:
	@$(DOCKER_COMPOSE) stop $(SVC_DB)

logs-api:
	@$(LOGS_CMD) $(SVC_API)

logs-db:
	@$(LOGS_CMD) $(SVC_DB)

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
	@$(DOCKER_COMPOSE) exec $(SVC_API) go test ./integrationTests/... -count=1 -p 1

.PHONY: test-integration-verbose
test-integration-verbose:
	@$(DOCKER_COMPOSE) exec $(SVC_API) go test ./integrationTests/... -v -count=1 -p 1

bash:
	@$(DOCKER_COMPOSE) exec $(SVC_API) bash

.PHONY: generate-mocks
generate-mocks:
	mockery --all
