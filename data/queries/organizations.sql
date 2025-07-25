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
WITH new_org AS (
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
        id
)
INSERT INTO
    public.organization_aliases (
        organization_id,
        organization_alias,
        created_at,
        updated_at
    )
SELECT
    id,
    $1,
    NOW(),
    NOW()
FROM
    new_org
RETURNING
    (
        SELECT
            id
        FROM
            new_org
    ) AS id;

-- name: OrganizationFetchByAliasMatchSingle :one
SELECT
    public.organization_aliases.organization_alias AS alias,
    public.organization.id AS id,
    public.organization.name AS name,
    public.organization.description AS description,
    public.organization.is_person AS is_person
FROM
    public.organization_aliases
    LEFT JOIN public.organization ON public.organization.id = public.organization_aliases.organization_id
WHERE
    public.organization_aliases.organization_alias = $1;

-- name: OrganizationFetchByAliasMatchAll :many
SELECT
    public.organization_aliases.organization_alias AS alias,
    public.organization.id AS id,
    public.organization.name AS name,
    public.organization.description AS description,
    public.organization.is_person AS is_person
FROM
    public.organization_aliases
    LEFT JOIN public.organization ON public.organization.id = public.organization_aliases.organization_id
WHERE
    public.organization_aliases.organization_alias = $1;

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

-- name: AliasOrganizationCreate :one
INSERT INTO
    public.organization_aliases (
        organization_id,
        organization_alias,
        created_at,
        updated_at
    )
VALUES
    ($1, $2, NOW(), NOW())
RETURNING
    public.organization_aliases.id;

-- name: AliasOrganizationDelete :one
INSERT INTO
    public.organization_aliases (
        organization_id,
        organization_alias,
        created_at,
        updated_at
    )
VALUES
    ($1, $2, NOW(), NOW())
RETURNING
    public.organization_aliases.id;

-- name: OrganizationAliasGetByItemID :many
SELECT
    public.organization_aliases.*
FROM
    public.organization_aliases
WHERE
    id = $1;

-- name: OrganizationAliasIdNameGet :many
SELECT
    public.organization_aliases.*
FROM
    public.organization_aliases
WHERE
    organization_id = $1
    AND organization_alias = $2;

-- name: OrganizationGetAllAliases :many
SELECT
    *
FROM
    public.organization_aliases
WHERE
    organization_id = $1;

-- name: OrganizationgGetConversationsAuthoredIn :many
SELECT
    public.organization.id AS organization_id,
    public.organization.name AS organization_name,
    public.relation_documents_organizations_authorship.document_id,
    public.docket_conversations.docket_gov_id AS docket_gov_id,
    public.docket_conversations.id AS conversation_uuid
FROM
    public.organization
    LEFT JOIN public.relation_documents_organizations_authorship ON public.organization.id = public.relation_documents_organizations_authorship.organization_id
    LEFT JOIN public.docket_documents ON public.relation_documents_organizations_authorship.document_id = public.docket_documents.file_id
    LEFT JOIN public.docket_conversations ON public.docket_documents.conversation_uuid = public.docket_conversations.id
WHERE
    public.organization.id = $1;

-- name: OrganizationSemiCompleteInfoListPaginated :many
SELECT
    org.id,
    COUNT(org_author.document_id) AS document_count,
    org.name,
    org.description,
    org.created_at,
    org.updated_at
FROM
    public.organization org
    LEFT JOIN public.relation_documents_organizations_authorship org_author ON org_author.organization_id = org.id
GROUP BY
    org.id
ORDER BY
    document_count DESC
LIMIT
    $1 OFFSET $2;

-- name: AuthorshipDocumentListOrganizations :many
SELECT
    rdoa.document_id,
    rdoa.organization_id,
    rdoa.is_primary_author,
    rdoa.created_at,
    rdoa.updated_at
FROM
    public.relation_documents_organizations_authorship rdoa
WHERE
    rdoa.document_id = $1
ORDER BY
    rdoa.is_primary_author DESC,
    rdoa.created_at ASC;
