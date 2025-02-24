-- name: GetFiltersByState :many
SELECT *
FROM filters
WHERE state = $1 AND is_active = true
ORDER BY created_at DESC;