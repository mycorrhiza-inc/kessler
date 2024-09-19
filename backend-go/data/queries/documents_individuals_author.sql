-- name: CreateDocumentAuthoredByIndividual :one
INSERT INTO public.relation_document_individual_author (
		document_id,
		individual_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id;
-- name: GetDocumentAuthors :one
SELECT individual_id
FROM public.relation_document_individual_author
WHERE document_id = $1;
-- name: ListDocumentsAuthoredByIndividual :many
SELECT *
FROM public.relation_document_individual_author
WHERE individual_id = $1
ORDER BY created_at DESC;
-- name: UpdateDocumentAuthoredByIndividual :one
UPDATE public.relation_document_individual_author
SET document_id = $1,
	individual_id = $2,
	sa_orm_sentinel = $3,
	updated_at = NOW()
WHERE id = $4
RETURNING id;
-- name: DeleteDocumentAuthoredByIndividual :one
DELETE FROM public.relation_document_individual_author
WHERE id = $1;