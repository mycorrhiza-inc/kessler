-- name: CreateJob :one 
INSERT INTO
	public.jobs (
		id,
		created_at,
		updated_at,
		priority,
		name,
		status,
		job_type,
		payload
	)
VALUES
	($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
RETURNING
	id;	

-- name: DeleteJob :exec
DELETE FROM
	public.jobs
WHERE
	id = $1;
	
