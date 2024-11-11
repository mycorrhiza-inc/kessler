-- name: CreateFile :one
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
RETURNING id;
-- name: GetFile :one
SELECT *
FROM public.file
WHERE id = $1;
-- name: GetFileWithMetadata :one
SELECT *
FROM public.file
  LEFT JOIN public.file_metadata ON public.file.id = public.file_metadata.id
WHERE public.file.id = $1;
-- name: HashGetFileID :many 
SELECT id
FROM public.file
Where public.file.hash = $1;
-- name: UpdateFile :one
UPDATE public.file
SET extension = $1,
  lang = $2,
  name = $3,
  isPrivate = $4,
  hash = $5,
  updated_at = NOW()
WHERE public.file.id = $6
RETURNING id;
-- name: ReadFile :one
SELECT *
FROM public.file
WHERE id = $1;
-- name: ListFiles :many
SELECT *
FROM public.file
ORDER BY updated_at DESC;
-- name: AddStageLog :one
-- used to log the state of a file processing stage and update filestage status
INSERT INTO public.stage_log (file_id, status, log)
VALUES ($1, $2, $3)
RETURNING id,
  file_id,
  status;
-- name: DeleteFile :exec
DELETE FROM public.file
WHERE id = $1;
-- name: InsertMetadata :one
INSERT INTO public.file_metadata (
    id,
    isPrivate,
    mdata,
    created_at,
    updated_at
  )
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id;
-- name: UpdateMetadata :one
UPDATE public.file_metadata
SET isPrivate = $1,
  mdata = $2,
  updated_at = NOW()
WHERE id = $3
RETURNING id;
-- name: FetchMetadata :many
SELECT *
FROM public.file_metadata
WHERE id = $1;
-- name: ExtrasFileCreate :one
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
RETURNING id;
-- name: ExtrasFileUpdate :one
UPDATE public.file_extras
SET isPrivate = $1,
  extra_obj = $2,
  updated_at = NOW()
WHERE id = $3
RETURNING id;
-- name: ExtrasFileFetch :one
SELECT *
FROM public.file_extras
WHERE id = $1;
-- name: GetFileMetadata :one
SELECT *
FROM public.file_metadata
WHERE id = $1;
