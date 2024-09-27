// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: files_events.sql

package dbstore

import (
	"context"

	"github.com/google/uuid"
)

const createFilesAssociatedWithEvent = `-- name: CreateFilesAssociatedWithEvent :one
INSERT INTO public.relation_files_events (
		file_id,
		event_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id
`

type CreateFilesAssociatedWithEventParams struct {
	FileID  uuid.UUID
	EventID uuid.UUID
}

func (q *Queries) CreateFilesAssociatedWithEvent(ctx context.Context, arg CreateFilesAssociatedWithEventParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createFilesAssociatedWithEvent, arg.FileID, arg.EventID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteFilesAssociatedWithEvent = `-- name: DeleteFilesAssociatedWithEvent :exec
DELETE FROM public.relation_files_events
WHERE id = $1
`

func (q *Queries) DeleteFilesAssociatedWithEvent(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteFilesAssociatedWithEvent, id)
	return err
}

const listFilesAssociatedWithEvent = `-- name: ListFilesAssociatedWithEvent :many
SELECT file_id, event_id, id, created_at, updated_at
FROM public.relation_files_events
ORDER BY created_at DESC
`

func (q *Queries) ListFilesAssociatedWithEvent(ctx context.Context) ([]RelationFilesEvent, error) {
	rows, err := q.db.QueryContext(ctx, listFilesAssociatedWithEvent)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RelationFilesEvent
	for rows.Next() {
		var i RelationFilesEvent
		if err := rows.Scan(
			&i.FileID,
			&i.EventID,
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const readFilesAssociatedWithEvent = `-- name: ReadFilesAssociatedWithEvent :one
SELECT file_id, event_id, id, created_at, updated_at
FROM public.relation_files_events
WHERE id = $1
`

func (q *Queries) ReadFilesAssociatedWithEvent(ctx context.Context, id uuid.UUID) (RelationFilesEvent, error) {
	row := q.db.QueryRowContext(ctx, readFilesAssociatedWithEvent, id)
	var i RelationFilesEvent
	err := row.Scan(
		&i.FileID,
		&i.EventID,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateFilesAssociatedWithEvent = `-- name: UpdateFilesAssociatedWithEvent :one
UPDATE public.relation_files_events
SET file_id = $1,
	event_id = $2,
	updated_at = NOW()
WHERE id = $3
RETURNING id
`

type UpdateFilesAssociatedWithEventParams struct {
	FileID  uuid.UUID
	EventID uuid.UUID
	ID      uuid.UUID
}

func (q *Queries) UpdateFilesAssociatedWithEvent(ctx context.Context, arg UpdateFilesAssociatedWithEventParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, updateFilesAssociatedWithEvent, arg.FileID, arg.EventID, arg.ID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}