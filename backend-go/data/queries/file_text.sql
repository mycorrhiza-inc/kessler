-- name: CreateFileTextSource :one
INSERT INTO public.file_text_source (
		file_id,
		is_original_text,
		language,
		text,
		created_at,
		updated_at
	)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id;
-- name: GetFileText :one
SELECT text
FROM public.file_text_source
WHERE file_id = $1;
-- name: GetFileLanguage :one 
SELECT language
FROM public.file_text_source
WHERE file_id = $1;
-- name: ListTextsOfFile :many
SELECT *
FROM public.file_text_source
WHERE file_id = $1
ORDER BY created_at DESC;
-- name: UpdateFileTextSource :one
UPDATE public.file_text_source
SET text = $1,
	updated_at = NOW()
WHERE id = $2
RETURNING id;
-- name: UpdateFileTextLanguage :one
UPDATE public.file_text_source
SET language = $1,
	updated_at = NOW()
WHERE file_id = $2;
-- name: UpdateIsNotOriginalText :one
UPDATE public.file_text_source
SET is_original_text = $1,
	updated_at = NOW()
WHERE id = $2;
RETURNING id;
-- name: DeleteFileTexts :one
DELETE FROM public.file_text_source
WHERE file_id = $1;