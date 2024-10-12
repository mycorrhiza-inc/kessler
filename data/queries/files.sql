-- name: CreateFile :one
INSERT INTO public.file (
		url,
		doctype,
		lang,
		name,
		source,
		hash,
		mdata,
		stage,
		summary,
		short_summary,
		created_at,
		updated_at
	)
VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		NOW(),
		NOW()
	)
RETURNING *;
-- name: ReadFile :one
SELECT *
FROM public.file
WHERE id = $1;
-- name: ListFiles :many
SELECT *
FROM public.file
ORDER BY created_at DESC;
-- name: UpdateFile :one
UPDATE public.file
SET url = $1,
	doctype = $2,
	lang = $3,
	name = $4,
	source = $5,
	hash = $6,
	mdata = $7,
	stage = $8,
	summary = $9,
	short_summary = $10,
	updated_at = NOW()
WHERE id = $11
RETURNING *;
-- name: DeleteFile :exec
DELETE FROM public.file
WHERE id = $1;
