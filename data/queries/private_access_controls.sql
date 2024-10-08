-- name: CreatePrivateAccessControl :one
INSERT INTO public.private_access_controls (
		operator_id,
    operator_table,
		object_id,
    object_table,
		created_at,
		updated_at
	)
VALUES ($1, $2,$3, $4, NOW(), NOW())
RETURNING id;
-- name: ListAcessesForOperator :many
SELECT *
FROM public.private_access_controls
WHERE operator_id = $1;
-- name: ListOperatorsCanAcessObject :many
SELECT *
FROM public.private_access_controls
WHERE object_id = $1;
-- name: DeleteAccessControl :exec
DELETE FROM public.private_access_controls
WHERE id = $1;
