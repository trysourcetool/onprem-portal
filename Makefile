PHONY: create-migrate
create-migrate:
	./devtools/create_migrate.sh $(name)

PHONY: migrate
migrate:
	./devtools/cmd/db/main migrate