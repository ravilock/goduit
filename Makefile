SVC_API := web
SVC_CONSUMER := feed-producer
SVC_DB := mongo mongo-express
SVC_QUEUE := queue
LOGS_CMD := docker-compose logs --follow --tail=5

run: run-all

run-all:
	@docker-compose up -d

run-api:
	@docker-compose up -d $(SVC_API)

run-db:
	@docker-compose up -d $(SVC_DB)

run-consumer:
	@docker-compose up -d $(SVC_CONSUMER)

run-queue:
	@docker-compose up -d $(SVC_QUEUE)

stop: stop-all

stop-all:
	docker-compose stop

stop-api:
	@docker-compose stop $(SVC_API)

stop-db:
	@docker-compose stop $(SVC_DB)

stop-consumer:
	@docker-compose stop $(SVC_CONSUMER)

stop-queue:
	@docker-compose stop $(SVC_QUEUE)

logs-api:
	@$(LOGS_CMD) $(SVC_API) $(SVC_CONSUMER)

logs-db:
	@$(LOGS_CMD) $(SVC_DB)

logs-all:
	@$(LOGS_CMD)

test:
	@docker-compose exec $(SVC_API) go test -p 1 ./... -count=1

test-verbose:
	@docker-compose exec $(SVC_API) go test -p 1 ./... -v -count=1

bash:
	@docker-compose exec $(SVC_API) sh
