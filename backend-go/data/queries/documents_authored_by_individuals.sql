-- name: CreateDocumentAuthoredByIndividual :one
INSERT INTO public.relation_document_authored_by_individual (
		document_id,
		individual_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id;
-- name: GetDocumentAuthors :one
SELECT individual_id
FROM public.relation_document_authored_by_individual
WHERE document_id = $1;
-- name: ListDocumentsAuthoredByIndividual :many
SELECT *
FROM public.relation_document_authored_by_individual
WHERE individual_id = $1
ORDER BY created_at DESC;
-- name: UpdateDocumentAuthoredByIndividual :one
UPDATE public.relation_document_authored_by_individual
SET document_id = $1,
	individual_id = $2,
	sa_orm_sentinel = $3,
	updated_at = NOW()
WHERE id = $4
RETURNING id;
-- name: DeleteDocumentAuthoredByIndividual :one
DELETE FROM public.relation_document_authored_by_individual
WHERE id = $1;