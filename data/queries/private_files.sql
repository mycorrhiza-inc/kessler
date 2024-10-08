-- name: CreateFile :one
INSERT INTO userfiles.private_file (
		url,
		doctype,
		lang,
		name,
		source,
		hash,
		mdata,
		stage,
		summary,
		short_summary,
		created_at,
		updated_at
	)
VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10,
		NOW(),
		NOW()
	)
RETURNING id;
-- name: ReadPrivateFile :one
SELECT *
FROM userfiles.private_file
WHERE id = $1;
-- name: ListPrivateFiles :many
SELECT *
FROM userfiles.private_file
ORDER BY created_at DESC;
-- name: UpdatePrivateFile :one
UPDATE userfiles.private_file
SET url = $1,
	doctype = $2,
	lang = $3,
	name = $4,
	source = $5,
	hash = $6,
	mdata = $7,
	stage = $8,
	summary = $9,
	short_summary = $10,
	updated_at = NOW()
WHERE id = $11
RETURNING id;
-- name: DeletePrivateFile :exec
DELETE FROM userfiles.private_file
WHERE id = $1;
