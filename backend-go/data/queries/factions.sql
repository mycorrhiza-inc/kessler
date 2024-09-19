-- name: CreateFaction :one
INSERT INTO public.faction (
		name,
		description,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id;
-- name: ReadFaction :one
SELECT *
FROM public.faction
WHERE id = $1;
-- name: ListFactions :many
SELECT *
FROM public.faction
ORDER BY created_at DESC;
-- name: UpdateFaction :one
UPDATE public.faction
SET name = $1,
	description = $2,
	updated_at = NOW()
WHERE id = $3
RETURNING id;
-- name: UpdateFactionName :one
UPDATE public.faction
SET name = $1,
	updated_at = NOW()
WHERE id = $2
RETURNING id;
-- name: UpdateFactionDescription :one
UPDATE public.faction
SET description = $1,
	updated_at = NOW()
WHERE id = $2
RETURNING id;
-- name: DeleteFaction :one
DELETE FROM public.faction
WHERE id = $1;