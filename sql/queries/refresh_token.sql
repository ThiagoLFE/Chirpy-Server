-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(
    token,
    user_id,

    created_at,
    updated_at,
    expires_at

) VALUES (
    $1,
    $2,

    now(),
    now(),
    $3
)
RETURNING *;

-- name: GetRefreshTokenByRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteRefreshTokenByUserID :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = now(), updated_at = now()
WHERE token = $1
RETURNING *;


-- name: ListRefreshTokens :many
SELECT * FROM refresh_tokens;
