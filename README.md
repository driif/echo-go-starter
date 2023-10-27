ECHO-GO-STARTER
===============

This is a starter project for Go. It includes a Makefile for building and testing, and a Dockerfile for building a Docker image.

This Server is built using the [Echo](https://echo.labstack.com/) framework.

##Requirements
- Go 1.21
- Docker
- Docker-compose
- Postgres 15

# for SQL Boiler
- `go install github.com/volatiletech/sqlboiler/v4`
- `go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql`

# for Golang Lint
- `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2`

# for Migrations 
- `go install github.com/rubenv/sql-migrate/sql-migrate`

# for OAPI Codegen
- `go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest`


## Makefile
Makefile contains some useful commands to build and run the project, some of them are:

- `make dc-up`               - Start the docker containers
- `make build`               - Build the go binary into `./bin`
- `make seed`                - Seed the database with some data
- `make run`                 - Run the go binary from `./bin`

## OAPI Codegen
OAPI Codegen is used to generate the API models and handlers from the `./api/paths/oapi_api.yaml` file. The generated files are placed in `./internal/types/oapi_api/`. The generated files should not be modified manually.

Commands:
- `make gen-oapi OAPI_SERVICE=oapi_api` - Generate the API models and handlers from the `./api/paths/oapi_api.yaml` file.
Instead of oapi_api you should use oapi_(service name) for example oapi_user. 
There is a possibility to create your command for generating service in Makefile:
```
gen-oapi-user: OAPI_SERVICE=oapi_user
gen-oapi-user: gen-oapi
```

`.env` is used to set the environment variables, should be placed to the root `./` directory.