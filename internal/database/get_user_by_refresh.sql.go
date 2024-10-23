// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: get_user_by_refresh.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const getUserByRefreshToken = `-- name: GetUserByRefreshToken :one
SELECT users.id, users.created_at, users.updated_at, users.email, users.hashed_password, users.is_chirpy_red, refresh_token.user_id, refresh_token.token
FROM
    refresh_token
INNER JOIN users ON refresh_token.user_id = users.id
WHERE refresh_token.token = $1
`

type GetUserByRefreshTokenRow struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Email          string
	HashedPassword string
	IsChirpyRed    bool
	UserID         uuid.UUID
	Token          string
}

func (q *Queries) GetUserByRefreshToken(ctx context.Context, token string) (GetUserByRefreshTokenRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByRefreshToken, token)
	var i GetUserByRefreshTokenRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
		&i.UserID,
		&i.Token,
	)
	return i, err
}
