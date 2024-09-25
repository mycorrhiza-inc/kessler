// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: individuals_organizations.sql

package dbstore

import (
	"context"

	"github.com/google/uuid"
)

const createIndividualsCurrentlyAssociatedWithOrganization = `-- name: CreateIndividualsCurrentlyAssociatedWithOrganization :one
INSERT INTO public.relation_individuals_organizations (
		individual_id,
		organization_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id
`

type CreateIndividualsCurrentlyAssociatedWithOrganizationParams struct {
	IndividualID   uuid.UUID
	OrganizationID uuid.UUID
}

func (q *Queries) CreateIndividualsCurrentlyAssociatedWithOrganization(ctx context.Context, arg CreateIndividualsCurrentlyAssociatedWithOrganizationParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createIndividualsCurrentlyAssociatedWithOrganization, arg.IndividualID, arg.OrganizationID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteIndividualsCurrentlyAssociatedWithOrganization = `-- name: DeleteIndividualsCurrentlyAssociatedWithOrganization :exec
DELETE FROM public.relation_individuals_organizations
WHERE id = $1
`

func (q *Queries) DeleteIndividualsCurrentlyAssociatedWithOrganization(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteIndividualsCurrentlyAssociatedWithOrganization, id)
	return err
}

const listIndividualsCurrentlyAssociatedWithOrganization = `-- name: ListIndividualsCurrentlyAssociatedWithOrganization :many
SELECT individual_id, organization_id, id, created_at, updated_at
FROM public.relation_individuals_organizations
ORDER BY created_at DESC
`

func (q *Queries) ListIndividualsCurrentlyAssociatedWithOrganization(ctx context.Context) ([]RelationIndividualsOrganization, error) {
	rows, err := q.db.QueryContext(ctx, listIndividualsCurrentlyAssociatedWithOrganization)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RelationIndividualsOrganization
	for rows.Next() {
		var i RelationIndividualsOrganization
		if err := rows.Scan(
			&i.IndividualID,
			&i.OrganizationID,
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const readIndividualsCurrentlyAssociatedWithOrganization = `-- name: ReadIndividualsCurrentlyAssociatedWithOrganization :one
SELECT individual_id, organization_id, id, created_at, updated_at
FROM public.relation_individuals_organizations
WHERE id = $1
`

func (q *Queries) ReadIndividualsCurrentlyAssociatedWithOrganization(ctx context.Context, id uuid.UUID) (RelationIndividualsOrganization, error) {
	row := q.db.QueryRowContext(ctx, readIndividualsCurrentlyAssociatedWithOrganization, id)
	var i RelationIndividualsOrganization
	err := row.Scan(
		&i.IndividualID,
		&i.OrganizationID,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateIndividualsCurrentlyAssociatedWithOrganization = `-- name: UpdateIndividualsCurrentlyAssociatedWithOrganization :one
UPDATE public.relation_individuals_organizations
SET individual_id = $1,
	organization_id = $2,
	updated_at = NOW()
WHERE id = $3
RETURNING id
`

type UpdateIndividualsCurrentlyAssociatedWithOrganizationParams struct {
	IndividualID   uuid.UUID
	OrganizationID uuid.UUID
	ID             uuid.UUID
}

func (q *Queries) UpdateIndividualsCurrentlyAssociatedWithOrganization(ctx context.Context, arg UpdateIndividualsCurrentlyAssociatedWithOrganizationParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, updateIndividualsCurrentlyAssociatedWithOrganization, arg.IndividualID, arg.OrganizationID, arg.ID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}
