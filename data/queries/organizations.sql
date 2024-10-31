-- name: AuthorshipDocumentOrganizationInsert: one
INSERT INTO public.relation_documents_organizations_authorship (
		document_id,
		organization_id,
    is_primary_author,
		created_at,
		updated_at
	)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id;
-- name: AuthorshipDocumentDeleteAll:exec
DELETE FROM public.relation_documents_organizations_authorship
WHERE document_id = $1;
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
-- name: DeleteOrganization :exec
DELETE FROM public.organization
WHERE id = $1;
