-- name: GetFilterString :one
SELECT *
from filter_map
where filter = $1;