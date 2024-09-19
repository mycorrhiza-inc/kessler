-- name: CreateOrganizationsInFaction :one
INSERT INTO public.relation_organizations_in_faction (
		faction_id,
		organization_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id;
-- name: ReadOrganizationsInFaction :one
SELECT *
FROM public.relation_organizations_in_faction
WHERE id = $1;
-- name: ListOrganizationsInFaction :many
SELECT *
FROM public.relation_organizations_in_faction
ORDER BY created_at DESC;
-- name: UpdateOrganizationsInFaction :one
UPDATE public.relation_organizations_in_faction
SET faction_id = $1,
	organization_id = $2,
	updated_at = NOW()
WHERE id = $3
RETURNING id;
-- name: DeleteOrganizationsInFaction :one
DELETE FROM public.relation_organizations_in_faction
WHERE id = $1;