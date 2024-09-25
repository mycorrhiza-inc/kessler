-- name: CreateIndividual :one
INSERT INTO public.individual (
		name,
		username,
		chosen_name,
		created_at,
		updated_at
	)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id;
-- name: ReadIndividual :one
SELECT *
FROM public.individual
WHERE id = $1;
-- name: ListIndividuals :many
SELECT *
FROM public.individual
ORDER BY created_at DESC;
-- name: UpdateIndividual :one
UPDATE public.individual
SET name = $1,
	username = $2,
	chosen_name = $3,
	updated_at = NOW()
WHERE id = $4
RETURNING id;
-- name: DeleteIndividual :exec
DELETE FROM public.individual
WHERE id = $1;