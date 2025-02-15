// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: jobs.sql

package dbstore

import (
	"context"

	"github.com/google/uuid"
)

const createJob = `-- name: CreateJob :one
INSERT INTO
    public.jobs (
        id,
        job_priority,
        job_name,
        job_status,
        job_type,
        job_data,
        created_at,
        updated_at
    )
VALUES
    ($1, $2, $3, $4, $5, $6, NOW(), NOW())
RETURNING
    id
`

type CreateJobParams struct {
	ID          uuid.UUID
	JobPriority int32
	JobName     string
	JobStatus   string
	JobType     string
	JobData     []byte
}

func (q *Queries) CreateJob(ctx context.Context, arg CreateJobParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createJob,
		arg.ID,
		arg.JobPriority,
		arg.JobName,
		arg.JobStatus,
		arg.JobType,
		arg.JobData,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteJob = `-- name: DeleteJob :exec
DELETE FROM
    public.jobs
WHERE
    id = $1
`

func (q *Queries) DeleteJob(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteJob, id)
	return err
}
