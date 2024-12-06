// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: organizations.sql

package dbstore

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const aliasOrganizationCreate = `-- name: AliasOrganizationCreate :one
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
    public.organization_aliases.id
`

type AliasOrganizationCreateParams struct {
	OrganizationID    uuid.UUID
	OrganizationAlias string
}

func (q *Queries) AliasOrganizationCreate(ctx context.Context, arg AliasOrganizationCreateParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, aliasOrganizationCreate, arg.OrganizationID, arg.OrganizationAlias)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const aliasOrganizationDelete = `-- name: AliasOrganizationDelete :one
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
    public.organization_aliases.id
`

type AliasOrganizationDeleteParams struct {
	OrganizationID    uuid.UUID
	OrganizationAlias string
}

func (q *Queries) AliasOrganizationDelete(ctx context.Context, arg AliasOrganizationDeleteParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, aliasOrganizationDelete, arg.OrganizationID, arg.OrganizationAlias)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const authorshipDocumentDeleteAll = `-- name: AuthorshipDocumentDeleteAll :exec
DELETE FROM
    public.relation_documents_organizations_authorship
WHERE
    document_id = $1
`

func (q *Queries) AuthorshipDocumentDeleteAll(ctx context.Context, documentID uuid.UUID) error {
	_, err := q.db.Exec(ctx, authorshipDocumentDeleteAll, documentID)
	return err
}

const authorshipDocumentOrganizationInsert = `-- name: AuthorshipDocumentOrganizationInsert :one
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
    id
`

type AuthorshipDocumentOrganizationInsertParams struct {
	DocumentID      uuid.UUID
	OrganizationID  uuid.UUID
	IsPrimaryAuthor pgtype.Bool
}

func (q *Queries) AuthorshipDocumentOrganizationInsert(ctx context.Context, arg AuthorshipDocumentOrganizationInsertParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, authorshipDocumentOrganizationInsert, arg.DocumentID, arg.OrganizationID, arg.IsPrimaryAuthor)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const authorshipOrganizationListDocuments = `-- name: AuthorshipOrganizationListDocuments :many
SELECT
    document_id, organization_id, id, created_at, updated_at, is_primary_author
FROM
    public.relation_documents_organizations_authorship
WHERE
    organization_id = $1
`

func (q *Queries) AuthorshipOrganizationListDocuments(ctx context.Context, organizationID uuid.UUID) ([]RelationDocumentsOrganizationsAuthorship, error) {
	rows, err := q.db.Query(ctx, authorshipOrganizationListDocuments, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RelationDocumentsOrganizationsAuthorship
	for rows.Next() {
		var i RelationDocumentsOrganizationsAuthorship
		if err := rows.Scan(
			&i.DocumentID,
			&i.OrganizationID,
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.IsPrimaryAuthor,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createOrganization = `-- name: CreateOrganization :one
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
    ) AS id
`

type CreateOrganizationParams struct {
	OrganizationAlias string
	Description       string
	IsPerson          pgtype.Bool
}

func (q *Queries) CreateOrganization(ctx context.Context, arg CreateOrganizationParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createOrganization, arg.OrganizationAlias, arg.Description, arg.IsPerson)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const organizationAliasGetByItemID = `-- name: OrganizationAliasGetByItemID :many
SELECT
    
FROM
    public.organization_aliases
WHERE
    id = $1
`

type OrganizationAliasGetByItemIDRow struct {
}

func (q *Queries) OrganizationAliasGetByItemID(ctx context.Context, id uuid.UUID) ([]OrganizationAliasGetByItemIDRow, error) {
	rows, err := q.db.Query(ctx, organizationAliasGetByItemID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrganizationAliasGetByItemIDRow
	for rows.Next() {
		var i OrganizationAliasGetByItemIDRow
		if err := rows.Scan(); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const organizationAliasIdNameGet = `-- name: OrganizationAliasIdNameGet :many
SELECT
    
FROM
    public.organization_aliases
WHERE
    organization_id = $1
    AND organization_alias = $2
`

type OrganizationAliasIdNameGetParams struct {
	OrganizationID    uuid.UUID
	OrganizationAlias string
}

type OrganizationAliasIdNameGetRow struct {
}

func (q *Queries) OrganizationAliasIdNameGet(ctx context.Context, arg OrganizationAliasIdNameGetParams) ([]OrganizationAliasIdNameGetRow, error) {
	rows, err := q.db.Query(ctx, organizationAliasIdNameGet, arg.OrganizationID, arg.OrganizationAlias)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrganizationAliasIdNameGetRow
	for rows.Next() {
		var i OrganizationAliasIdNameGetRow
		if err := rows.Scan(); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const organizationDelete = `-- name: OrganizationDelete :exec
DELETE FROM
    public.organization
WHERE
    id = $1
`

func (q *Queries) OrganizationDelete(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, organizationDelete, id)
	return err
}

const organizationFetchByNameMatch = `-- name: OrganizationFetchByNameMatch :many
SELECT
    name, description, id, created_at, updated_at, is_person
FROM
    public.organization
WHERE
    name = $1
`

func (q *Queries) OrganizationFetchByNameMatch(ctx context.Context, name string) ([]Organization, error) {
	rows, err := q.db.Query(ctx, organizationFetchByNameMatch, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Organization
	for rows.Next() {
		var i Organization
		if err := rows.Scan(
			&i.Name,
			&i.Description,
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.IsPerson,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const organizationGetAllAliases = `-- name: OrganizationGetAllAliases :many
SELECT
    
FROM
    public.organization_aliases
WHERE
    organization_id = $1
`

type OrganizationGetAllAliasesRow struct {
}

func (q *Queries) OrganizationGetAllAliases(ctx context.Context, organizationID uuid.UUID) ([]OrganizationGetAllAliasesRow, error) {
	rows, err := q.db.Query(ctx, organizationGetAllAliases, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrganizationGetAllAliasesRow
	for rows.Next() {
		var i OrganizationGetAllAliasesRow
		if err := rows.Scan(); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const organizationList = `-- name: OrganizationList :many
SELECT
    name, description, id, created_at, updated_at, is_person
FROM
    public.organization
ORDER BY
    created_at DESC
`

func (q *Queries) OrganizationList(ctx context.Context) ([]Organization, error) {
	rows, err := q.db.Query(ctx, organizationList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Organization
	for rows.Next() {
		var i Organization
		if err := rows.Scan(
			&i.Name,
			&i.Description,
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.IsPerson,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const organizationRead = `-- name: OrganizationRead :one
SELECT
    name, description, id, created_at, updated_at, is_person
FROM
    public.organization
WHERE
    id = $1
`

func (q *Queries) OrganizationRead(ctx context.Context, id uuid.UUID) (Organization, error) {
	row := q.db.QueryRow(ctx, organizationRead, id)
	var i Organization
	err := row.Scan(
		&i.Name,
		&i.Description,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.IsPerson,
	)
	return i, err
}

const organizationSemiCompleteInfoListPaginated = `-- name: OrganizationSemiCompleteInfoListPaginated :many
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
    $1 OFFSET $2
`

type OrganizationSemiCompleteInfoListPaginatedParams struct {
	Limit  int32
	Offset int32
}

type OrganizationSemiCompleteInfoListPaginatedRow struct {
	ID            uuid.UUID
	DocumentCount int64
	Name          string
	Description   string
	CreatedAt     pgtype.Timestamptz
	UpdatedAt     pgtype.Timestamptz
}

func (q *Queries) OrganizationSemiCompleteInfoListPaginated(ctx context.Context, arg OrganizationSemiCompleteInfoListPaginatedParams) ([]OrganizationSemiCompleteInfoListPaginatedRow, error) {
	rows, err := q.db.Query(ctx, organizationSemiCompleteInfoListPaginated, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrganizationSemiCompleteInfoListPaginatedRow
	for rows.Next() {
		var i OrganizationSemiCompleteInfoListPaginatedRow
		if err := rows.Scan(
			&i.ID,
			&i.DocumentCount,
			&i.Name,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const organizationUpdate = `-- name: OrganizationUpdate :one
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
    id
`

type OrganizationUpdateParams struct {
	Name        string
	Description string
	IsPerson    pgtype.Bool
	ID          uuid.UUID
}

func (q *Queries) OrganizationUpdate(ctx context.Context, arg OrganizationUpdateParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, organizationUpdate,
		arg.Name,
		arg.Description,
		arg.IsPerson,
		arg.ID,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const organizationgGetConversationsAuthoredIn = `-- name: OrganizationgGetConversationsAuthoredIn :many
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
    public.organization.id = $1
`

type OrganizationgGetConversationsAuthoredInRow struct {
	OrganizationID   uuid.UUID
	OrganizationName string
	DocumentID       uuid.UUID
	DocketID         pgtype.Text
	ConversationUuid uuid.UUID
}

func (q *Queries) OrganizationgGetConversationsAuthoredIn(ctx context.Context, id uuid.UUID) ([]OrganizationgGetConversationsAuthoredInRow, error) {
	rows, err := q.db.Query(ctx, organizationgGetConversationsAuthoredIn, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrganizationgGetConversationsAuthoredInRow
	for rows.Next() {
		var i OrganizationgGetConversationsAuthoredInRow
		if err := rows.Scan(
			&i.OrganizationID,
			&i.OrganizationName,
			&i.DocumentID,
			&i.DocketID,
			&i.ConversationUuid,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
