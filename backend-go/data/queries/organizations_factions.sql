-- name: CreateOrganizationsInFaction :one
INSERT INTO public.relation_organizations_factions (
		faction_id,
		organization_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id;
-- name: ReadOrganizationsInFaction :one
SELECT *
FROM public.relation_organizations_factions
WHERE id = $1;
-- name: ListOrganizationsInFaction :many
SELECT *
FROM public.relation_organizations_factions
ORDER BY created_at DESC;
-- name: UpdateOrganizationsInFaction :one
UPDATE public.relation_organizations_factions
SET faction_id = $1,
	organization_id = $2,
	updated_at = NOW()
WHERE id = $3
RETURNING id;
-- name: DeleteOrganizationsInFaction :exec
DELETE FROM public.relation_organizations_factions
WHERE id = $1;