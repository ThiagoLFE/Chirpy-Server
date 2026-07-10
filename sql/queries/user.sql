-- name: CreateUser :one
INSERT INTO users ( id, email, created_at, updated_at )
VALUES (
    gen_random_uuid(),
    $1,
    now(),
    now()
)
RETURNING *;

-- name: ClearUsers :exec
DELETE FROM users;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at ASC;
