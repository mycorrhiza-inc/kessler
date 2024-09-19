-- name: AssociateDocumentWithOrganization :one
INSERT INTO public.relation_document_organization (
		document_id,
		organization_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id;
-- name: ListDocumentIdsByOrganization :many
SELECT document_id
FROM public.relation_document_organization
WHERE organization_id = $1;
-- name: ListOrganizationIdsByDocument :many
SELECT organization_id
FROM public.relation_document_organization
WHERE document_id = $1;
-- name: UpdateDocumentAssociatedWithOrganization :one
UPDATE public.relation_document_organization
SET document_id = $1,
	organization_id = $2,
	updated_at = NOW()
WHERE id = $3
RETURNING id;
-- name: DeleteDocumentAssociatedWithOrganization :one
DELETE FROM public.relation_document_organization
WHERE id = $1;