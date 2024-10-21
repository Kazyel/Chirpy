-- name: GetRefreshToken :one
SELECT * FROM refresh_token
WHERE token = $1;