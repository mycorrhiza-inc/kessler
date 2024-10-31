-- name: CreateOrganization :one
INSERT INTO public.organization (
		name,
		description,
    is_person,
		created_at,
		updated_at
	)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id;
-- name: OrganizationFetchByNameMatch : many
SELECT *
FROM public.organization
WHERE name = $1;
-- name: ReadOrganization :one
SELECT *
FROM public.organization
WHERE id = $1;
-- name: ListOrganizations :many
SELECT *
FROM public.organization
ORDER BY created_at DESC;
-- name: UpdateOrganization :one
UPDATE public.organization
SET name = $1,
	description = $2,
  is_person = $3
	updated_at = NOW()
WHERE id = $4
RETURNING id;
-- name: UpdateOrganizationName :one
UPDATE public.organization
SET name = $1,
	updated_at = NOW()
WHERE id = $2
RETURNING id;
-- name: UpdateOrganizationDescription :one
UPDATE public.organization
SET description = $1,
	updated_at = NOW()
WHERE id = $2
RETURNING id;
-- name: DeleteOrganization :exec
DELETE FROM public.organization
WHERE id = $1;
