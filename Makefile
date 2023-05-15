APP_NAME=echo-go-starter

build: ##- Format, Lint, Test, Build.
	@$(MAKE) go-build
	
fmt: ##- Format code.
	@go fmt ./...

go-build: ##- Build binary.
	@go build -o bin/$(APP_NAME)
test: ##- Run tests, output by package, print coverage.
	@go test ./... --race

run: ##- Run the app.
	@./bin/$(APP_NAME) run

seed: ##- Seed the database.
	@./bin/$(APP_NAME) seed

lint: ##- Lint code.
	@golangci-lint run ./...

## Docker-compose commands ( dc = docker-compose )
dc-build: ##- Build docker image.
	@docker-compose build

dc-up: ##- Run docker-compose.
	@docker-compose up -d

dc-down: ##- Stop docker-compose.
	@docker-compose down

dc-logs: ##- Show docker-compose logs.
	@docker-compose logs -f

dc-restart: ##- Restart docker-compose.
	@docker-compose restart

dc-clean: ##- Remove docker-compose containers.
	@docker-compose rm -f

dc-clean-all: ##- Remove docker-compose containers and images.
	@docker-compose rm -f
	@docker rmi -f $(APP_NAME)

dc-shell: ##- Run shell in docker container.
	@docker-compose exec $(APP_NAME) sh

dc-test: ##- Run tests in docker container.
	@docker-compose exec $(APP_NAME) make test

dc-lint: ##- Run lint in docker container.
	@docker-compose exec $(APP_NAME) make lint

dc-build-run: ##- Build and run docker-compose.
	@$(MAKE) docker-build
	@$(MAKE) docker-up

dc-clean-run: ##- Clean and run docker-compose.
	@$(MAKE) docker-clean
	@$(MAKE) docker-up

dc-clean-all-run: ##- Clean and run docker-compose.
	@$(MAKE) docker-clean-all
	@$(MAKE) docker-up

### -----------------------
# --- SQL
### -----------------------

sql-reset: ##- Wizard to drop and create our development database.
	@echo "DROP & CREATE database:"
	@echo "  PGHOST=${PGHOST} PGDATABASE=${PGDATABASE}" PGUSER=${PGUSER}
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	psql -d postgres -c 'DROP DATABASE IF EXISTS "${PGDATABASE}";'
	psql -d postgres -c 'CREATE DATABASE "${PGDATABASE}" WITH OWNER ${PGUSER} TEMPLATE "template0";'

sql-drop-all: ##- Wizard to drop ALL databases: spec, development and tracked by integresql.
	@echo "DROP ALL:"
	TO_DROP=$$(psql -qtz0 -d postgres -c "SELECT 'DROP DATABASE \"' || datname || '\";' FROM pg_database WHERE datistemplate = FALSE AND datname != 'postgres';")
	@echo "$$TO_DROP"
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	@echo "Resetting integresql..."
	curl --fail -X DELETE http://integresql:5000/api/v1/admin/templates
	@echo "Drop databases..."
	echo $$TO_DROP | psql -tz0 -d postgres
	@echo "Done. Please run 'make sql-reset && make sql-spec-reset && make sql-spec-migrate' to reinitialize."

# This step is only required to be executed when the "migrations" folder has changed!
sql: ##- Runs sql format, all sql related checks and finally generates internal/models/*.go.
	@$(MAKE) sql-format
	@$(MAKE) sql-regenerate

sql-regenerate: ##- (opt) Runs sql related checks and finally generates internal/models/*.go.
	@$(MAKE) sql-check-files
	@$(MAKE) sql-spec-reset
	@$(MAKE) sql-spec-migrate
	@$(MAKE) sql-check-and-generate

sql-check-and-generate: sql-check-structure sql-boiler ##- (opt) Runs make sql-check-structure and sql-boiler.

sql-boiler: ##- (opt) Runs sql-boiler introspects the spec db to generate internal/models/*.go.
	@echo "make sql-boiler"
	sqlboiler psql

sql-format: ##- (opt) Formats all *.sql files.
	@echo "make sql-format"
	@find ${PWD} -path "*/tmp/*" -prune -name ".*" -prune -o -type f -iname "*.sql" -print \
		| grep --invert "/app/dumps/" \
		| grep --invert "/app/test/" \
		| xargs -i pg_format --inplace {}

sql-check-files: sql-check-syntax sql-check-migrations-unnecessary-null ##- (opt) Check syntax and unnecessary use of NULL keyword.

# check syntax via the real database
# https://stackoverflow.com/questions/8271606/postgresql-syntax-check-without-running-the-query
sql-check-syntax: ##- (opt) Checks syntax of all *.sql files.
	@echo "make sql-check-syntax"
	@find ${PWD} -path "*/tmp/*" -prune -name ".*" -prune -path ./dumps -prune -false -o -type f -iname "*.sql" -print \
		| grep --invert "/app/dumps/" \
		| grep --invert "/app/test/" \
		| xargs -i sed '1s#^#DO $$SYNTAX_CHECK$$ BEGIN RETURN;#; $$aEND; $$SYNTAX_CHECK$$;' {} \
		| psql -d postgres --quiet -v ON_ERROR_STOP=1

sql-check-migrations-unnecessary-null: ##- (opt) Checks migrations/*.sql for unnecessary use of NULL keywords.
	@echo "make sql-check-migrations-unnecessary-null"
	@(grep -R "NULL" ./migrations/ | grep --invert "DEFAULT NULL" | grep --invert "NOT NULL" | grep --invert "WITH NULL" | grep --invert "NULL, " | grep --invert ", NULL" | grep --invert "RETURN NULL" | grep --invert "SET NULL") \
		&& exit 1 || exit 0

sql-spec-reset: ##- (opt) Drop and creates our spec database.
	@echo "make sql-spec-reset"
	@psql --quiet -d postgres -c 'DROP DATABASE IF EXISTS "${PSQL_DBNAME}";'
	@psql --quiet -d postgres -c 'CREATE DATABASE "${PSQL_DBNAME}" WITH OWNER ${PSQL_USER} TEMPLATE "template0";'

sql-spec-migrate: ##- (opt) Applies migrations/*.sql to our spec database.
	@echo "make sql-spec-migrate"
	@sql-migrate up -env spec | xargs -i echo "[spec DB]" {}

sql-check-structure: sql-check-structure-fk-missing-index sql-check-structure-default-zero-values ##- (opt) Runs make sql-check-structure-*.

sql-check-structure-fk-missing-index: ##- (opt) Ensures spec database objects have FK-indices set.
	@echo "make sql-check-structure-fk-missing-index"
	@cat scripts/sql/fk_missing_index.sql | psql -qtz0 --no-align -d  "${PSQL_DBNAME}" -v ON_ERROR_STOP=1

sql-check-structure-default-zero-values: ##- (opt) Ensures spec database objects default values match go zero values.
	@echo "make sql-check-structure-default-zero-values"
	@cat scripts/sql/default_zero_values.sql | psql -qtz0 --no-align -d "${PSQL_DBNAME}" -v ON_ERROR_STOP=1

dumpfile := /app/dumps/development_$(shell date '+%Y-%m-%d-%H-%M-%S').sql
sql-dump: ##- Dumps the development database to '/app/dumps/development_YYYY-MM-DD-hh-mm-ss.sql'.
	@mkdir -p /app/dumps
	@pg_dump development --format=p --clean --if-exists > $(dumpfile)
	@echo "Dumped '$(dumpfile)'. Use 'cat $(dumpfile) | psql' to restore"

watch-sql: ##- Watches *.sql files in /migrations and runs 'make sql-regenerate' on modifications.
	@echo Watching /migrations. Use Ctrl-c to stop a run or exit.
	watchexec -p -w migrations --exts sql $(MAKE) sql-regenerate

## Help
help: ##- Show this help.
