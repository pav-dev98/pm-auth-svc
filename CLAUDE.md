# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

gRPC authentication microservice in Go (part of the "Pure Minerals" microservices backend). Exposes Register/Login/Refresh over gRPC, persists credentials and sessions in PostgreSQL (GORM with AutoMigrate), and issues HS256 JWTs. The protobuf definitions live in an external module: `github.com/pav-dev98/pm-proto` (`pb "github.com/pav-dev98/pm-proto/auth"`).

## Environment notes

- The machine uses **Podman**, not Docker. `docker` may be aliased to `podman` in the user's interactive shell, but in non-interactive shells only `podman` exists — use `podman` directly.
- The compose file is named **`Docker-compose.yml`** (capital D), so compose commands need `-f Docker-compose.yml`.

## Commands

```bash
go build ./...                  # build everything
go run ./cmd/main.go            # run the service (gRPC on :50051)
go test ./...                   # all tests (integration tests need a container runtime)
go test ./internal/application/usecase/                          # unit tests only
go test ./internal/application/usecase/ -run TestLogin -v        # single test
go generate ./...               # regenerate mocks (mockgen, into internal/application/ports/mocks/)

podman compose -f Docker-compose.yml up --build   # full stack (postgres + app)
podman compose -f Docker-compose.yml up -d postgres  # DB only, run app locally
```

Integration tests (`internal/infrastructure/persistence/postgress/`) use testcontainers-go to spin up an ephemeral PostgreSQL container, so they require the container runtime to be available.

Manual testing: the gRPC server has reflection enabled, so `grpcurl -plaintext -d '{...}' localhost:50051 auth.AuthService/Register` works (see README).

## Architecture

Clean/hexagonal architecture. Dependencies point inward; the domain knows nothing about infrastructure.

```
cmd/main.go        → composition root: wires config → infra → use cases → gRPC server
config/            → loads .env (godotenv) + env vars into a Config struct (DSN, JWT settings)
internal/
  domain/          → entities (AuthCredential, Session) and domain errors (ErrNotFound, ErrDatabase, ...)
  application/
    ports/         → interfaces (AuthRepository, PasswordHasher, TokenService) with //go:generate mockgen directives
    ports/mocks/   → generated mocks (go.uber.org/mock) used by unit tests
    usecase/       → RegisterCredential, LoginCredential, RefreshCredential — depend only on ports
  infrastructure/
    persistence/postgress/  → GORM AuthRepository (note the package spelling "postgress")
    security/bcrypt/        → PasswordHasher implementation
    security/jwt/           → TokenService implementation (access + refresh tokens)
  interfaces/grpc/ → AuthServer adapter: maps proto requests to use cases and domain errors to gRPC status codes
```

Conventions to follow:

- **Persistence mapping**: GORM models (`AuthCredentialModel`, `SessionModel`) are separate from domain entities. Conversions are struct methods sharing the same names across models: `(m *XModel) toDomain()` and `(m *XModel) fromDomain(d)`. Keep this pattern for new models — `fromDomain` fills the receiver (it can't be a method on domain types without an import cycle).
- **New use case flow**: define/extend a port interface in `application/ports` (with a `go:generate mockgen` line), implement it in `infrastructure/`, add the use case in `application/usecase`, expose it via `interfaces/grpc`, and wire it in `cmd/main.go`.
- Repository methods translate `gorm.ErrRecordNotFound` into `domain.ErrNotFound` and other DB errors into `domain.ErrDatabase`; gRPC handlers map domain errors to status codes.
- Code comments and log messages are in Spanish; commit messages follow Conventional Commits in English.
