-- name: CreateIndividualsInFaction :one
INSERT INTO public.relation_individuals_factions (
		faction_id,
		individual_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id;
-- name: ReadIndividualsInFaction :one
SELECT *
FROM public.relation_individuals_factions
WHERE id = $1;
-- name: ListIndividualsInFaction :many
SELECT *
FROM public.relation_individuals_factions
ORDER BY created_at DESC;
-- name: UpdateIndividualsInFaction :one
UPDATE public.relation_individuals_factions
SET faction_id = $1,
	individual_id = $2,
	updated_at = NOW()
WHERE id = $3
RETURNING id;
-- name: DeleteIndividualsInFaction :exec
DELETE FROM public.relation_individuals_factions
WHERE id = $1;