PHONY: create-migrate
create-migrate:
	./devtools/create_migrate.sh $(name)

PHONY: migrate
migrate:
	./devtools/cmd/db/main migrate

PHONY: lint
lint:
	@echo "Running linters on codebase..."
	@gofumpt -l -w . && \
		golangci-lint run --print-issued-lines --fix --go=1.23