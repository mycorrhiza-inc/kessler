// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: file_text.sql

package dbstore

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createFileTextSource = `-- name: CreateFileTextSource :one
INSERT INTO public.file_text_source (
		file_id,
		is_original_text,
		language,
		text,
		created_at,
		updated_at
	)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING file_id, is_original_text, language, text, id, created_at, updated_at
`

type CreateFileTextSourceParams struct {
	FileID         pgtype.UUID
	IsOriginalText bool
	Language       string
	Text           pgtype.Text
}

func (q *Queries) CreateFileTextSource(ctx context.Context, arg CreateFileTextSourceParams) (FileTextSource, error) {
	row := q.db.QueryRow(ctx, createFileTextSource,
		arg.FileID,
		arg.IsOriginalText,
		arg.Language,
		arg.Text,
	)
	var i FileTextSource
	err := row.Scan(
		&i.FileID,
		&i.IsOriginalText,
		&i.Language,
		&i.Text,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteFileTexts = `-- name: DeleteFileTexts :exec
DELETE FROM public.file_text_source
WHERE file_id = $1
`

func (q *Queries) DeleteFileTexts(ctx context.Context, fileID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteFileTexts, fileID)
	return err
}

const listTextsOfFile = `-- name: ListTextsOfFile :many
SELECT file_id, is_original_text, language, text, id, created_at, updated_at
FROM public.file_text_source
WHERE file_id = $1
`

func (q *Queries) ListTextsOfFile(ctx context.Context, fileID pgtype.UUID) ([]FileTextSource, error) {
	rows, err := q.db.Query(ctx, listTextsOfFile, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FileTextSource
	for rows.Next() {
		var i FileTextSource
		if err := rows.Scan(
			&i.FileID,
			&i.IsOriginalText,
			&i.Language,
			&i.Text,
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

const listTextsOfFileOriginal = `-- name: ListTextsOfFileOriginal :many
SELECT file_id, is_original_text, language, text, id, created_at, updated_at
FROM public.file_text_source
WHERE file_id = $1 and is_original_text = true
`

func (q *Queries) ListTextsOfFileOriginal(ctx context.Context, fileID pgtype.UUID) ([]FileTextSource, error) {
	rows, err := q.db.Query(ctx, listTextsOfFileOriginal, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FileTextSource
	for rows.Next() {
		var i FileTextSource
		if err := rows.Scan(
			&i.FileID,
			&i.IsOriginalText,
			&i.Language,
			&i.Text,
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

const listTextsOfFileWithLanguage = `-- name: ListTextsOfFileWithLanguage :many
SELECT file_id, is_original_text, language, text, id, created_at, updated_at
FROM public.file_text_source
WHERE file_id = $1 and language = $2
`

type ListTextsOfFileWithLanguageParams struct {
	FileID   pgtype.UUID
	Language string
}

func (q *Queries) ListTextsOfFileWithLanguage(ctx context.Context, arg ListTextsOfFileWithLanguageParams) ([]FileTextSource, error) {
	rows, err := q.db.Query(ctx, listTextsOfFileWithLanguage, arg.FileID, arg.Language)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FileTextSource
	for rows.Next() {
		var i FileTextSource
		if err := rows.Scan(
			&i.FileID,
			&i.IsOriginalText,
			&i.Language,
			&i.Text,
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
