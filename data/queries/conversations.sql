-- type DocketConversation struct {
-- 	ID            uuid.UUID
-- 	DocketGovID   Not null varchar 
-- 	State         Not null varchar 
-- 	CreatedAt     pgtype.Timestamp
-- 	UpdatedAt     pgtype.Timestamp
-- 	DeletedAt     pgtype.Timestamp
-- 	Name          Not null varchar 
-- 	Description   Not null varchar 
-- 	MatterType    Not null varchar 
-- 	IndustryType  Not null varchar 
-- 	Metadata      Not null varchar 
-- 	Extra         JSONB
-- 	DatePublished pgtype.Timestamptz
-- }
-- name: DocketDocumentInsert :one
INSERT INTO
    public.docket_documents (
        conversation_uuid,
        file_id,
        created_at,
        updated_at
    )
VALUES
    ($1, $2, NOW(), NOW())
RETURNING
    conversation_uuid;

-- name: DocketDocumentUpdate :one
UPDATE
    public.docket_documents
SET
    conversation_uuid = $1,
    updated_at = NOW()
WHERE
    file_id = $2
RETURNING
    file_id;

-- name: DocketDocumentDeleteAll :exec
DELETE FROM
    public.docket_documents
WHERE
    conversation_uuid = $1;

-- name: DocketConversationCreate :one
INSERT INTO
    public.docket_conversations (
        docket_gov_id,
        name,
        description,
        state,
        created_at,
        updated_at
    )
VALUES
    ($1, $2, $3, $4, NOW(), NOW())
RETURNING
    id;

-- name: DocketConversationFetchByDocketIdMatch :many
SELECT
    *
FROM
    public.docket_conversations
WHERE
    docket_gov_id = $1;

-- name: DocketConversationRead :one
SELECT
    *
FROM
    public.docket_conversations
WHERE
    id = $1;

-- name: DocketConversationList :many
SELECT
    *
FROM
    public.docket_conversations
ORDER BY
    created_at DESC;

-- name: DocketConversationUpdate :one
UPDATE
    public.docket_conversations
SET
    docket_gov_id = $1,
    state = $2,
    name = $3,
    description = $4,
    updated_at = NOW()
WHERE
    id = $5
RETURNING
    id;

-- name: DocketConversationDelete :exec
DELETE FROM
    public.docket_conversations
WHERE
    id = $1;

-- name: ConversationSemiCompleteInfoList :many
SELECT
    dc.id,
    dc.docket_gov_id,
    COUNT(dd.file_id) AS document_count,
    dc."name",
    dc.description,
    dc.created_at,
    dc.updated_at
FROM
    public.docket_conversations dc
    LEFT JOIN public.docket_documents dd ON dd.docket_gov_id = dc.id
GROUP BY
    dc.id
ORDER BY
    document_count DESC;

-- name: ConversationSemiCompleteInfoListPaginated :many
SELECT
    dc.id,
    dc.docket_gov_id,
    COUNT(dd.file_id) AS document_count,
    dc."name",
    dc.state,
    dc.description,
    dc.created_at,
    dc.updated_at
FROM
    public.docket_conversations dc
    LEFT JOIN public.docket_documents dd ON dd.docket_gov_id = dc.id
GROUP BY
    dc.id
ORDER BY
    document_count DESC
LIMIT
    $1 OFFSET $2;
