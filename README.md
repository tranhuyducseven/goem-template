# Golang with Echo and Mongo Template

This is a template for a Golang project using the Echo framework and MongoDB. It includes a basic project structure, configuration, and Docker setup.

## Table of Contents

1. [Overview](#overview)
2. [Installation](#installation)
3. [Usage](#usage)
4. [Configuration](#configuration)
5. [Testing](#testing)
6. [Deployment](#deployment)
7. [Contributing](#contributing)
8. [License](#license)

## Overview

Structure of the project:

```
┣📦 .github                      # Github actions
┃
┣📦 cmd                          # Command line
┃ ┗ 📂 api
┃    ┗ 📂 server                 # Setup routes, middlewares, services,...
┃    ┣ 📂 test                   # Integration tests
┃    ┗ 📜 main.go                # Api entry point
┣ 📦 config                        # Configuration
┃    ┣ 📜 env.go                 # Environment setup
┃    ┗ 📜 logger.go              # Logger setup: zerolog, openobserve, stdout
┣ 📦 constants                    # Constants: errors, mail, ...
┣ 📦 docs                         # Documents
┣ 📦 internal                     # Internal packages
┃    ┣ 📂 health                 # Health module
┃    ┃   ┣ 📜 health.go          # Health check handler
┃    ┃   ┗ 📜 route.go           # Health check route
┃    ┣ 📂 user
┃    ┃    ┣ 📂 services
┃    ┃    ┃   ┣ 📜 user.go       # User service
┃    ┃    ┃   ┗ 📜 user_test.go  # User service test
┃    ┃    ┗ 📂 handler
┃    ┃        ┣ 📜 user.go       # User handler
┃    ┃        ┗ 📜 user_test.go  # User handler test
┃    ┗ 📂 middlleware            # Middlewares
┣ 📦 pkg                          # Public packages
┃    ┣ 📂 db
┃    ┃   ┣ 📂 models             # Database models
┃    ┃   ┗ 📜 init.go            # Database connection
┃    ┣ 📂 openobserve            # Openobserve
┃    ┗ 📂 utils                  # Utilities
┃
┣ 📜 .air.toml                    # Air configuration
┣ 📜 .env.example                 # Env example
┣ 📜 .gitignore                   # Git ignore
┣ 📜 app.compose.yml              # App docker compose
┣ 📜 compose.yml                  # Docker compose
┣ 📜 Dockerfile                    # Dockerfile
┣ 📜 go.mod                       # Go modules
┣ 📜 go.sum                       # Go modules
┣ 📜 Makefile                      # Makefile
┗ 📜 README.md                    # Readme
```

## Installation

1. Clone the repository
2. Install dependencies
- `go mod tidy`
- Install Air for hot reload: `go install github.com/cosmtrek/air@latest`
3. Set up configuration (if any)
- Copy `.env.example` to `.env` and update the values

## Usage
- Run bootstrapping: `make compose` to start the services
    1. This will start the Openobserve service that runs on port 5080
        - You can access the Openobserve dashboard at `http://localhost:5080`
        - Login with the default email and password in `compose.yml`, or update the values
        - Access the [Ingestion API - Trace](http://localhost:5080/web/ingestion/custom/traces/opentelemetry) 
        - Copy the `Authorization` header value and update the `OPENOBSERVE_CREDENTIAL` in the `.env` file
        - Access the [Trace Tab](http://localhost:5080/web/traces) to view the traces
    2. This will start the MongoDB service that runs on port 27017
        - Access the MongoDB connection string and update the `MONGO_URI` in the `.env` file
- Run the application: `make run`
- Run the application with hot reload: `make watch`
- Run application with Docker: `make start`
- Stop the docker container: `make stop`
- Shutdown and clean up: `make shutdown`
- Restart the application: `make restart`

## Configuration

- Update the configuration in the `.env` file
- Update the logger configuration in `config/logger.go`
- Update the environment setup in `config/env.go`
- Update the database connection in `pkg/db/init.go`
- Update the database models in `pkg/db/models`
- Add new routes in `cmd/api/server/routes.go`
- Add new services in `internal/<module>/services`
- Add new handlers in `internal/<module>/handlers`
- Add new middlewares in `internal/middleware`
- Add new constants in `constants`
- If you want to add new packages, add them in `pkg` such as `pkg/<package>`: jwt, cookie, mail, cache, oauth, s3, ...

## Testing
- Run tests: `make test`
- Run tests with coverage: `make coverage`

## Deployment

- With existing Dockerfile, you can deploy the application to any cloud provider
- Update the Dockerfile if needed

## Contributing
- Fork the repository
- Create a new branch
- Make your changes
- Create a pull request

## License