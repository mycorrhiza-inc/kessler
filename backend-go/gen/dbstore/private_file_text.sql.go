// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: private_file_text.sql

package dbstore

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPrivateFileTextSource = `-- name: CreatePrivateFileTextSource :one
INSERT INTO userfiles.private_file_text_source (
    file_id,
    is_original_text,
    language,
    text,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id
`

type CreatePrivateFileTextSourceParams struct {
	FileID         pgtype.UUID
	IsOriginalText bool
	Language       string
	Text           pgtype.Text
}

func (q *Queries) CreatePrivateFileTextSource(ctx context.Context, arg CreatePrivateFileTextSourceParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, createPrivateFileTextSource,
		arg.FileID,
		arg.IsOriginalText,
		arg.Language,
		arg.Text,
	)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const deletePrivateFileTexts = `-- name: DeletePrivateFileTexts :exec
DELETE FROM userfiles.private_file_text_source
WHERE file_id = $1
`

func (q *Queries) DeletePrivateFileTexts(ctx context.Context, fileID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deletePrivateFileTexts, fileID)
	return err
}

const listPrivateTextsOfFile = `-- name: ListPrivateTextsOfFile :many
SELECT file_id, is_original_text, language, text, id, created_at, updated_at
FROM userfiles.private_file_text_source
WHERE file_id = $1
`

func (q *Queries) ListPrivateTextsOfFile(ctx context.Context, fileID pgtype.UUID) ([]UserfilesPrivateFileTextSource, error) {
	rows, err := q.db.Query(ctx, listPrivateTextsOfFile, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserfilesPrivateFileTextSource
	for rows.Next() {
		var i UserfilesPrivateFileTextSource
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

const listPrivateTextsOfFileOriginal = `-- name: ListPrivateTextsOfFileOriginal :many
SELECT file_id, is_original_text, language, text, id, created_at, updated_at
FROM userfiles.private_file_text_source
WHERE file_id = $1 and is_original_text = true
`

func (q *Queries) ListPrivateTextsOfFileOriginal(ctx context.Context, fileID pgtype.UUID) ([]UserfilesPrivateFileTextSource, error) {
	rows, err := q.db.Query(ctx, listPrivateTextsOfFileOriginal, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserfilesPrivateFileTextSource
	for rows.Next() {
		var i UserfilesPrivateFileTextSource
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

const listPrivateTextsOfFileWithLanguage = `-- name: ListPrivateTextsOfFileWithLanguage :many
SELECT file_id, is_original_text, language, text, id, created_at, updated_at
FROM userfiles.private_file_text_source
WHERE file_id = $1 and language = $2
`

type ListPrivateTextsOfFileWithLanguageParams struct {
	FileID   pgtype.UUID
	Language string
}

func (q *Queries) ListPrivateTextsOfFileWithLanguage(ctx context.Context, arg ListPrivateTextsOfFileWithLanguageParams) ([]UserfilesPrivateFileTextSource, error) {
	rows, err := q.db.Query(ctx, listPrivateTextsOfFileWithLanguage, arg.FileID, arg.Language)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserfilesPrivateFileTextSource
	for rows.Next() {
		var i UserfilesPrivateFileTextSource
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
