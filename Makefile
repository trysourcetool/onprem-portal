.PHONY: start stop rebuild clean logs status migrate-new migrate lint

DOCKER_COMPOSE = docker compose

start:
	$(DOCKER_COMPOSE) up -d

stop:
	$(DOCKER_COMPOSE) down

rebuild:
	$(DOCKER_COMPOSE) build

clean:
	$(DOCKER_COMPOSE) down -v

logs:
	$(DOCKER_COMPOSE) logs -f

status:
	$(DOCKER_COMPOSE) ps

create-migrate:
	./devtools/create_migrate.sh $(name)

migrate:
	./devtools/cmd/db/main migrate

lint:
	@echo "Running linters on codebase..."
	@gofumpt -l -w . && \
		golangci-lint run --print-issued-lines --fix --go=1.23