-- name: CreateIndividualsCurrentlyAssociatedWithOrganization :one
INSERT INTO public.relation_individual_organization (individual_id, organization_id, sa_orm_sentinel, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id;

-- name: ReadIndividualsCurrentlyAssociatedWithOrganization :one
SELECT * FROM public.relation_individual_organization WHERE id = $1;

-- name: ListIndividualsCurrentlyAssociatedWithOrganization :many
SELECT * FROM public.relation_individual_organization ORDER BY created_at DESC;

-- name: UpdateIndividualsCurrentlyAssociatedWithOrganization :one
UPDATE public.relation_individual_organization
SET individual_id = $1, organization_id = $2, sa_orm_sentinel = $3, updated_at = NOW()
WHERE id = $4 RETURNING id;

-- name: DeleteIndividualsCurrentlyAssociatedWithOrganization :one
DELETE FROM public.relation_individual_organization WHERE id = $1;