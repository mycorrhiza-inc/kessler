-- name: CreateIndividualsInFaction :one
INSERT INTO public.relation_individual_faction (
		faction_id,
		individual_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id;
-- name: ReadIndividualsInFaction :one
SELECT *
FROM public.relation_individual_faction
WHERE id = $1;
-- name: ListIndividualsInFaction :many
SELECT *
FROM public.relation_individual_faction
ORDER BY created_at DESC;
-- name: UpdateIndividualsInFaction :one
UPDATE public.relation_individual_faction
SET faction_id = $1,
	individual_id = $2,
	sa_orm_sentinel = $3,
	updated_at = NOW()
WHERE id = $4
RETURNING id;
-- name: DeleteIndividualsInFaction :one
DELETE FROM public.relation_individual_faction
WHERE id = $1;