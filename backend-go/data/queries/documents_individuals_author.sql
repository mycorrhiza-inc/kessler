-- name: LinkDocumentToIndividual :one
INSERT INTO public.relation_documents_individuals_author (
		document_id,
		individual_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING *;
-- name: GetDocumentAuthors :many
SELECT individual_id
FROM public.relation_documents_individuals_author
WHERE document_id = $1;
-- name: ListDocumentsAuthoredByIndividual :many
SELECT *
FROM public.relation_documents_individuals_author
WHERE individual_id = $1
ORDER BY created_at DESC;
-- name: UnlinkDocumentFromIndividual :exec
DELETE FROM public.relation_documents_individuals_author
WHERE document_id = $1
	AND individual_id = $2;