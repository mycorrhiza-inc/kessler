// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: thaum_apikeys.sql

package dbstore

import (
	"context"
)

const checkIfThaumaturgyAPIKeyExists = `-- name: CheckIfThaumaturgyAPIKeyExists :one
SELECT
    key_name, key_blake3_hash, id, created_at, updated_at
FROM
    userfiles.thaumaturgy_api_keys
WHERE
    key_blake3_hash = $1
`

func (q *Queries) CheckIfThaumaturgyAPIKeyExists(ctx context.Context, keyBlake3Hash string) (UserfilesThaumaturgyApiKey, error) {
	row := q.db.QueryRow(ctx, checkIfThaumaturgyAPIKeyExists, keyBlake3Hash)
	var i UserfilesThaumaturgyApiKey
	err := row.Scan(
		&i.KeyName,
		&i.KeyBlake3Hash,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
