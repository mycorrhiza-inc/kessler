-- name: CheckIfThaumaturgyAPIKeyExists :one
SELECT 1
FROM userfiles.thaumaturgy_api_keys
WHERE key_blake3_hash = $1;
