-- name: GetFiltersByState :many
SELECT
    *
FROM
    filters
WHERE
    state = $1
    AND is_active = TRUE
ORDER BY
    created_at DESC;