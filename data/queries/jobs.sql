-- name: CreateJob :one
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
    id;


-- name: DeleteJob :exec
DELETE FROM
    public.jobs
WHERE
    id = $1;
