-- name: CheckIfThaumaturgyAPIKeyExists :one
SELECT *
FROM userfiles.thaumaturgy_api_keys
WHERE key_blake3_hash = $1;
