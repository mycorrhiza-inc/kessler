-- name: AssociateDocumentWithOrganization :one
INSERT INTO public.relation_documents_organizations (
		document_id,
		organization_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING *;
-- name: ListDocumentIdsByOrganization :many
SELECT document_id
FROM public.relation_documents_organizations
WHERE organization_id = $1;
-- name: ListOrganizationIdsByDocument :many
SELECT organization_id
FROM public.relation_documents_organizations
WHERE document_id = $1;
-- name: UpdateDocumentAssociatedWithOrganization :exec
UPDATE public.relation_documents_organizations
SET document_id = $1,
	organization_id = $2,
	updated_at = NOW()
WHERE id = $3
RETURNING *;
-- name: DeleteDocumentAssociatedWithOrganization :one
DELETE FROM public.relation_documents_organizations
WHERE id = $1
RETURNING *;
