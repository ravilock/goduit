SVC_API := goduit-api
SVC_DB := mongo mongo-express
SVC_QUEUE := goduit-queue
LOGS_CMD := docker-compose logs --follow --tail=5

DB_EXEC_CMD := docker-compose exec mongo bash -c

run: run-all

run-all: run-db run-queue run-api

run-api:
	@docker-compose up -d $(SVC_API)

run-article-feed-worker:
	@docker-compose exec $(SVC_API) go run ./cmd/article-feed-worker/article-feed-worker.go

run-db:
	@docker-compose up -d $(SVC_DB)

run-queue:
	@docker-compose up -d $(SVC_QUEUE)

stop: stop-all

stop-all:
	@docker-compose stop

stop-api:
	@docker-compose stop $(SVC_API)

stop-db:
	@docker-compose stop $(SVC_DB)

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
	@docker-compose exec $(SVC_API) go test -count=1 `go list ./... | grep -v integrationTests`

.PHONY: test-verbose
test-verbose:
	@docker-compose exec $(SVC_API) go test -v -count=1 `go list ./... | grep -v integrationTests`

.PHONY: test-integration
test-integration:
	@docker-compose exec $(SVC_API) go test ./integrationTests/... -count=1 -p 1

.PHONY: test-integration-verbose
test-integration-verbose:
	@docker-compose exec $(SVC_API) go test ./integrationTests/... -v -count=1 -p 1

bash:
	@docker-compose exec $(SVC_API) bash

generate-mocks:
	@docker run -v "$PWD":/src -w /src vektra/mockery --all
