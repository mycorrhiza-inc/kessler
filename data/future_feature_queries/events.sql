-- name: CreateEvent :one
INSERT INTO public.event (
		date,
		name,
		description,
		created_at,
		updated_at
	)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING *;
-- name: ReadEvent :one
SELECT *
FROM public.event
WHERE id = $1;
-- name: ListEvents :many
SELECT *
FROM public.event
ORDER BY date DESC;
-- name: UpdateEvent :one
UPDATE public.event
SET date = $1,
	name = $2,
	description = $3,
	updated_at = NOW()
WHERE id = $4
RETURNING id;
-- name: UpdateEventDate :one
UPDATE public.event
SET date = $1,
	updated_at = NOW()
WHERE id = $2
RETURNING id;
-- name: UpdateEventName :one
UPDATE public.event
SET name = $1,
	updated_at = NOW()
WHERE id = $2
RETURNING id;
-- name: DeleteEvent :exec
DELETE FROM public.event
WHERE id = $1;