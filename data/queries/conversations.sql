
-- name: DocketDocumentInsert :one
INSERT INTO public.docket_documents (
		docket_id,
		file_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING docket_id;

-- name: DocketDocumentUpdate :one
UPDATE public.docket_documents 
SET
		docket_id = $1,
		updated_at = NOW()
WHERE file_id = $2
RETURNING file_id;

-- name: DocketDocumentDeleteAll :exec
DELETE FROM public.docket_documents
WHERE docket_id = $1;

-- name: DocketConversationCreate :one
INSERT INTO public.docket_conversations (
		docket_id,
		state,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id;

-- name: DocketConversationFetchByDocketIdMatch :many
SELECT *
FROM public.docket_conversations
WHERE docket_id = $1;

-- name: DocketConversationRead :one
SELECT *
FROM public.docket_conversations
WHERE id = $1;


-- name: DocketConversationUpdate :one
UPDATE public.docket_conversations
SET docket_id = $1,
	state = $2,
	updated_at = NOW()
WHERE id = $3
RETURNING id;

-- name: DocketConversationDelete :exec
DELETE FROM public.docket_conversations
WHERE id = $1;