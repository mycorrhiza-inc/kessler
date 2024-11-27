-- name: OrganizationDeduplicateCascade :exec
DO
$$
BEGIN
-- Update all foreign key references from organization A to B
-- Add new UPDATE statements here when new tables with FK references are added
UPDATE
    public.relation_documents_organizations_authorship
SET
    organization_id = $2
WHERE
    organization_id = $1;

UPDATE
    public.organization_aliases
SET
    organization_id = $2
WHERE
    organization_id = $1;

-- Finally delete the source organization
DELETE FROM
    public.organization
WHERE
    id = $1;

END
$$
;

-- name: ConversationDeduplicateCascade :exec
DO
$$
BEGIN
-- Update all foreign key references from organization A to B
-- Add new UPDATE statements here when new tables with FK references are added
UPDATE
    public.docket_documents
SET
    docket_id = $2
WHERE
    docket_id = $1;

-- Finally delete the source organization
DELETE FROM
    public.docket_conversations
WHERE
    id = $1;

END
$$
;