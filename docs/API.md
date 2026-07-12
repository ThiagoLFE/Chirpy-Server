# Chirpy API reference

Base URL: `http://localhost:8080`

All request and response bodies use JSON unless stated otherwise. Errors use this shape:

```json
{
  "error": "message"
}
```

## Authentication

Login returns two tokens:

- `token`: a JWT access token valid for one hour. Send it as `Authorization: Bearer <token>`.
- `refresh_token`: a database-backed token valid for 60 days. Send it as `Authorization: Bearer <refresh_token>` to `/api/refresh` or `/api/revoke`.

The user and chirp endpoints that require an access token are marked **Access token** below. The refresh-token endpoints do not accept an access token in place of a refresh token.

## Endpoint summary

| Method | Path | Authentication | Description |
| --- | --- | --- | --- |
| `GET` | `/app/` | None | Serves the static frontend |
| `GET` | `/api/healthz` | None | Readiness check |
| `GET` | `/admin/metrics` | None | HTML page showing `/app` visit count |
| `POST` | `/admin/reset` | None | Clears users and resets metrics in development |
| `POST` | `/api/users` | None | Creates a user |
| `GET` | `/api/users` | None | Lists users |
| `PUT` | `/api/users` | Access token | Updates the authenticated user |
| `POST` | `/api/login` | None | Authenticates a user |
| `POST` | `/api/refresh` | Refresh token | Issues a new access token |
| `POST` | `/api/revoke` | Refresh token | Revokes a refresh token |
| `GET` | `/api/refresh_tokens` | None | Lists stored refresh tokens |
| `POST` | `/api/chirps` | Access token | Creates a chirp |
| `GET` | `/api/chirps` | None | Lists chirps |
| `GET` | `/api/chirps/{id}` | None | Gets one chirp |
| `DELETE` | `/api/chirps/{id}` | Access token | Deletes one of the authenticated user's chirps |
| `POST` | `/api/polka/webhooks` | Polka API key | Marks a user as a Chirpy Red member |

## Health and administration

### `GET /api/healthz`

Returns plain text `OK` with status `200`.

### `GET /admin/metrics`

Returns an HTML page containing the number of requests served under `/app`.

### `POST /admin/reset`

Resets the `/app` hit counter and deletes all users. It returns `403` unless `PLATFORM=dev`.

Successful response:

```json
{
  "status": "clear"
}
```

This also deletes related chirps and refresh tokens through the database foreign keys.

## Users and authentication

### `POST /api/users`

Creates a user. The current handler requires a non-empty email and a password with at least five characters.

Request:

```json
{
  "email": "reader@example.com",
  "password": "bananaSplit"
}
```

Response: `201 Created`

```json
{
  "id": "3311741c-680c-4546-99f3-fc9efac2036c",
  "email": "reader@example.com",
  "is_chirpy_red": false,
  "created_at": "2026-07-12T12:00:00Z",
  "updated_at": "2026-07-12T12:00:00Z"
}
```

### `GET /api/users`

Returns users ordered by creation time. Password hashes are not included in the response.

Response: `200 OK`

```json
[
  {
    "id": "3311741c-680c-4546-99f3-fc9efac2036c",
    "email": "reader@example.com",
    "is_chirpy_red": false,
    "created_at": "2026-07-12T12:00:00Z",
    "updated_at": "2026-07-12T12:00:00Z"
  }
]
```

### `PUT /api/users`

Updates the email and password of the user identified by the access token. The request body has the same shape as user creation.

Header:

```text
Authorization: Bearer ACCESS_TOKEN
```

Response: `200 OK` with the updated user object.

### `POST /api/login`

Authenticates a user.

Request:

```json
{
  "email": "reader@example.com",
  "password": "bananaSplit"
}
```

Response: `200 OK`

```json
{
  "id": "3311741c-680c-4546-99f3-fc9efac2036c",
  "email": "reader@example.com",
  "is_chirpy_red": false,
  "created_at": "2026-07-12T12:00:00Z",
  "updated_at": "2026-07-12T12:00:00Z",
  "token": "JWT_ACCESS_TOKEN",
  "refresh_token": "REFRESH_TOKEN"
}
```

### `POST /api/refresh`

Exchanges a valid, non-expired, non-revoked refresh token for a new access token.

Header:

```text
Authorization: Bearer REFRESH_TOKEN
```

Response: `200 OK`

```json
{
  "token": "NEW_JWT_ACCESS_TOKEN"
}
```

### `POST /api/revoke`

Revokes a refresh token. The endpoint returns `204 No Content` on success.

Header:

```text
Authorization: Bearer REFRESH_TOKEN
```

### `GET /api/refresh_tokens`

Lists stored refresh-token records. This is a development/admin inspection endpoint and currently has no authentication middleware. Treat its response as sensitive.

## Chirps

### `POST /api/chirps`

Creates a chirp for the authenticated user. Leading and trailing whitespace is removed, and the body must contain between 1 and 140 characters.

Header:

```text
Authorization: Bearer ACCESS_TOKEN
```

Request:

```json
{
  "body": "Hello from Chirpy!"
}
```

Response: `201 Created`

```json
{
  "id": "c4bcd112-16bb-4109-8c4c-3c4c84ddf488",
  "body": "Hello from Chirpy!",
  "user_id": "3311741c-680c-4546-99f3-fc9efac2036c",
  "created_at": "2026-07-12T12:00:00Z",
  "updated_at": "2026-07-12T12:00:00Z"
}
```

### `GET /api/chirps`

Lists chirps ordered from oldest to newest by default.

Query parameters:

| Parameter | Description |
| --- | --- |
| `author_id` | UUID used to request chirps for one author |
| `sort=desc` | Sort the returned list from newest to oldest |

Response: `200 OK` with an array of chirp objects.

### `GET /api/chirps/{id}`

Returns one chirp by UUID. Returns `404 Not Found` when the chirp does not exist.

### `DELETE /api/chirps/{id}`

Deletes a chirp only when its `user_id` matches the authenticated user. Returns `204 No Content` on success and `403 Forbidden` when another user owns the chirp.

## Polka webhook

### `POST /api/polka/webhooks`

Authenticates with the configured Polka API key:

```text
Authorization: ApiKey POLKA_KEY
```

For a membership upgrade, send:

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "3311741c-680c-4546-99f3-fc9efac2036c"
  }
}
```

The user is marked as Chirpy Red and the endpoint returns `204 No Content`. Events other than `user.upgraded` are acknowledged with `204` without changing the user.

## Current implementation notes

- The password validation code requires five characters, although some error text says eight. This documentation reflects the actual validation rule.
- The `author_id` route parameter is exposed by the HTTP handler, but the current SQL query in [`sql/queries/chirp.sql`](../sql/queries/chirp.sql) compares the chirp `id` column instead of `user_id`. If author filtering is required, update that query and regenerate the sqlc package.
