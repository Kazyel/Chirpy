-- name: UpdateUser :exec
UPDATE users
SET email = $2, hashed_password = $3 
WHERE id = $1;