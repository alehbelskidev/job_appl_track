# Job Application Tracker — API

Go REST API for tracking job applications. JWT authentication, PostgreSQL, modular architecture.

## Stack

- **Go** — net/http, Chi router
- **PostgreSQL** — via pgx/v5
- **sqlc** — type-safe SQL query generation
- **golang-migrate** — database migrations
- **JWT** — access + refresh token auth

## Project Structure

```
cmd/
  app/        — server entrypoint
  migrate/    — migration runner
internal/
  auth/       — register, login, refresh (handler, service, module)
  applications/ — CRUD for job applications (handler, module)
  repo/       — sqlc generated code
  shared/     — config, middleware
sql/
  migrations/ — up/down migration files
  queries/    — sqlc query definitions
```

## Prerequisites

- Go 1.25+
- Docker + Docker Compose
- sqlc: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

## Local Development

```bash
cp .env.example .env
docker compose up
```

The API starts with hot reload via [air](https://github.com/air-verse/air) on `localhost:3001`.

## Environment Variables

| Variable | Description |
|---|---|
| `POSTGRES_USER` | Database user |
| `POSTGRES_PASSWORD` | Database password |
| `POSTGRES_DB` | Database name |
| `POSTGRES_HOST` | Database host (`postgres` in Docker, `localhost` for migrations) |
| `POSTGRES_HOST_MIGRATE` | Host for migration runner |
| `JWT_SECRET` | Secret key for signing JWT tokens |
| `ALLOWED_ORIGINS` | Comma-separated list of allowed CORS origins |

## Database Migrations

Migrations are in `sql/migrations/` following the `000N_name.up.sql` / `000N_name.down.sql` convention.

Run locally:
```bash
go run ./cmd/migrate
```

## Regenerating sqlc Types

After modifying SQL queries or schema:
```bash
sqlc generate
```

## Production Build

```bash
docker build -f Dockerfile.prod -t job-tracker-api:latest .
```

Produces two binaries: `server` and `migrate`.

## API Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/auth/register` | — | Register new user |
| POST | `/auth/login` | — | Login, returns tokens |
| POST | `/auth/refresh` | — | Refresh access token |
| GET | `/api/applications` | ✓ | List user's applications |
| POST | `/api/applications` | ✓ | Create application |
| PATCH | `/api/applications/:id` | ✓ | Update application |
