-- name: CreateDocketConversation :one
INSERT INTO public.docket_conversations (docket_id, created_at, updated_at)
SELECT $1,
	NOW(),
	NOW()
WHERE NOT EXISTS (
		SELECT 1
		FROM public.docket_conversations
		WHERE docket_id = $1
	)
RETURNING id;
-- name: GetDocketConversation :one
SELECT *
FROM public.docket_conversations
WHERE docket_id = $1;
-- name: AddFileToDocket :exec
INSERT INTO public.docket_documents (docket_id, file_id, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW());
-- name: IsDocumentInDocket :one
SELECT EXISTS (
		SELECT 1
		FROM public.docket_documents
		WHERE docket_id = $1
			AND file_id = $2
	);