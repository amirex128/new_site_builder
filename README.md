# Go Backend Service

## Overview

This project is a backend service, built with Go, designed to power a robust and scalable platform. It appears to be tailored for e-commerce or a similar system requiring management of products, orders, users, and content.

Key features suggested by the project structure include:
- Product Catalog Management
- Order Processing and Management
- Customer and User Authentication/Management
- Content Management (Articles/Pages)
- Payment Processing
- Notification System (Email/SMS - inferred from `cmd/consumer/notification/`)
- Discount and Coupon Management
- Site and Website Configuration

The service leverages a modern technology stack:
- **Go**: For the core backend logic.
- **Docker**: For containerization and ease of deployment (as seen in `docker/Dockerfile` and `docker/docker-compose.yml`).
- **gRPC**: For inter-service communication (indicated by `protobuf/` and the gRPC provider in `bootstrap/service_provider/`).
- **REST APIs**: For client-server communication (evident from `internal/api/router/` and HTTP handlers).
- **Databases & Brokers**:
    - MySQL (indicated by `docker/mysql/init.sql` and `bootstrap/service_provider/mysql_provider.go`)
    - MongoDB (seen with `docker/mongodb/` and `bootstrap/service_provider/mongo_provider.go`)
    - Kafka (inferred from `bootstrap/service_provider/kafka_provider.go`)
    - RabbitMQ (inferred from `bootstrap/service_provider/rabbit_provider.go`)
    - Elasticsearch (inferred from `bootstrap/service_provider/elastic_provider.go`)
    - Redis (inferred from `bootstrap/service_provider/redis_provider.go`)

## Project Structure

The project follows a structured layout to organize its components:

-   `bootstrap/`: Contains the application initialization logic, including setting up consumers, dependency injection containers, and service providers.
-   `cmd/`: Holds the main entry points for the application. This includes:
    -   `cli/`: Command-line interface tools.
    -   `server/`: The main server application (e.g., HTTP/gRPC server).
    -   `consumer/`: Background workers or message queue consumers.
-   `config/`: Manages application configuration, including loading configuration files (e.g., `config/files/config-local.json`).
-   `docker/`: Includes `Dockerfile` for building service images and `docker-compose.yml` for orchestrating multi-container setups.
-   `internal/`: Houses the core business logic and private code of the application, not intended for import by other projects.
    -   `api/`: Defines API routes (HTTP/gRPC), request/response handlers, and middleware.
    -   `application/`: Contains the application-specific business logic, use cases (services), and Data Transfer Objects (DTOs).
    -   `contract/`: Defines interfaces (contracts) for components like repositories, services, and other abstractions, promoting loose coupling.
    -   `domain/`: Represents the core business entities, enums, value objects, and domain-specific logic.
    -   `infra/`: Provides concrete implementations of the interfaces defined in `contract/`. This includes database repositories, clients for external services, etc.
-   `pkg/`: Contains shared utility packages that can be used across different parts of the project or potentially in other projects (e.g., service discovery utilities).
-   `protobuf/`: Stores Protocol Buffer (`.proto`) definition files for defining gRPC services and message types.
-   `vendor/`: Manages project dependencies, typically populated by Go modules.

## Configuration

Application configuration is managed through JSON files located in the `config/files/` directory, such as `config-local.json` and `config-stage.json`. The application likely loads a specific configuration file based on an environment variable (e.g., `APP_ENV`, `GO_ENV`, or a similar custom variable) that determines which `config-*.json` file to use.

The `config/config.go` file defines a struct that holds configuration parameters, which are populated from environment variables. These environment variables may be sourced from the loaded JSON configuration files or set directly in the deployment environment. The `bootstrap/service_provider/config_provider.go` file is likely involved in the process of loading these configuration files and making them available to the application.

For local development, you can typically copy an existing configuration file (e.g., `config-local.json`), rename it if necessary (e.g., to `config-custom.json`), and modify its values to match your local setup. Ensure the appropriate environment variable is set to make the application load your custom file, or that your environment variables directly match those expected by `config/config.go`.

Key configuration areas (primarily managed via environment variables defined in `config/config.go`):
- Application settings (e.g., `APP_LOG_LEVEL`, `APP_PORT`)
- Database connections (MongoDB, MySQL)
- Message broker details (RabbitMQ; Kafka provider also exists)
- Caching services (Redis)
- Search services (Elasticsearch provider exists)
- Storage service credentials (S3-like object storage)
- JWT settings (secret token, issuer, audience)
- Service discovery endpoints (Consul, etcd, Zookeeper - inferred from `pkg/service_discovery/` and providers in `bootstrap/service_provider/`)

## Building and Running the Application

### Using Docker (Recommended)

The easiest way to get the application and its dependencies up and running is by using Docker and Docker Compose. The project includes a `docker-compose.yml` file in the `docker/` directory.

To build images and start all services (application, databases, message brokers, etc.):
```bash
# Navigate to the docker directory
cd docker/

# Run docker-compose (use 'docker compose' for newer Docker versions)
docker-compose up -d --build --remove-orphans
```
Alternatively, the `Makefile` (in the project root) provides a convenient target that starts specific services for development:
```bash
make up
```
This command typically handles starting services like Redis, MongoDB, MySQL, MinIO, and RabbitMQ, along with the main application, in detached mode. Check the `Makefile` for the exact services started by this command.

### Running Locally with Go

To run the application directly using Go:

**Prerequisites:**
-   **Go**: Version `1.24.2` or newer (as specified in `go.mod`).
-   **External Services**: If you are not using Docker, you'll need to manually set up and configure all external services (e.g., MySQL, MongoDB, Redis, RabbitMQ). Ensure their connection details match your application configuration.

**Running the Server:**
1.  Ensure your configuration is correctly set up (e.g., via environment variables matching those in `config/config.go`, or by ensuring your chosen `config-*.json` file is loaded).
2.  The project uses `air` for live reloading during development (see `Makefile`'s `init` target for installation). If you have `air` installed, you can run from the project root:
    ```bash
    air
    ```
    This will watch for file changes and automatically rebuild and restart the application.
3.  Alternatively, to run the main server application without live reload from the project root:
    ```bash
    go run cmd/server/main.go
    ```

### Building the Binary

To build the application into a single executable binary:

1.  Use the `make build` command from the project root, which utilizes the `Makefile`:
    ```bash
    make build
    ```
    This command will compile the application and create an executable file named `server` (or as specified in the `Makefile`) in the project's root directory.
2.  Alternatively, you can run the Go build command directly from the project root:
    ```bash
    go build -o server cmd/server/main.go
    ```
    (Replace `server` with your desired output binary name).

After building, you can run the application directly (e.g., `./server`). Remember to have the necessary configurations available (either through environment variables or a configuration file).

## API Endpoints

The application exposes both HTTP/RESTful APIs and gRPC services.

-   **HTTP/RESTful APIs**:
    -   The main API route definitions can be found in `internal/api/router/` (e.g., in files like `api.go` or `router_v1.go`).
    -   HTTP request handlers that implement the business logic for these routes are located in `internal/api/handler/`.
    -   The project uses `swag` to generate Swagger/OpenAPI documentation from code annotations. You can generate this documentation using the `make docs` command (as seen in the `Makefile`). Once generated (typically into a `docs/` directory), the documentation (often a `swagger.json` or `swagger.yaml` file and a UI) will provide detailed information about all available RESTful endpoints, request/response payloads, and authentication methods.

-   **gRPC Services**:
    -   Service definitions using Protocol Buffers (`.proto` files) are located in the `protobuf/` directory.
    -   The gRPC services are set up and registered within the application, typically configured via `bootstrap/service_provider/grpc_provider.go` and `bootstrap/service_provider/grpc_router_provider.go`.
    -   To understand the available gRPC services, methods, and message types, refer to the `.proto` files in the `protobuf/` directory.

For a comprehensive understanding of all available API endpoints and their usage, it is recommended to:
1.  Generate and consult the Swagger/OpenAPI documentation for RESTful APIs.
2.  Examine the `.proto` files for gRPC service definitions.
3.  Explore the routing and handler code in `internal/api/` as needed.

It is not practical to list all specific endpoints in this README; the sources mentioned above are the most accurate and up-to-date references.

## Testing

This project uses Go's built-in testing capabilities. Tests are typically written in files named `*_test.go` and are located in the same package as the code they test.

To run all tests in the project, you can use the standard Go test command from the project root:
```bash
go test ./...
```
This command will discover and execute tests in all subdirectories.

For more specific testing:
- To run tests for a specific package: `go test ./path/to/package` (e.g., `go test ./internal/application/usecase/user`)
- To run a specific test function: `go test ./path/to/package -run TestSpecificFunctionName`

While the `Makefile` includes a `vet` target (`make vet`) for static analysis with `go vet ./...`, there isn't an explicit `make test` target defined. Use the `go test` commands directly for executing tests.

To generate a test coverage report:
```bash
# From the project root
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```
This will run the tests, create a coverage profile (`coverage.out`), and then open an HTML report in your browser showing code coverage.

## Contributing

Contributions are welcome! If you'd like to improve this project, please follow these steps:

1.  **Fork the repository** on GitHub.
2.  **Create a new branch** for your feature or bug fix:
    ```bash
    git checkout -b feature/your-feature-name
    ```
    Or for a bug fix:
    ```bash
    git checkout -b fix/issue-description
    ```
3.  **Make your changes**: Implement your feature or fix the bug.
4.  **Test your changes**: Ensure all existing tests pass and add new tests for your changes if applicable. Run linters and static analysis tools:
    ```bash
    # Run all tests from the project root
    go test ./...

    # Run Go vet for static analysis (as per Makefile) from the project root
    make vet
    ```
5.  **Commit your changes**: Write clear and descriptive commit messages.
    ```bash
    git commit -am 'feat: Add some amazing feature' # Example for a new feature
    git commit -am 'fix: Resolve issue #123'   # Example for a bug fix
    ```
    While the GitHub Actions workflow (`.github/workflows/commit.yml`) primarily handles notifications, adhering to good commit message practices (e.g., a clear summary, detailed explanation if needed) is encouraged.
6.  **Push to your branch**:
    ```bash
    git push origin feature/your-feature-name
    ```
7.  **Create a new Pull Request (PR)**: Go to the original repository on GitHub and open a PR from your forked branch. Provide a clear title and description for your PR.

Your contribution will be reviewed, and once approved, it will be merged into the main codebase. Thank you for helping to improve the project!

## License

This project does not currently have a specific `LICENSE` file in its root directory. Users should consider adding one to clarify the terms under which the software can be used, modified, and distributed (e.g., MIT, Apache 2.0).

For now, please assume the code is proprietary and all rights are reserved by the original authors, unless otherwise explicitly stated by the project maintainers. If you plan to use or contribute to this project, it is advisable to contact the maintainers to clarify licensing terms.
