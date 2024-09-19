-- name: CreateOrganizationsAssociatedWithEvent :one
INSERT INTO public.relation_organization_event (
		organization_id,
		event_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id;
-- name: ReadOrganizationsAssociatedWithEvent :one
SELECT *
FROM public.relation_organization_event
WHERE id = $1;
-- name: ListOrganizationsAssociatedWithEvent :many
SELECT *
FROM public.relation_organization_event
ORDER BY created_at DESC;
-- name: UpdateOrganizationsAssociatedWithEvent :one
UPDATE public.relation_organization_event
SET organization_id = $1,
	event_id = $2,
	updated_at = NOW()
WHERE id = $3
RETURNING id;
-- name: DeleteOrganizationsAssociatedWithEvent :one
DELETE FROM public.relation_organization_event
WHERE id = $1;