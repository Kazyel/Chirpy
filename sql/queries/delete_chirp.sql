-- name: DeleteChirp :one
DELETE FROM chirps WHERE id = $1 and user_id = $2
RETURNING *;