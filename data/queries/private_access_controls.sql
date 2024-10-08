-- name: CreatePrivateAccessControl :one
INSERT INTO public.private_access_controls (
		operator_id,
		object_id,
    object_table,
		created_at,
		updated_at
	)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id;
-- name: ListAcessesForOperator :one
SELECT *
FROM public.private_access_controls
WHERE operator_id = $1;
-- name: ListIndividualsInFaction :many
SELECT *
FROM public.private_access_controls
ORDER BY created_at DESC;
-- name: UpdateIndividualsInFaction :one
UPDATE public.relation_individuals_factions
SET faction_id = $1,
	individual_id = $2,
	updated_at = NOW()
WHERE id = $3
RETURNING id;
-- name: DeleteAccessControl :exec
DELETE FROM public.private_access_controls
WHERE id = $1;
