-- name: CreateEncounter :one
INSERT INTO public.encounter (
		name,
		description,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING *;
-- name: GetEncounter :one
SELECT *
FROM public.encounter
WHERE id = $1;
-- name: ListEncounters :many
SELECT *
FROM public.encounter
ORDER BY created_at DESC;
-- name: UpdateEncounter :exec
UPDATE public.encounter
SET name = $1,
	description = $2,
	updated_at = NOW()
WHERE id = $3;
-- name: UpdateEncounterName :exec
UPDATE public.encounter
SET name = $1,
	updated_at = NOW()
WHERE id = $2;
-- name: UpdateEncounterDescription :exec
UPDATE public.encounter
SET description = $1,
	updated_at = NOW()
WHERE id = $2;
-- name: DeleteEncounter :exec
DELETE FROM public.encounter
WHERE id = $1;