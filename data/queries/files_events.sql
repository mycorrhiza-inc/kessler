-- name: CreateFilesAssociatedWithEvent :one
INSERT INTO public.relation_files_events (
		file_id,
		event_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING *;
-- name: ReadFilesAssociatedWithEvent :one
SELECT *
FROM public.relation_files_events
WHERE id = $1;
-- name: ListFilesAssociatedWithEvent :many
SELECT *
FROM public.relation_files_events
ORDER BY created_at DESC;
-- name: UpdateFilesAssociatedWithEvent :one
UPDATE public.relation_files_events
SET file_id = $1,
	event_id = $2,
	updated_at = NOW()
WHERE id = $3
RETURNING *;
-- name: DeleteFilesAssociatedWithEvent :exec
DELETE FROM public.relation_files_events
WHERE id = $1;
