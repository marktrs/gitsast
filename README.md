# GitSAST

## Background

The aim of this project is to build a simple code scanning application that can detect sensitive keywords in public git repositories. The application will allow users to create, read, update, and delete repositories. Each repository will be identified by a unique name and a link to the repository on GitHub.

Users will be able to trigger a scan against a specific repository in order to detect any potential security issues. The scanning process will involve iterating through the codebase and looking for keywords that indicate the presence of sensitive information. In this case, the application will look for the presence of `public_key` or `private_key` prefixes with allow user to add more rules in further

Once the scan is complete, users will be able to view a Security Scan Result List. This list will show the repositories that have been scanned and the results of each scan. If any sensitive keywords are detected in a repository, this will be indicated in the scan result.

## Architecture

The API server is designed using a layered architecture pattern. The layered architecture pattern separates an application into logical layers that interact with each other to perform a specific task.

The layers of the API server are following:

- Presentation layer: This layer handle the incoming HTTP requests and outgoing HTTP responses. The Go standard library will be used to create an HTTP server that listens to requests and sends responses. The server use a routing package called Bun Router to route requests to the appropriate handler.

- Service layer: This layer handle the business logic of the application. It interact with the data access layer to perform CRUD (Create, Read, Update, Delete) operations on the data.

- Data access layer: This layer will interact with the database to perform CRUD operations. The Bun golang-ORM is used to execute SQL queries against the database.

To ensure the application's scalability and maintainability, we implemented with the following design principles:

- Dependency injection: We use dependency injection to decouple the components of the application and improve testability.

- Separation of concerns: We will ensure that each component of the application is responsible for only one aspect of the application's functionality.

- Single Responsibility Principle (SRP): We will ensure that each component of the application has only one responsibility.

- OpenAPI specification: We use the OpenAPI specification to document the API server's endpoints and responses. This will improve the API server's usability and maintainability.

Tools :

- Database : Postgres
  - ORM : bun golang-ORM
- HTTP Router : bun router
- Task Queue Handler : taskq
- Queue Store : Redis

## API Document

Full API document can be found in [/apidoc](https://github.com/marktrs/gitsast/tree/main/apidoc) directory

## Get Started

### Prerequisites

- Installed [Golang 1.19](https://golang.org/)
- Docker [Docker](https://www.docker.com/)

### Run Test

Unit testing

> make test

or get test coverage profile

> make test-coverage

## Start Postgres, Redis, API server using docker-compose

Build docker image and run using docker compose

> docker compose up -d

Stop running services

> docker-compose down

### Example flow

Add a repository with some files contain keyword `public_key` or `private_key`

```
curl --location 'http://127.0.0.1:8080/api/v1/repository' \
--header 'Content-Type: application/json' \
--data '{
    "name": "wireguard_exporter",
    "remote_url": "https://github.com/mdlayher/wireguard_exporter"
}'
```

Trigger scanner to enqueue analyzer task using newly generated ID

```
curl --location --request POST 'http://127.0.0.1:8080/api/v1/repository/98b57e1c-eb0f-40ea-a690-b7df6a0946e7/scan'
```

Get report status and result using repository ID

```
curl --location 'http://127.0.0.1:8080/api/v1/repository/98b57e1c-eb0f-40ea-a690-b7df6a0946e7/report'
```

## Start API server and db migration with command

### Build GitSAST as an executable file

> make build

### Initialize Postgres Database

Drop table if exist, migrated tables and insert initial data

> make db

### Start API Server

> make start

### CLI Command Reference

```
$ gitsast help

NAME:
   GitSAST - A new cli application

USAGE:
   GitSAST [global options] command [command options] [arguments...]

COMMANDS:
   api      start GitSAST API server
   db       manage database migrations
   help, h  Shows a list of commands or help for one command
```

### Project Layout

```tree
.
├── apidoc
├── app
│   ├── embed
│   │   └── config
│   └── middleware
├── cmd
│   ├── api
│   └── database
├── docker-compose.yml
├── entrypoint.sh
├── internal
│   ├── model
│   ├── queue
│   │   └── task
│   │       └── analyzer
│   ├── recover
│   └── repository
├── scripts
├── testutil
└── workflows

```

A brief description of the layout:

`apidoc` - Contains the importable documentation for the API viewer like OpenAPI, Postman collection.

`app` - Holds the main application code and middleware

`config` - Contains environment specific configuration used by the application

`cmd` - Contains list of CLI command for running the application and other command-line tools.

`internal` - Contains private implementation details of the application that are not intended to be used outside the application itself

`model` - Contains the application's data models and database schema

`queue` - Contains any queue-related code, such as workers and job processing

`recover` - Contains code related to error handling and recovery

`testutil` - Contains utilities used for testing, such as mock objects and test fixtures

`.github` - Contains github workflow files and scripts for automating continuous integration and deployment.

## Future works

- More coverage test % beside the core functionalities
- API Caching
- Distributed queue
- Endpoint for API document
- Speed up scanner for the large git repository

## Notes

- Makefile **MUST NOT** change well-defined command semantics, see Makefile for details.
