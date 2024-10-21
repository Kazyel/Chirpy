-- name: RevokeToken :exec
UPDATE refresh_token
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;