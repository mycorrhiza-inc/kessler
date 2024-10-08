-- name: CreateFileTextSource :one
INSERT INTO public.file_text_source (
		file_id,
		is_original_text,
		language,
		text,
		created_at,
		updated_at
	)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id;
-- name: ListTextsOfFile :many
SELECT *
FROM public.file_text_source
WHERE file_id = $1
-- name: ListTextsOfFileWithLang :many
SELECT *
FROM public.file_text_source
WHERE file_id = $1 and language = $2
-- name: ListTextsOfFileOriginal
SELECT *
FROM public.file_text_source
WHERE file_id = $1 and is_original_text = true
-- name: DeleteFileTexts :exec
DELETE FROM public.file_text_source
WHERE file_id = $1;
