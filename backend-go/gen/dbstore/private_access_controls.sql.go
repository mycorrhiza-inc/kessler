// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: private_access_controls.sql

package dbstore

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const checkOperatorAccessToObject = `-- name: CheckOperatorAccessToObject :one
SELECT EXISTS(
    SELECT 1
    FROM userfiles.private_access_controls
    WHERE operator_id = $1 AND object_id = $2
)
`

type CheckOperatorAccessToObjectParams struct {
	OperatorID pgtype.UUID
	ObjectID   pgtype.UUID
}

func (q *Queries) CheckOperatorAccessToObject(ctx context.Context, arg CheckOperatorAccessToObjectParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkOperatorAccessToObject, arg.OperatorID, arg.ObjectID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createPrivateAccessControl = `-- name: CreatePrivateAccessControl :one
INSERT INTO userfiles.private_access_controls (
		operator_id,
    operator_table,
		object_id,
    object_table,
		created_at,
		updated_at
	)
VALUES ($1, $2,$3, $4, NOW(), NOW())
RETURNING operator_id, operator_table, object_id, object_table, id, created_at, updated_at
`

type CreatePrivateAccessControlParams struct {
	OperatorID    pgtype.UUID
	OperatorTable string
	ObjectID      pgtype.UUID
	ObjectTable   string
}

func (q *Queries) CreatePrivateAccessControl(ctx context.Context, arg CreatePrivateAccessControlParams) (UserfilesPrivateAccessControl, error) {
	row := q.db.QueryRow(ctx, createPrivateAccessControl,
		arg.OperatorID,
		arg.OperatorTable,
		arg.ObjectID,
		arg.ObjectTable,
	)
	var i UserfilesPrivateAccessControl
	err := row.Scan(
		&i.OperatorID,
		&i.OperatorTable,
		&i.ObjectID,
		&i.ObjectTable,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAccessControl = `-- name: DeleteAccessControl :exec
DELETE FROM userfiles.private_access_controls
WHERE id = $1
`

func (q *Queries) DeleteAccessControl(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteAccessControl, id)
	return err
}

const listAcessesForOperator = `-- name: ListAcessesForOperator :many
SELECT operator_id, operator_table, object_id, object_table, id, created_at, updated_at
FROM userfiles.private_access_controls
WHERE operator_id = $1
`

func (q *Queries) ListAcessesForOperator(ctx context.Context, operatorID pgtype.UUID) ([]UserfilesPrivateAccessControl, error) {
	rows, err := q.db.Query(ctx, listAcessesForOperator, operatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserfilesPrivateAccessControl
	for rows.Next() {
		var i UserfilesPrivateAccessControl
		if err := rows.Scan(
			&i.OperatorID,
			&i.OperatorTable,
			&i.ObjectID,
			&i.ObjectTable,
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listOperatorsCanAcessObject = `-- name: ListOperatorsCanAcessObject :many
SELECT operator_id, operator_table, object_id, object_table, id, created_at, updated_at
FROM userfiles.private_access_controls
WHERE object_id = $1
`

func (q *Queries) ListOperatorsCanAcessObject(ctx context.Context, objectID pgtype.UUID) ([]UserfilesPrivateAccessControl, error) {
	rows, err := q.db.Query(ctx, listOperatorsCanAcessObject, objectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserfilesPrivateAccessControl
	for rows.Next() {
		var i UserfilesPrivateAccessControl
		if err := rows.Scan(
			&i.OperatorID,
			&i.OperatorTable,
			&i.ObjectID,
			&i.ObjectTable,
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const revokeAccessForOperatorOnObject = `-- name: RevokeAccessForOperatorOnObject :exec
DELETE FROM userfiles.private_access_controls
WHERE operator_id = $1 AND object_id = $2
`

type RevokeAccessForOperatorOnObjectParams struct {
	OperatorID pgtype.UUID
	ObjectID   pgtype.UUID
}

func (q *Queries) RevokeAccessForOperatorOnObject(ctx context.Context, arg RevokeAccessForOperatorOnObjectParams) error {
	_, err := q.db.Exec(ctx, revokeAccessForOperatorOnObject, arg.OperatorID, arg.ObjectID)
	return err
}
