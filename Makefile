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

## Help
help: ##- Show this help.
