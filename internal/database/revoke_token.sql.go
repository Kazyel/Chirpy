// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: revoke_token.sql

package database

import (
	"context"
)

const revokeToken = `-- name: RevokeToken :exec
UPDATE refresh_token
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1
`

func (q *Queries) RevokeToken(ctx context.Context, token string) error {
	_, err := q.db.ExecContext(ctx, revokeToken, token)
	return err
}
