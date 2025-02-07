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

-- name: SaveJob :exec 
UPDATE
    public.jobs
set job_priority = $2,
    job_name = $3,
    job_status = $4,
    job_type = $5,
    job_data = $6,
    updated_at = NOW()
WHERE
    id = $1;

-- name: GetJobById :one
SELECT
    id,
    job_priority,
    job_name,
    job_status,
    job_type,
    job_data,
    created_at,
    updated_at
FROM public.jobs
WHERE
    id = $1;

-- name: SetJobStatus :exec
UPDATE
    public.jobs
SET job_status = $2
WHERE
    id = $1;


-- name: DeleteJob :exec
DELETE FROM
    public.jobs
WHERE
    id = $1;
