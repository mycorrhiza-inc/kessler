-- name: CreatePrivateAccessControl :one
INSERT INTO userfiles.private_access_controls (
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
FROM userfiles.private_access_controls
WHERE operator_id = $1;
-- name: ListOperatorsCanAcessObject :many
SELECT *
FROM userfiles.private_access_controls
WHERE object_id = $1;
-- name: CheckOperatorAccessToObject :one
SELECT EXISTS(
    SELECT 1
    FROM userfiles.private_access_controls
    WHERE operator_id = $1 AND object_id = $2
);
-- name: DeleteAccessControl :exec
DELETE FROM userfiles.private_access_controls
WHERE id = $1;
-- name: RevokeAccessForOperatorOnObject :exec
DELETE FROM userfiles.private_access_controls
WHERE operator_id = $1 AND object_id = $2;
