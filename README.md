# Authentication Service

Authentication service for a Hotel Management System, built with Go and designed as a standalone backend service responsible for authentication and authorization.

## Overview

This project provides the authentication layer for the hotel management system. It is designed as an independent service so that other services in the system can rely on it for user authentication and authorization.

The project is also structured as a portfolio/hobby project to demonstrate production-oriented Go backend development, testing, code quality, and CI practices.

## Features

- User authentication
- Password-based authentication
- JWT-based authorization
- Secure password handling
- Token validation
- RESTful APIs
- Database-backed user management
- Input validation
- Unit testing
- Static code analysis

## Tech Stack

| Technology | Purpose |
| ---------- | ------- |
| Go | Backend service |
| PostgreSQL | Persistent data storage |
| JWT | Authentication and authorization |
| golangci-lint | Static analysis and linting |
| gofmt | Code formatting |
| go vet | Go static analysis |

## Project Structure

```
authentication/
├── cmd/
│   └── server/
├── database/
├── dto/
├── handler/
├── migrations/
├── mocks/
├── model/
├── repository/
├── router/
├── service/
├── Makefile
├── env.vars
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## Getting Started

### Prerequisites

Make sure you have the following installed:
- Go 1.25+
- PostgreSQL
- Git
- [golang-migrate](https://github.com/golang-migrate/migrate) (for database migrations)

### Clone the Repository

```bash
git clone <repository-url>
cd authentication
```

### Environment Configuration

The project uses `env.vars` for environment configuration. Ensure the variables match your local environment.

Example `env.vars`:
```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/authentication?sslmode=disable"
export JWT_SECRET_STRING="your-jwt-secret-string"
```

Never commit real secrets to the repository. Note that `env.vars` in this repo currently contains local development defaults.

### Install Dependencies

```bash
go mod download
```

### Database Migrations

Run the migrations to set up the database schema:
```bash
make migrate-up
```

### Run the Application

You can run the application using Make:
```bash
make run
```
Or directly using Go:
```bash
go run ./cmd/server/main.go
```

## Testing

Run tests using Make:
```bash
make test
```
Or directly using Go:
```bash
go test ./repository/...
go test ./service/...
```

## Code Quality

The project uses standard Go tooling and golangci-lint. A `Makefile` is provided for convenience.

### Code Formatting

Automatically format the code:
```bash
make fmt
```

### Static Analysis and Linting

Run all quality checks (formatting, vet, linting, tests):
```bash
make check
```

Run Go Vet:
```bash
make vet
```

Run GolangCI-Lint:
```bash
make lint
```

## API

The service exposes REST APIs for authentication.

Endpoints:
- `POST /signup` - Register a new user
- `POST /login` - Authenticate and receive a JWT

## Security

Authentication-related services handle sensitive information. The following practices are followed:
- Passwords are never stored as plaintext (hashed securely).
- Authentication tokens are signed using a secret key.
- Secrets are provided through environment variables.
- Authentication endpoints are validated before processing requests.

For production deployment, use a properly managed secret-management solution rather than committing secrets to configuration files.

## License

This project is licensed under the MIT License.
See the [LICENSE](LICENSE) file for details.

## Author

**Gokul Sujan**

This project is developed as a hobby/portfolio project to demonstrate backend engineering, Go development, authentication architecture, testing, and CI practices.
