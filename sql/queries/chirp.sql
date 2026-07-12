-- name: CreateChirp :one
INSERT INTO chirps(
    id,
    body,
    user_id,
    created_at,
    updated_at
) VALUES (
    gen_random_uuid(),
    $1,
    $2,
    now(),
    now()
)
RETURNING *;

-- name: ListChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: ListChirpsFromUserID :many
SELECT * FROM chirps
WHERE id = $1
ORDER BY created_at ASC;

-- name: GetChirpByID :one
SELECT * FROM chirps
WHERE id = $1;

-- name: DeleteChirpByID :exec
DELETE FROM chirps
WHERE id = $1;
