-- name: GetFiltersByDataset :many
SELECT
    *
FROM
    filters
WHERE
    dataset = $1
    AND is_active = TRUE
ORDER BY
    created_at DESC;