// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: semi_complete_file.sql

package dbstore

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const getFileListWithMetadata = `-- name: GetFileListWithMetadata :many
SELECT
    public.file.id,
    public.file.name,
    public.file.extension,
    public.file.lang,
    public.file.verified,
    public.file.hash,
    public.file.isPrivate,
    public.file.created_at,
    public.file.updated_at,
    public.file_metadata.mdata
FROM
    public.file
    LEFT JOIN public.file_metadata ON public.file.id = public.file_metadata.id
WHERE
    public.file.id = ANY($1 :: UUID [])
`

type GetFileListWithMetadataRow struct {
	ID        uuid.UUID
	Name      string
	Extension string
	Lang      string
	Verified  pgtype.Bool
	Hash      string
	Isprivate pgtype.Bool
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	Mdata     []byte
}

func (q *Queries) GetFileListWithMetadata(ctx context.Context, dollar_1 []uuid.UUID) ([]GetFileListWithMetadataRow, error) {
	rows, err := q.db.Query(ctx, getFileListWithMetadata, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFileListWithMetadataRow
	for rows.Next() {
		var i GetFileListWithMetadataRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Extension,
			&i.Lang,
			&i.Verified,
			&i.Hash,
			&i.Isprivate,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Mdata,
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

const getFileWithMetadata = `-- name: GetFileWithMetadata :one
SELECT
    public.file.id,
    public.file.name,
    public.file.extension,
    public.file.lang,
    public.file.verified,
    public.file.hash,
    public.file.isPrivate,
    public.file.created_at,
    public.file.updated_at,
    public.file_metadata.mdata
FROM
    public.file
    LEFT JOIN public.file_metadata ON public.file.id = public.file_metadata.id
WHERE
    public.file.id = $1
`

type GetFileWithMetadataRow struct {
	ID        uuid.UUID
	Name      string
	Extension string
	Lang      string
	Verified  pgtype.Bool
	Hash      string
	Isprivate pgtype.Bool
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	Mdata     []byte
}

func (q *Queries) GetFileWithMetadata(ctx context.Context, id uuid.UUID) (GetFileWithMetadataRow, error) {
	row := q.db.QueryRow(ctx, getFileWithMetadata, id)
	var i GetFileWithMetadataRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Extension,
		&i.Lang,
		&i.Verified,
		&i.Hash,
		&i.Isprivate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Mdata,
	)
	return i, err
}

const semiCompleteFileGet = `-- name: SemiCompleteFileGet :many
SELECT
    public.file.id,
    public.file.name,
    public.file.extension,
    public.file.lang,
    public.file.verified,
    public.file.hash,
    public.file.created_at,
    public.file.updated_at,
    public.file_metadata.mdata,
    public.file_extras.extra_obj,
    public.docket_documents.conversation_uuid AS docket_uuid,
    public.relation_documents_organizations_authorship.is_primary_author,
    public.organization.id AS organization_id,
    public.organization.name AS organization_name,
    public.organization.is_person
FROM
    public.file
    LEFT JOIN public.file_metadata ON public.file.id = public.file_metadata.id
    LEFT JOIN public.file_extras ON public.file.id = public.file_extras.id
    LEFT JOIN public.docket_documents ON public.file.id = public.docket_documents.file_id
    LEFT JOIN public.relation_documents_organizations_authorship ON public.file.id = public.relation_documents_organizations_authorship.document_id
    LEFT JOIN public.organization ON public.relation_documents_organizations_authorship.organization_id = public.organization.id
WHERE
    public.file.id = $1
`

type SemiCompleteFileGetRow struct {
	ID               uuid.UUID
	Name             string
	Extension        string
	Lang             string
	Verified         pgtype.Bool
	Hash             string
	CreatedAt        pgtype.Timestamptz
	UpdatedAt        pgtype.Timestamptz
	Mdata            []byte
	ExtraObj         []byte
	DocketUuid       uuid.UUID
	IsPrimaryAuthor  pgtype.Bool
	OrganizationID   uuid.UUID
	OrganizationName pgtype.Text
	IsPerson         pgtype.Bool
}

func (q *Queries) SemiCompleteFileGet(ctx context.Context, id uuid.UUID) ([]SemiCompleteFileGetRow, error) {
	rows, err := q.db.Query(ctx, semiCompleteFileGet, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SemiCompleteFileGetRow
	for rows.Next() {
		var i SemiCompleteFileGetRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Extension,
			&i.Lang,
			&i.Verified,
			&i.Hash,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Mdata,
			&i.ExtraObj,
			&i.DocketUuid,
			&i.IsPrimaryAuthor,
			&i.OrganizationID,
			&i.OrganizationName,
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

const semiCompleteFileListGet = `-- name: SemiCompleteFileListGet :many
SELECT
    public.file.id,
    public.file.name,
    public.file.extension,
    public.file.lang,
    public.file.verified,
    public.file.hash,
    public.file.created_at,
    public.file.updated_at,
    public.file_metadata.mdata,
    public.file_extras.extra_obj,
    public.docket_documents.conversation_uuid AS docket_uuid,
    array_agg(
        public.organization.id
        ORDER BY
            public.organization.id
    ) :: uuid [] AS organization_ids,
    array_agg(
        public.organization.name
        ORDER BY
            public.organization.id
    ) :: text [] AS organization_names,
    array_agg(
        public.organization.is_person
        ORDER BY
            public.organization.id
    ) :: boolean [] AS is_person_list
FROM
    public.file
    LEFT JOIN public.file_metadata ON public.file.id = public.file_metadata.id
    LEFT JOIN public.file_extras ON public.file.id = public.file_extras.id
    LEFT JOIN public.docket_documents ON public.file.id = public.docket_documents.file_id
    LEFT JOIN public.relation_documents_organizations_authorship ON public.file.id = public.relation_documents_organizations_authorship.document_id
    LEFT JOIN public.organization ON public.relation_documents_organizations_authorship.organization_id = public.organization.id
WHERE
    public.file.id = ANY($1 :: UUID [])
GROUP BY
    FILE.id,
    FILE.name,
    FILE.extension,
    FILE.lang,
    FILE.verified,
    FILE.hash,
    FILE.created_at,
    FILE.updated_at,
    file_metadata.mdata,
    file_extras.extra_obj,
    docket_documents.conversation_uuid
`

type SemiCompleteFileListGetRow struct {
	ID                uuid.UUID
	Name              string
	Extension         string
	Lang              string
	Verified          pgtype.Bool
	Hash              string
	CreatedAt         pgtype.Timestamptz
	UpdatedAt         pgtype.Timestamptz
	Mdata             []byte
	ExtraObj          []byte
	DocketUuid        uuid.UUID
	OrganizationIds   []uuid.UUID
	OrganizationNames []string
	IsPersonList      []bool
}

func (q *Queries) SemiCompleteFileListGet(ctx context.Context, dollar_1 []uuid.UUID) ([]SemiCompleteFileListGetRow, error) {
	rows, err := q.db.Query(ctx, semiCompleteFileListGet, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SemiCompleteFileListGetRow
	for rows.Next() {
		var i SemiCompleteFileListGetRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Extension,
			&i.Lang,
			&i.Verified,
			&i.Hash,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Mdata,
			&i.ExtraObj,
			&i.DocketUuid,
			&i.OrganizationIds,
			&i.OrganizationNames,
			&i.IsPersonList,
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
