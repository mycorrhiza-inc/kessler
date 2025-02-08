-- name: FilePrecomputedQuickwitListGetPaginated :many
SELECT
  *
FROM
  public.testmat
LIMIT
    $1 OFFSET $2;
