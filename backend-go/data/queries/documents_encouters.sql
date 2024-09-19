-- name: AddDocumentToEncounter :one
INSERT INTO public.relation_documents_in_encounter (
		document_id,
		encounter_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id;
-- name: ListDocumentsInEncounter :many
SELECT document_id
FROM public.relation_documents_in_encounter
WHERE encounter_id = $1;
-- name: DeleteDocumentsInEncounter :one
DELETE FROM public.relation_documents_in_encounter
WHERE id = $1;