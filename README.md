# Chirpy Server

Chirpy is a small social API written in Go. Users can create accounts, sign in, publish short messages called **chirps**, and manage their own content. The project also includes refresh-token authentication, a protected user profile update endpoint, a development reset endpoint, and a simulated Polka membership webhook.

The project is intentionally small and uses the Go standard library for its HTTP server and routing. PostgreSQL stores users, chirps, and refresh tokens; `sqlc` generates the database access layer from the SQL files in [`sql/`](sql/).

## Features

- HTTP API built with `net/http`
- PostgreSQL persistence
- Argon2id password hashing
- One-hour JWT access tokens
- Sixty-day refresh tokens with revocation
- User management and chirp create/list/read/delete operations
- Chirpy Red membership updates through a Polka-style API-key webhook
- Static frontend served at `/app/`
- Bruno request collection for manual API testing

## Requirements

- Go 1.25 or newer
- PostgreSQL with support for the functions used by [`sql/schema/`](sql/schema/)
- Optional: [Bruno](https://www.usebruno.com/) for running the request collection

## Getting started

### 1. Configure PostgreSQL

Create a database and copy the example environment file:

```sh
cp .env.example .env
```

Update `DB_URL` in `.env` with your PostgreSQL connection string. The server reads these variables using [`godotenv`](https://github.com/joho/godotenv):

| Variable | Required | Description |
| --- | --- | --- |
| `DB_URL` | Yes | PostgreSQL connection string |
| `PLATFORM` | No | Set to `dev` to enable `POST /admin/reset` |
| `TOKEN_SECRET` | Recommended | Secret used to sign JWT access tokens |
| `POLKA_KEY` | Recommended | API key accepted by the membership webhook |

Apply the schema files in order. The repository does not currently include a migration CLI, so `psql` can be used directly:

```sh
psql "$DB_URL" -f sql/schema/001_users_table.sql
psql "$DB_URL" -f sql/schema/002_chirps.sql
psql "$DB_URL" -f sql/schema/003_users_password.sql
psql "$DB_URL" -f sql/schema/004_refresh_token.sql
psql "$DB_URL" -f sql/schema/005_users_chirpy_red.sql
```

### 2. Start the server

```sh
go mod download
go run .
```

The API listens on `http://localhost:8080`.

Check that it is running:

```sh
curl http://localhost:8080/api/healthz
# OK
```

### 3. Run tests

```sh
go test ./...
```

## Quick API example

Create a user, log in, then use the returned access token to create a chirp:

```sh
curl -X POST http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -d '{"email":"reader@example.com","password":"bananaSplit"}'

curl -X POST http://localhost:8080/api/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"reader@example.com","password":"bananaSplit"}'

curl -X POST http://localhost:8080/api/chirps \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer ACCESS_TOKEN' \
  -d '{"body":"Hello from Chirpy!"}'
```

The complete endpoint reference, request bodies, responses, and authentication flow are in [`docs/API.md`](docs/API.md).

## Bruno collection

The [`bruno/`](bruno/) directory contains requests grouped by authentication, users, chirps, admin operations, and the Polka webhook.

To use it:

1. Open the `bruno/` directory as a collection in Bruno.
2. Start the server and select the `chirpy` environment.
3. Run **Create User**, then **Login**. The login request stores `accessToken` and `refreshToken` automatically.
4. Run the authenticated user and chirp requests.

The collection uses `http://localhost:8080` directly and expects the server to be running locally.

## Project layout

```text
.
├── main.go                 # Server setup and route registration
├── handle_*.go             # HTTP handlers
├── middleware_auth.go      # JWT authentication middleware
├── internal/auth/          # Password, JWT, API-key, and refresh-token helpers
├── internal/database/      # sqlc-generated PostgreSQL access layer
├── sql/schema/             # Database schema migrations
├── sql/queries/            # SQL queries used by sqlc
├── bruno/                  # Bruno API collection
├── docs/API.md             # Detailed API reference
└── index.html              # Static content served under /app/
```

## Generating database code

The database package is generated with [sqlc](https://sqlc.dev/) from [`sqlc.yaml`](sqlc.yaml):

```sh
sqlc generate
```

Install `sqlc` separately if it is not already available on your machine.

## Design notes

- Passwords are stored as Argon2id hashes; plaintext passwords are never returned by the API.
- Access tokens are JWTs signed with `TOKEN_SECRET` and expire after one hour.
- Refresh tokens are stored in PostgreSQL, expire after 60 days, and can be revoked.
- `POST /admin/reset` is intentionally restricted to `PLATFORM=dev`; do not expose it in production.
- All JSON errors use the shape `{ "error": "message" }`.

## License

No license has been declared for this repository yet.
