-- name: OrganizationDeduplicateCascade :exec
WITH update_authorship AS (
    UPDATE
        public.relation_documents_organizations_authorship
    SET
        organization_id = $2
    WHERE
        organization_id = $1
    RETURNING
        1
),
update_aliases AS (
    UPDATE
        public.organization_aliases
    SET
        organization_id = $2
    WHERE
        organization_id = $1
    RETURNING
        1
)
DELETE FROM
    public.organization
WHERE
    public.organization.id = $1;

-- name: ConversationDeduplicateCascade :exec
WITH update_documents AS (
    UPDATE
        public.docket_documents
    SET
        conversation_uuid = $2
    WHERE
        conversation_uuid = $1
    RETURNING
        1
)
DELETE FROM
    public.docket_conversations
WHERE
    id = $1;

-- name: FileCheckForDuplicates :many
SELECT
    public.file.id,
    public.file.name,
    public.file.extension,
    public.file.lang,
    public.file.verified,
    public.file.hash,
    public.file.created_at,
    public.file.updated_at,
    public.file.date_published,
    public.file_metadata.mdata,
    public.file_extras.extra_obj,
    public.docket_documents.conversation_uuid,
    public.docket_conversations.docket_gov_id
FROM
    public.file
    LEFT JOIN public.file_metadata ON public.file.id = public.file_metadata.id
    LEFT JOIN public.file_extras ON public.file.id = public.file_extras.id
    LEFT JOIN public.docket_documents ON public.file.id = public.docket_documents.file_id
    LEFT JOIN public.docket_conversations ON public.docket_documents.conversation_uuid = public.docket_conversations.id
WHERE
    public.file.name = $1
    AND public.file.extension = $2
    AND public.docket_conversations.docket_gov_id = $3;