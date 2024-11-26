-- name: AuthorshipDocumentOrganizationInsert :one
INSERT INTO
    public.relation_documents_organizations_authorship (
        document_id,
        organization_id,
        is_primary_author,
        created_at,
        updated_at
    )
VALUES
    ($1, $2, $3, NOW(), NOW())
RETURNING
    id;

-- name: AuthorshipOrganizationListDocuments :many
SELECT
    *
FROM
    public.relation_documents_organizations_authorship
WHERE
    organization_id = $1;

-- name: AuthorshipDocumentDeleteAll :exec
DELETE FROM
    public.relation_documents_organizations_authorship
WHERE
    document_id = $1;

-- name: CreateOrganization :one
INSERT INTO
    public.organization (
        name,
        description,
        is_person,
        created_at,
        updated_at
    )
VALUES
    ($1, $2, $3, NOW(), NOW())
RETURNING
    id;

-- name: OrganizationFetchByNameMatch :many
SELECT
    *
FROM
    public.organization
WHERE
    name = $1;

-- name: OrganizationRead :one
SELECT
    *
FROM
    public.organization
WHERE
    id = $1;

-- name: OrganizationList :many
SELECT
    *
FROM
    public.organization
ORDER BY
    created_at DESC;

-- name: OrganizationUpdate :one
UPDATE
    public.organization
SET
    name = $1,
    description = $2,
    is_person = $3,
    updated_at = NOW()
WHERE
    id = $4
RETURNING
    id;

-- name: OrganizationDelete :exec
DELETE FROM
    public.organization
WHERE
    id = $1;


-- name: OrganizationgGetConversationsAuthoredIn :many
SELECT
    public.organization.id AS organization_id,
    public.organization.name AS organization_name,
    public.relation_documents_organizations_authorship.document_id,
    public.docket_conversations.docket_id AS docket_id,
    public.docket_conversations.id AS conversation_uuid
FROM
    public.organization
    LEFT JOIN public.relation_documents_organizations_authorship ON public.organization.id = public.relation_documents_organizations_authorship.organization_id
    LEFT JOIN public.docket_documents ON public.relation_documents_organizations_authorship.document_id = public.docket_documents.file_id
    LEFT JOIN public.docket_conversations ON public.docket_documents.docket_id = public.docket_conversations.id
WHERE
    public.organization.id = $1;