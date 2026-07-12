-- name: CreateUser :one
INSERT INTO users ( id, email, hashed_password, created_at, updated_at )
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    now(),
    now()
)
RETURNING *;

-- name: ClearUsers :exec
DELETE FROM users;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at ASC;

-- name: UpdateUser :one
UPDATE users
SET email = $2, hashed_password = $3, updated_at = now()
WHERE id = $1
RETURNING *;
