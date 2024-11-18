SVC_API := goduit-api
SVC_DB := mongo mongo-express
LOGS_CMD := docker-compose logs --follow --tail=5

run: run-all

run-all: run-db run-api

run-api:
	@docker-compose up -d $(SVC_API)

run-db:
	@docker-compose up -d $(SVC_DB)

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
