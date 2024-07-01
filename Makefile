SVC_API := web
SVC_DB := mongo mongo-express
LOGS_CMD := docker-compose logs --follow --tail=5
MOCKS_COMMAND := $(shell { command -v mockery; } 2>/dev/null)

.PHONY: mocks

setup:
	@docker-compose build

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

bash:
	@docker-compose exec $(SVC_API) sh

mocks:
ifndef MOCKS_COMMAND
	@echo "\nCommand 'mockery' not found!\n"
	@echo "Please, run the following command to install it:"
	@echo "\nMacOSX:"
	@echo "brew install mockery"
	@echo "\nGNU/Linux:"
	@echo "More info, take a look at: https://vektra.github.io/mockery/latest/installation/#github-release-recommended"
	@exit 1
endif
	@rm -rf mocks
	@mockery --all --with-expecter
