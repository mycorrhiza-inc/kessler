// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: files.sql

package dbstore

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addMetadataToFile = `-- name: AddMetadataToFile :one
UPDATE public.file
SET metadata_id = $1
WHERE id = $2
RETURNING id
`

type AddMetadataToFileParams struct {
	MetadataID pgtype.UUID
	ID         pgtype.UUID
}

func (q *Queries) AddMetadataToFile(ctx context.Context, arg AddMetadataToFileParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, addMetadataToFile, arg.MetadataID, arg.ID)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const addStageLog = `-- name: AddStageLog :one
INSERT INTO public.stage_log (
    file_id,
    status,
    log
  )
VALUES ($1, $2, $3)
RETURNING id,
  file_id,
  status
`

type AddStageLogParams struct {
	FileID pgtype.UUID
	Status NullStageState
	Log    []byte
}

type AddStageLogRow struct {
	ID     pgtype.UUID
	FileID pgtype.UUID
	Status NullStageState
}

// used to log the state of a file processing stage and update filestage status
func (q *Queries) AddStageLog(ctx context.Context, arg AddStageLogParams) (AddStageLogRow, error) {
	row := q.db.QueryRow(ctx, addStageLog, arg.FileID, arg.Status, arg.Log)
	var i AddStageLogRow
	err := row.Scan(&i.ID, &i.FileID, &i.Status)
	return i, err
}

const createFile = `-- name: CreateFile :one
INSERT INTO public.file (
		id,
		extension,
		lang,
		name,
		isPrivate,
    hash,
		created_at,
		updated_at
	)
VALUES (
		gen_random_uuid(),
		$1,
		$2,
		$3,
		$4,
		$5,
		NOW(),
		NOW()
	)
RETURNING id
`

type CreateFileParams struct {
	Extension pgtype.Text
	Lang      pgtype.Text
	Name      pgtype.Text
	Isprivate pgtype.Bool
	Hash      pgtype.Text
}

func (q *Queries) CreateFile(ctx context.Context, arg CreateFileParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, createFile,
		arg.Extension,
		arg.Lang,
		arg.Name,
		arg.Isprivate,
		arg.Hash,
	)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteFile = `-- name: DeleteFile :exec
DELETE FROM public.file WHERE id = $1
`

func (q *Queries) DeleteFile(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteFile, id)
	return err
}

const extrasFileCreate = `-- name: ExtrasFileCreate :one
INSERT INTO public.file_extras (
    id,
    isPrivate,
    summary,
    short_summary,
    purpose,
    created_at,
    updated_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    NOW(),
    NOW()
)
RETURNING id
`

type ExtrasFileCreateParams struct {
	ID           pgtype.UUID
	Isprivate    pgtype.Bool
	Summary      pgtype.Text
	ShortSummary pgtype.Text
	Purpose      pgtype.Text
}

func (q *Queries) ExtrasFileCreate(ctx context.Context, arg ExtrasFileCreateParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, extrasFileCreate,
		arg.ID,
		arg.Isprivate,
		arg.Summary,
		arg.ShortSummary,
		arg.Purpose,
	)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const extrasFileFetch = `-- name: ExtrasFileFetch :one
SELECT id, isprivate, mdata, created_at, updated_at
FROM public.file_metadata
WHERE id = $1
`

func (q *Queries) ExtrasFileFetch(ctx context.Context, id pgtype.UUID) (FileMetadatum, error) {
	row := q.db.QueryRow(ctx, extrasFileFetch, id)
	var i FileMetadatum
	err := row.Scan(
		&i.ID,
		&i.Isprivate,
		&i.Mdata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const extrasFileUpdate = `-- name: ExtrasFileUpdate :one
UPDATE public.file_extras
SET isPrivate = $1,
    summary = $2,
    short_summary = $3,
    purpose = $4,
    updated_at = NOW()
WHERE id = $5
RETURNING id
`

type ExtrasFileUpdateParams struct {
	Isprivate    pgtype.Bool
	Summary      pgtype.Text
	ShortSummary pgtype.Text
	Purpose      pgtype.Text
	ID           pgtype.UUID
}

func (q *Queries) ExtrasFileUpdate(ctx context.Context, arg ExtrasFileUpdateParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, extrasFileUpdate,
		arg.Isprivate,
		arg.Summary,
		arg.ShortSummary,
		arg.Purpose,
		arg.ID,
	)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const fetchMetadata = `-- name: FetchMetadata :one
SELECT id, isprivate, mdata, created_at, updated_at
FROM public.file_metadata
WHERE id = $1
`

func (q *Queries) FetchMetadata(ctx context.Context, id pgtype.UUID) (FileMetadatum, error) {
	row := q.db.QueryRow(ctx, fetchMetadata, id)
	var i FileMetadatum
	err := row.Scan(
		&i.ID,
		&i.Isprivate,
		&i.Mdata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getFileMetadata = `-- name: GetFileMetadata :one
SELECT id, isprivate, mdata, created_at, updated_at FROM public.file_metadata WHERE id = $1
`

func (q *Queries) GetFileMetadata(ctx context.Context, id pgtype.UUID) (FileMetadatum, error) {
	row := q.db.QueryRow(ctx, getFileMetadata, id)
	var i FileMetadatum
	err := row.Scan(
		&i.ID,
		&i.Isprivate,
		&i.Mdata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const insertMetadata = `-- name: InsertMetadata :one
INSERT INTO public.file_metadata (
    id,
    isPrivate,
    mdata,
    created_at,
    updated_at
)
VALUES (
    $1,
    $2,
    $3,
    NOW(),
    NOW()
)
RETURNING id
`

type InsertMetadataParams struct {
	ID        pgtype.UUID
	Isprivate pgtype.Bool
	Mdata     []byte
}

func (q *Queries) InsertMetadata(ctx context.Context, arg InsertMetadataParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, insertMetadata, arg.ID, arg.Isprivate, arg.Mdata)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const listFiles = `-- name: ListFiles :many
SELECT id, lang, name, extension, isprivate, created_at, updated_at, hash, metadata_id
FROM public.file
ORDER BY updated_at DESC
`

func (q *Queries) ListFiles(ctx context.Context) ([]File, error) {
	rows, err := q.db.Query(ctx, listFiles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []File
	for rows.Next() {
		var i File
		if err := rows.Scan(
			&i.ID,
			&i.Lang,
			&i.Name,
			&i.Extension,
			&i.Isprivate,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Hash,
			&i.MetadataID,
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

const readFile = `-- name: ReadFile :one
SELECT id, lang, name, extension, isprivate, created_at, updated_at, hash, metadata_id
FROM public.file
WHERE id = $1
`

func (q *Queries) ReadFile(ctx context.Context, id pgtype.UUID) (File, error) {
	row := q.db.QueryRow(ctx, readFile, id)
	var i File
	err := row.Scan(
		&i.ID,
		&i.Lang,
		&i.Name,
		&i.Extension,
		&i.Isprivate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Hash,
		&i.MetadataID,
	)
	return i, err
}

const updateFile = `-- name: UpdateFile :one
UPDATE public.file
SET extension = $1,
	lang = $2,
	name = $3,
	isPrivate = $4,
  hash = $5,
	updated_at = NOW()
WHERE public.file.id = $6
RETURNING id
`

type UpdateFileParams struct {
	Extension pgtype.Text
	Lang      pgtype.Text
	Name      pgtype.Text
	Isprivate pgtype.Bool
	Hash      pgtype.Text
	ID        pgtype.UUID
}

func (q *Queries) UpdateFile(ctx context.Context, arg UpdateFileParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, updateFile,
		arg.Extension,
		arg.Lang,
		arg.Name,
		arg.Isprivate,
		arg.Hash,
		arg.ID,
	)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const updateMetadata = `-- name: UpdateMetadata :one
UPDATE public.file_metadata
SET isPrivate = $1,
    mdata = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING id
`

type UpdateMetadataParams struct {
	Isprivate pgtype.Bool
	Mdata     []byte
	ID        pgtype.UUID
}

func (q *Queries) UpdateMetadata(ctx context.Context, arg UpdateMetadataParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, updateMetadata, arg.Isprivate, arg.Mdata, arg.ID)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}
