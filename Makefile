SVC_API := web
SVC_DB := mongo mongo-express

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

test:
	go test ./...
