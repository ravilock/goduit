SVC_API := web
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

test:
	@docker-compose exec $(SVC_API) go test -p 1 ./... -count=1

test-verbose:
	@docker-compose exec $(SVC_API) go test -p 1 ./... -v -count=1

test-integration:
	@docker-compose exec $(SVC_API) go test ./integrationTests/... -count=1 -p 1

test-integration-verbose:
	@docker-compose exec $(SVC_API) go test ./integrationTests/... -v -count=1 -p 1

bash:
	@docker-compose exec $(SVC_API) bash
