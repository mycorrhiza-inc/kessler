// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: organizations.sql

package dbstore

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createOrganization = `-- name: CreateOrganization :one
INSERT INTO public.organization (
		name,
		description,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING name, description, id, created_at, updated_at
`

type CreateOrganizationParams struct {
	Name        string
	Description pgtype.Text
}

func (q *Queries) CreateOrganization(ctx context.Context, arg CreateOrganizationParams) (Organization, error) {
	row := q.db.QueryRow(ctx, createOrganization, arg.Name, arg.Description)
	var i Organization
	err := row.Scan(
		&i.Name,
		&i.Description,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteOrganization = `-- name: DeleteOrganization :exec
DELETE FROM public.organization
WHERE id = $1
`

func (q *Queries) DeleteOrganization(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteOrganization, id)
	return err
}

const listOrganizations = `-- name: ListOrganizations :many
SELECT name, description, id, created_at, updated_at
FROM public.organization
ORDER BY created_at DESC
`

func (q *Queries) ListOrganizations(ctx context.Context) ([]Organization, error) {
	rows, err := q.db.Query(ctx, listOrganizations)
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

const readOrganization = `-- name: ReadOrganization :one
SELECT name, description, id, created_at, updated_at
FROM public.organization
WHERE id = $1
`

func (q *Queries) ReadOrganization(ctx context.Context, id pgtype.UUID) (Organization, error) {
	row := q.db.QueryRow(ctx, readOrganization, id)
	var i Organization
	err := row.Scan(
		&i.Name,
		&i.Description,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateOrganization = `-- name: UpdateOrganization :one
UPDATE public.organization
SET name = $1,
	description = $2,
	updated_at = NOW()
WHERE id = $3
RETURNING name, description, id, created_at, updated_at
`

type UpdateOrganizationParams struct {
	Name        string
	Description pgtype.Text
	ID          pgtype.UUID
}

func (q *Queries) UpdateOrganization(ctx context.Context, arg UpdateOrganizationParams) (Organization, error) {
	row := q.db.QueryRow(ctx, updateOrganization, arg.Name, arg.Description, arg.ID)
	var i Organization
	err := row.Scan(
		&i.Name,
		&i.Description,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateOrganizationDescription = `-- name: UpdateOrganizationDescription :one
UPDATE public.organization
SET description = $1,
	updated_at = NOW()
WHERE id = $2
RETURNING name, description, id, created_at, updated_at
`

type UpdateOrganizationDescriptionParams struct {
	Description pgtype.Text
	ID          pgtype.UUID
}

func (q *Queries) UpdateOrganizationDescription(ctx context.Context, arg UpdateOrganizationDescriptionParams) (Organization, error) {
	row := q.db.QueryRow(ctx, updateOrganizationDescription, arg.Description, arg.ID)
	var i Organization
	err := row.Scan(
		&i.Name,
		&i.Description,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateOrganizationName = `-- name: UpdateOrganizationName :one
UPDATE public.organization
SET name = $1,
	updated_at = NOW()
WHERE id = $2
RETURNING name, description, id, created_at, updated_at
`

type UpdateOrganizationNameParams struct {
	Name string
	ID   pgtype.UUID
}

func (q *Queries) UpdateOrganizationName(ctx context.Context, arg UpdateOrganizationNameParams) (Organization, error) {
	row := q.db.QueryRow(ctx, updateOrganizationName, arg.Name, arg.ID)
	var i Organization
	err := row.Scan(
		&i.Name,
		&i.Description,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
