-- name: GetUserByRefreshToken :one
SELECT users.*, refresh_token.user_id, refresh_token.token
FROM
    refresh_token
INNER JOIN users ON refresh_token.user_id = users.id
WHERE refresh_token.token = $1;