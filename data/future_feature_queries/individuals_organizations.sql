-- name: CreateIndividualsCurrentlyAssociatedWithOrganization :one
INSERT INTO public.relation_individuals_organizations (
		individual_id,
		organization_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id;
-- name: ListIndividualsCurrentlyAssociatedWithOrganization :many
SELECT *
FROM public.relation_individuals_organizations
WHERE organization_id = $1;
-- name: UpdateIndividualsCurrentlyAssociatedWithOrganization :one
UPDATE public.relation_individuals_organizations
SET individual_id = $1,
	organization_id = $2,
	updated_at = NOW()
WHERE id = $3
RETURNING id;
-- name: DeleteIndividualsCurrentlyAssociatedWithOrganization :exec
DELETE FROM public.relation_individuals_organizations
WHERE id = $1;
