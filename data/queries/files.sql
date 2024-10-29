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
-- name: ReadFile :one
SELECT *
FROM public.file
WHERE id = $1;
-- name: ListFiles :many
SELECT *
FROM public.file
ORDER BY updated_at DESC;
-- TODO: Update these to work with the deletion of filestage
-- -- name: ListUnprocessedFiles :many
-- -- get all files that are not in the completed stage
-- -- use a left join to get the stage status
-- SELECT *
-- FROM public.file
-- 	LEFT JOIN public.stage_log sl ON f.stage_id = sl.stage_id
-- WHERE sl.status != 'completed'
-- ORDER BY sl.created_at DESC;
--
-- -- name: ListUnprocessedFilesPagnated :many
-- SELECT *
-- FROM public.file
-- 	LEFT JOIN public.stage_log sl ON f.stage_id = sl.stage_id
-- WHERE sl.status != 'completed'
-- ORDER BY sl.created_at DESC
-- LIMIT $1 OFFSET $2;
-- name: AddStageLog :one
-- used to log the state of a file processing stage and update filestage status
WITH inserted_log AS (
	INSERT INTO public.stage_log (
			file_id,
			status,
			log
		)
	VALUES ($1, $2, $3)
	RETURNING id,
		file_id,
		status
)
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
VALUES (
    $1,
    $2,
    $3,
    NOW(),
    NOW()
)
RETURNING id;
-- name: UpdateMetadata :one
UPDATE public.file_metadata
SET isPrivate = $1,
    mdata = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING id;
-- name: FetchMetadata :one
SELECT *
FROM public.file_metadata
WHERE id = $1;
