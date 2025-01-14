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
        docket_gov_id = $2
    WHERE
        docket_gov_id = $1
    RETURNING
        1
)
DELETE FROM
    public.docket_conversations
WHERE
    id = $1;
