// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: conversations.sql

package dbstore

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const conversationIDFetchFromFileID = `-- name: ConversationIDFetchFromFileID :many
SELECT
    conversation_uuid, file_id, created_at, updated_at
FROM
    public.docket_documents
WHERE
    file_id = $1
`

func (q *Queries) ConversationIDFetchFromFileID(ctx context.Context, fileID uuid.UUID) ([]DocketDocument, error) {
	rows, err := q.db.Query(ctx, conversationIDFetchFromFileID, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DocketDocument
	for rows.Next() {
		var i DocketDocument
		if err := rows.Scan(
			&i.ConversationUuid,
			&i.FileID,
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

const conversationSemiCompleteInfoList = `-- name: ConversationSemiCompleteInfoList :many
SELECT
    dc.id,
    dc.docket_gov_id,
    dc.state,
    COUNT(dd.file_id) AS document_count,
    dc."name",
    dc.description,
    dc.matter_type,
    dc.industry_type,
    dc.metadata,
    dc.extra,
    dc.date_published,
    dc.created_at,
    dc.updated_at
FROM
    public.docket_conversations dc
    LEFT JOIN public.docket_documents dd ON dd.conversation_uuid = dc.id
GROUP BY
    dc.id
ORDER BY
    document_count DESC
`

type ConversationSemiCompleteInfoListRow struct {
	ID            uuid.UUID
	DocketGovID   string
	State         string
	DocumentCount int64
	Name          string
	Description   string
	MatterType    string
	IndustryType  string
	Metadata      string
	Extra         string
	DatePublished pgtype.Timestamptz
	CreatedAt     pgtype.Timestamp
	UpdatedAt     pgtype.Timestamp
}

func (q *Queries) ConversationSemiCompleteInfoList(ctx context.Context) ([]ConversationSemiCompleteInfoListRow, error) {
	rows, err := q.db.Query(ctx, conversationSemiCompleteInfoList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ConversationSemiCompleteInfoListRow
	for rows.Next() {
		var i ConversationSemiCompleteInfoListRow
		if err := rows.Scan(
			&i.ID,
			&i.DocketGovID,
			&i.State,
			&i.DocumentCount,
			&i.Name,
			&i.Description,
			&i.MatterType,
			&i.IndustryType,
			&i.Metadata,
			&i.Extra,
			&i.DatePublished,
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

const conversationSemiCompleteInfoListPaginated = `-- name: ConversationSemiCompleteInfoListPaginated :many
SELECT
    dc.id,
    dc.docket_gov_id,
    dc.state,
    COUNT(dd.file_id) AS document_count,
    dc."name",
    dc.description,
    dc.matter_type,
    dc.industry_type,
    dc.metadata,
    dc.extra,
    dc.date_published,
    dc.created_at,
    dc.updated_at
FROM
    public.docket_conversations dc
    LEFT JOIN public.docket_documents dd ON dd.conversation_uuid = dc.id
GROUP BY
    dc.id
ORDER BY
    document_count DESC
LIMIT
    $1 OFFSET $2
`

type ConversationSemiCompleteInfoListPaginatedParams struct {
	Limit  int32
	Offset int32
}

type ConversationSemiCompleteInfoListPaginatedRow struct {
	ID            uuid.UUID
	DocketGovID   string
	State         string
	DocumentCount int64
	Name          string
	Description   string
	MatterType    string
	IndustryType  string
	Metadata      string
	Extra         string
	DatePublished pgtype.Timestamptz
	CreatedAt     pgtype.Timestamp
	UpdatedAt     pgtype.Timestamp
}

func (q *Queries) ConversationSemiCompleteInfoListPaginated(ctx context.Context, arg ConversationSemiCompleteInfoListPaginatedParams) ([]ConversationSemiCompleteInfoListPaginatedRow, error) {
	rows, err := q.db.Query(ctx, conversationSemiCompleteInfoListPaginated, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ConversationSemiCompleteInfoListPaginatedRow
	for rows.Next() {
		var i ConversationSemiCompleteInfoListPaginatedRow
		if err := rows.Scan(
			&i.ID,
			&i.DocketGovID,
			&i.State,
			&i.DocumentCount,
			&i.Name,
			&i.Description,
			&i.MatterType,
			&i.IndustryType,
			&i.Metadata,
			&i.Extra,
			&i.DatePublished,
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

const docketConversationCreate = `-- name: DocketConversationCreate :one
INSERT INTO
    public.docket_conversations (
        docket_gov_id,
        state,
        name,
        description,
        matter_type,
        industry_type,
        metadata,
        extra,
        date_published,
        created_at,
        updated_at
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
RETURNING
    id
`

type DocketConversationCreateParams struct {
	DocketGovID   string
	State         string
	Name          string
	Description   string
	MatterType    string
	IndustryType  string
	Metadata      string
	Extra         string
	DatePublished pgtype.Timestamptz
}

func (q *Queries) DocketConversationCreate(ctx context.Context, arg DocketConversationCreateParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, docketConversationCreate,
		arg.DocketGovID,
		arg.State,
		arg.Name,
		arg.Description,
		arg.MatterType,
		arg.IndustryType,
		arg.Metadata,
		arg.Extra,
		arg.DatePublished,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const docketConversationDelete = `-- name: DocketConversationDelete :exec
DELETE FROM
    public.docket_conversations
WHERE
    id = $1
`

func (q *Queries) DocketConversationDelete(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, docketConversationDelete, id)
	return err
}

const docketConversationFetchByDocketIdMatch = `-- name: DocketConversationFetchByDocketIdMatch :many
SELECT
    id, docket_gov_id, state, created_at, updated_at, name, description, matter_type, industry_type, metadata, extra, date_published
FROM
    public.docket_conversations
WHERE
    docket_gov_id = $1
`

func (q *Queries) DocketConversationFetchByDocketIdMatch(ctx context.Context, docketGovID string) ([]DocketConversation, error) {
	rows, err := q.db.Query(ctx, docketConversationFetchByDocketIdMatch, docketGovID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DocketConversation
	for rows.Next() {
		var i DocketConversation
		if err := rows.Scan(
			&i.ID,
			&i.DocketGovID,
			&i.State,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Description,
			&i.MatterType,
			&i.IndustryType,
			&i.Metadata,
			&i.Extra,
			&i.DatePublished,
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

const docketConversationList = `-- name: DocketConversationList :many
SELECT
    id, docket_gov_id, state, created_at, updated_at, name, description, matter_type, industry_type, metadata, extra, date_published
FROM
    public.docket_conversations
ORDER BY
    created_at DESC
`

func (q *Queries) DocketConversationList(ctx context.Context) ([]DocketConversation, error) {
	rows, err := q.db.Query(ctx, docketConversationList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DocketConversation
	for rows.Next() {
		var i DocketConversation
		if err := rows.Scan(
			&i.ID,
			&i.DocketGovID,
			&i.State,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Description,
			&i.MatterType,
			&i.IndustryType,
			&i.Metadata,
			&i.Extra,
			&i.DatePublished,
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

const docketConversationRead = `-- name: DocketConversationRead :one
SELECT
    id, docket_gov_id, state, created_at, updated_at, name, description, matter_type, industry_type, metadata, extra, date_published
FROM
    public.docket_conversations
WHERE
    id = $1
`

func (q *Queries) DocketConversationRead(ctx context.Context, id uuid.UUID) (DocketConversation, error) {
	row := q.db.QueryRow(ctx, docketConversationRead, id)
	var i DocketConversation
	err := row.Scan(
		&i.ID,
		&i.DocketGovID,
		&i.State,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Description,
		&i.MatterType,
		&i.IndustryType,
		&i.Metadata,
		&i.Extra,
		&i.DatePublished,
	)
	return i, err
}

const docketConversationUpdate = `-- name: DocketConversationUpdate :one
UPDATE
    public.docket_conversations
SET
    docket_gov_id = $1,
    state = $2,
    name = $3,
    description = $4,
    matter_type = $5,
    industry_type = $6,
    metadata = $7,
    extra = $8,
    date_published = $9,
    updated_at = NOW()
WHERE
    id = $10
RETURNING
    id
`

type DocketConversationUpdateParams struct {
	DocketGovID   string
	State         string
	Name          string
	Description   string
	MatterType    string
	IndustryType  string
	Metadata      string
	Extra         string
	DatePublished pgtype.Timestamptz
	ID            uuid.UUID
}

func (q *Queries) DocketConversationUpdate(ctx context.Context, arg DocketConversationUpdateParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, docketConversationUpdate,
		arg.DocketGovID,
		arg.State,
		arg.Name,
		arg.Description,
		arg.MatterType,
		arg.IndustryType,
		arg.Metadata,
		arg.Extra,
		arg.DatePublished,
		arg.ID,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const docketDocumentDeleteAll = `-- name: DocketDocumentDeleteAll :exec
DELETE FROM
    public.docket_documents
WHERE
    conversation_uuid = $1
`

func (q *Queries) DocketDocumentDeleteAll(ctx context.Context, conversationUuid uuid.UUID) error {
	_, err := q.db.Exec(ctx, docketDocumentDeleteAll, conversationUuid)
	return err
}

const docketDocumentInsert = `-- name: DocketDocumentInsert :one
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
    conversation_uuid
`

type DocketDocumentInsertParams struct {
	ConversationUuid uuid.UUID
	FileID           uuid.UUID
}

func (q *Queries) DocketDocumentInsert(ctx context.Context, arg DocketDocumentInsertParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, docketDocumentInsert, arg.ConversationUuid, arg.FileID)
	var conversation_uuid uuid.UUID
	err := row.Scan(&conversation_uuid)
	return conversation_uuid, err
}

const docketDocumentUpdate = `-- name: DocketDocumentUpdate :one
UPDATE
    public.docket_documents
SET
    conversation_uuid = $1,
    updated_at = NOW()
WHERE
    file_id = $2
RETURNING
    file_id
`

type DocketDocumentUpdateParams struct {
	ConversationUuid uuid.UUID
	FileID           uuid.UUID
}

func (q *Queries) DocketDocumentUpdate(ctx context.Context, arg DocketDocumentUpdateParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, docketDocumentUpdate, arg.ConversationUuid, arg.FileID)
	var file_id uuid.UUID
	err := row.Scan(&file_id)
	return file_id, err
}
