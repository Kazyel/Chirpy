-- name: GetChirpsDesc :many
SELECT
    *
FROM
    chirps
ORDER BY
    created_at DESC;