include internal/config/.env
export

PWD = $(shell pwd)
ACCTPATH = $(PWD)/
MPATH = $(ACCTPATH)internal/usecase/repo/postgres/migrations

# Default number of migrations to execute up or down
N = 1

migrate-create:
	@echo "---Creating migration files---"
	migrate create -ext sql -dir $(MPATH) -seq -digits 5 $(NAME)
.PHONY: migrate-create

migrate-up:
	migrate -source file://$(MPATH) -database $(DATABASE_URI) up $(N)
.PHONY: migrate-up

migrate-down:
	migrate -source file://$(MPATH) -database $(DATABASE_URI) down $(N)
.PHONY: migrate-down

migrate-force:
	migrate -source file://$(MPATH) -database $(DATABASE_URI) force $(VERSION)
.PHONY: migrate-force