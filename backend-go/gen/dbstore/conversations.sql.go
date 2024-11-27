// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: conversations.sql

package dbstore

import (
	"context"

	"github.com/google/uuid"
)

const docketConversationCreate = `-- name: DocketConversationCreate :one
INSERT INTO
    public.docket_conversations (
        docket_id,
        state,
        created_at,
        updated_at
    )
VALUES
    ($1, $2, NOW(), NOW())
RETURNING
    id
`

type DocketConversationCreateParams struct {
	DocketID string
	State    string
}

func (q *Queries) DocketConversationCreate(ctx context.Context, arg DocketConversationCreateParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, docketConversationCreate, arg.DocketID, arg.State)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const docketConversationDelete = `-- name: DocketConversationDelete :exec
DELETE FROM
    public.docket_conversations
WHERE
    id = $1
`

func (q *Queries) DocketConversationDelete(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, docketConversationDelete, id)
	return err
}

const docketConversationFetchByDocketIdMatch = `-- name: DocketConversationFetchByDocketIdMatch :many
SELECT
    id, docket_id, state, created_at, updated_at, deleted_at, name, description
FROM
    public.docket_conversations
WHERE
    docket_id = $1
`

func (q *Queries) DocketConversationFetchByDocketIdMatch(ctx context.Context, docketID string) ([]DocketConversation, error) {
	rows, err := q.db.Query(ctx, docketConversationFetchByDocketIdMatch, docketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DocketConversation
	for rows.Next() {
		var i DocketConversation
		if err := rows.Scan(
			&i.ID,
			&i.DocketID,
			&i.State,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Name,
			&i.Description,
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

const docketConversationList = `-- name: DocketConversationList :many
SELECT
    id, docket_id, state, created_at, updated_at, deleted_at, name, description
FROM
    public.docket_conversations
ORDER BY
    created_at DESC
`

func (q *Queries) DocketConversationList(ctx context.Context) ([]DocketConversation, error) {
	rows, err := q.db.Query(ctx, docketConversationList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DocketConversation
	for rows.Next() {
		var i DocketConversation
		if err := rows.Scan(
			&i.ID,
			&i.DocketID,
			&i.State,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Name,
			&i.Description,
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

const docketConversationRead = `-- name: DocketConversationRead :one
SELECT
    id, docket_id, state, created_at, updated_at, deleted_at, name, description
FROM
    public.docket_conversations
WHERE
    id = $1
`

func (q *Queries) DocketConversationRead(ctx context.Context, id uuid.UUID) (DocketConversation, error) {
	row := q.db.QueryRow(ctx, docketConversationRead, id)
	var i DocketConversation
	err := row.Scan(
		&i.ID,
		&i.DocketID,
		&i.State,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Name,
		&i.Description,
	)
	return i, err
}

const docketConversationUpdate = `-- name: DocketConversationUpdate :one
UPDATE
    public.docket_conversations
SET
    docket_id = $1,
    state = $2,
    updated_at = NOW()
WHERE
    id = $3
RETURNING
    id
`

type DocketConversationUpdateParams struct {
	DocketID string
	State    string
	ID       uuid.UUID
}

func (q *Queries) DocketConversationUpdate(ctx context.Context, arg DocketConversationUpdateParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, docketConversationUpdate, arg.DocketID, arg.State, arg.ID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const docketDocumentDeleteAll = `-- name: DocketDocumentDeleteAll :exec
DELETE FROM
    public.docket_documents
WHERE
    docket_id = $1
`

func (q *Queries) DocketDocumentDeleteAll(ctx context.Context, docketID uuid.UUID) error {
	_, err := q.db.Exec(ctx, docketDocumentDeleteAll, docketID)
	return err
}

const docketDocumentInsert = `-- name: DocketDocumentInsert :one
INSERT INTO
    public.docket_documents (
        docket_id,
        file_id,
        created_at,
        updated_at
    )
VALUES
    ($1, $2, NOW(), NOW())
RETURNING
    docket_id
`

type DocketDocumentInsertParams struct {
	DocketID uuid.UUID
	FileID   uuid.UUID
}

func (q *Queries) DocketDocumentInsert(ctx context.Context, arg DocketDocumentInsertParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, docketDocumentInsert, arg.DocketID, arg.FileID)
	var docket_id uuid.UUID
	err := row.Scan(&docket_id)
	return docket_id, err
}

const docketDocumentUpdate = `-- name: DocketDocumentUpdate :one
UPDATE
    public.docket_documents
SET
    docket_id = $1,
    updated_at = NOW()
WHERE
    file_id = $2
RETURNING
    file_id
`

type DocketDocumentUpdateParams struct {
	DocketID uuid.UUID
	FileID   uuid.UUID
}

func (q *Queries) DocketDocumentUpdate(ctx context.Context, arg DocketDocumentUpdateParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, docketDocumentUpdate, arg.DocketID, arg.FileID)
	var file_id uuid.UUID
	err := row.Scan(&file_id)
	return file_id, err
}
