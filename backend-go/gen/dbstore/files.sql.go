// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: files.sql

package dbstore

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createFile = `-- name: CreateFile :one
INSERT INTO public.file (
    id,
    extension,
    lang,
    name,
    isPrivate,
    hash,
    verified,
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
    $6,
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
	Verified  pgtype.Bool
}

func (q *Queries) CreateFile(ctx context.Context, arg CreateFileParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, createFile,
		arg.Extension,
		arg.Lang,
		arg.Name,
		arg.Isprivate,
		arg.Hash,
		arg.Verified,
	)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteFile = `-- name: DeleteFile :exec
DELETE FROM public.file
WHERE id = $1
`

func (q *Queries) DeleteFile(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteFile, id)
	return err
}

const extrasFileCreate = `-- name: ExtrasFileCreate :one
INSERT INTO public.file_extras (
    id,
    isPrivate,
    extra_obj,
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

type ExtrasFileCreateParams struct {
	ID        pgtype.UUID
	Isprivate pgtype.Bool
	ExtraObj  []byte
}

func (q *Queries) ExtrasFileCreate(ctx context.Context, arg ExtrasFileCreateParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, extrasFileCreate, arg.ID, arg.Isprivate, arg.ExtraObj)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const extrasFileFetch = `-- name: ExtrasFileFetch :one
SELECT id, isprivate, created_at, updated_at, extra_obj
FROM public.file_extras
WHERE id = $1
`

func (q *Queries) ExtrasFileFetch(ctx context.Context, id pgtype.UUID) (FileExtra, error) {
	row := q.db.QueryRow(ctx, extrasFileFetch, id)
	var i FileExtra
	err := row.Scan(
		&i.ID,
		&i.Isprivate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExtraObj,
	)
	return i, err
}

const extrasFileUpdate = `-- name: ExtrasFileUpdate :one
UPDATE public.file_extras
SET isPrivate = $1,
  extra_obj = $2,
  updated_at = NOW()
WHERE id = $3
RETURNING id
`

type ExtrasFileUpdateParams struct {
	Isprivate pgtype.Bool
	ExtraObj  []byte
	ID        pgtype.UUID
}

func (q *Queries) ExtrasFileUpdate(ctx context.Context, arg ExtrasFileUpdateParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, extrasFileUpdate, arg.Isprivate, arg.ExtraObj, arg.ID)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const fetchMetadata = `-- name: FetchMetadata :many
SELECT id, isprivate, mdata, created_at, updated_at
FROM public.file_metadata
WHERE id = $1
`

func (q *Queries) FetchMetadata(ctx context.Context, id pgtype.UUID) ([]FileMetadatum, error) {
	rows, err := q.db.Query(ctx, fetchMetadata, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FileMetadatum
	for rows.Next() {
		var i FileMetadatum
		if err := rows.Scan(
			&i.ID,
			&i.Isprivate,
			&i.Mdata,
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

const fileVerifiedUpdate = `-- name: FileVerifiedUpdate :one
UPDATE public.file
SET verified = $1,
  updated_at = NOW()
WHERE public.file.id = $2
RETURNING id
`

type FileVerifiedUpdateParams struct {
	Verified pgtype.Bool
	ID       pgtype.UUID
}

func (q *Queries) FileVerifiedUpdate(ctx context.Context, arg FileVerifiedUpdateParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, fileVerifiedUpdate, arg.Verified, arg.ID)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const getFile = `-- name: GetFile :one
SELECT id, lang, name, extension, isprivate, created_at, updated_at, hash, verified
FROM public.file
WHERE id = $1
`

func (q *Queries) GetFile(ctx context.Context, id pgtype.UUID) (File, error) {
	row := q.db.QueryRow(ctx, getFile, id)
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
		&i.Verified,
	)
	return i, err
}

const getFileWithMetadata = `-- name: GetFileWithMetadata :one
SELECT file.id, lang, name, extension, file.isprivate, file.created_at, file.updated_at, hash, verified, file_metadata.id, file_metadata.isprivate, mdata, file_metadata.created_at, file_metadata.updated_at
FROM public.file
  LEFT JOIN public.file_metadata ON public.file.id = public.file_metadata.id
WHERE public.file.id = $1
`

type GetFileWithMetadataRow struct {
	ID          pgtype.UUID
	Lang        pgtype.Text
	Name        pgtype.Text
	Extension   pgtype.Text
	Isprivate   pgtype.Bool
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	Hash        pgtype.Text
	Verified    pgtype.Bool
	ID_2        pgtype.UUID
	Isprivate_2 pgtype.Bool
	Mdata       []byte
	CreatedAt_2 pgtype.Timestamptz
	UpdatedAt_2 pgtype.Timestamptz
}

func (q *Queries) GetFileWithMetadata(ctx context.Context, id pgtype.UUID) (GetFileWithMetadataRow, error) {
	row := q.db.QueryRow(ctx, getFileWithMetadata, id)
	var i GetFileWithMetadataRow
	err := row.Scan(
		&i.ID,
		&i.Lang,
		&i.Name,
		&i.Extension,
		&i.Isprivate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Hash,
		&i.Verified,
		&i.ID_2,
		&i.Isprivate_2,
		&i.Mdata,
		&i.CreatedAt_2,
		&i.UpdatedAt_2,
	)
	return i, err
}

const hashGetFileID = `-- name: HashGetFileID :many
SELECT id
FROM public.file
Where public.file.hash = $1
`

func (q *Queries) HashGetFileID(ctx context.Context, hash pgtype.Text) ([]pgtype.UUID, error) {
	rows, err := q.db.Query(ctx, hashGetFileID, hash)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.UUID
	for rows.Next() {
		var id pgtype.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertMetadata = `-- name: InsertMetadata :one
INSERT INTO public.file_metadata (
    id,
    isPrivate,
    mdata,
    created_at,
    updated_at
  )
VALUES ($1, $2, $3, NOW(), NOW())
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
SELECT id, lang, name, extension, isprivate, created_at, updated_at, hash, verified
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
			&i.Verified,
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
SELECT id, lang, name, extension, isprivate, created_at, updated_at, hash, verified
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
		&i.Verified,
	)
	return i, err
}

const stageLogAdd = `-- name: StageLogAdd :one
INSERT INTO public.stage_log (file_id, status, log)
VALUES ($1, $2, $3)
RETURNING id,
  file_id,
  status
`

type StageLogAddParams struct {
	FileID pgtype.UUID
	Status NullStageState
	Log    []byte
}

type StageLogAddRow struct {
	ID     pgtype.UUID
	FileID pgtype.UUID
	Status NullStageState
}

// used to log the state of a file processing stage and update filestage status
func (q *Queries) StageLogAdd(ctx context.Context, arg StageLogAddParams) (StageLogAddRow, error) {
	row := q.db.QueryRow(ctx, stageLogAdd, arg.FileID, arg.Status, arg.Log)
	var i StageLogAddRow
	err := row.Scan(&i.ID, &i.FileID, &i.Status)
	return i, err
}

const stageLogFileGetLatest = `-- name: StageLogFileGetLatest :one
SELECT id, status, log, created_at, file_id
FROM public.stage_log
WHERE file_id = $1
ORDER BY created_at DESC
LIMIT 1
`

func (q *Queries) StageLogFileGetLatest(ctx context.Context, fileID pgtype.UUID) (StageLog, error) {
	row := q.db.QueryRow(ctx, stageLogFileGetLatest, fileID)
	var i StageLog
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Log,
		&i.CreatedAt,
		&i.FileID,
	)
	return i, err
}

const updateFile = `-- name: UpdateFile :exec
UPDATE public.file
SET extension = $1,
  lang = $2,
  name = $3,
  isPrivate = $4,
  hash = $5,
  verified = $6,
  updated_at = NOW()
WHERE public.file.id = $7
`

type UpdateFileParams struct {
	Extension pgtype.Text
	Lang      pgtype.Text
	Name      pgtype.Text
	Isprivate pgtype.Bool
	Hash      pgtype.Text
	Verified  pgtype.Bool
	ID        pgtype.UUID
}

func (q *Queries) UpdateFile(ctx context.Context, arg UpdateFileParams) error {
	_, err := q.db.Exec(ctx, updateFile,
		arg.Extension,
		arg.Lang,
		arg.Name,
		arg.Isprivate,
		arg.Hash,
		arg.Verified,
		arg.ID,
	)
	return err
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
