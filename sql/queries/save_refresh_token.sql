-- name: SaveRefreshToken :one
INSERT INTO refresh_token(token, created_at, updated_at, user_id, expires_at)
VALUES ($1, NOW(), NOW(), $2, $3)
RETURNING *;