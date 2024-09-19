-- name: CreateFilesAssociatedWithEvent :one
INSERT INTO public.relation_file_event (file_id, event_id, sa_orm_sentinel, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id;

-- name: ReadFilesAssociatedWithEvent :one
SELECT * FROM public.relation_file_event WHERE id = $1;

-- name: ListFilesAssociatedWithEvent :many
SELECT * FROM public.relation_file_event ORDER BY created_at DESC;

-- name: UpdateFilesAssociatedWithEvent :one
UPDATE public.relation_file_event
SET file_id = $1, event_id = $2, sa_orm_sentinel = $3, updated_at = NOW()
WHERE id = $4 RETURNING id;

-- name: DeleteFilesAssociatedWithEvent :one
DELETE FROM public.relation_file_event WHERE id = $1;