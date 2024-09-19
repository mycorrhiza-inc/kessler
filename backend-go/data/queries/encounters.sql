-- name: CreateEncounter :one
INSERT INTO public.encounter (
		name,
		description,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id;
-- name: ReadEncounter :one
SELECT *
FROM public.encounter
WHERE id = $1;
-- name: ListEncounters :many
SELECT *
FROM public.encounter
ORDER BY created_at DESC;
-- name: UpdateEncounter :one
UPDATE public.encounter
SET name = $1,
	description = $2,
	updated_at = NOW()
WHERE id = $3
RETURNING id;
-- name: UpdateEncounterName :one
UPDATE public.encounter
SET name = $1,
	updated_at = NOW()
WHERE id = $2
RETURNING id;
-- name: UpdateEncounterDescription :one
UPDATE public.encounter
SET description = $1,
	updated_at = NOW()
WHERE id = $2
RETURNING id;
-- name: DeleteEncounter :one
DELETE FROM public.encounter
WHERE id = $1;