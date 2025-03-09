-- name: AttachmentCreate :one
INSERT INTO
    public.attachment (
        file_id,
        lang,
        name,
        extension,
        hash,
        mdata
    )
VALUES
    ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: AttachmentUpdate :one
UPDATE
    public.attachment
SET
    lang = COALESCE($2, lang),
    name = COALESCE($3, name),
    extension = COALESCE($4, extension),
    hash = COALESCE($5, hash),
    mdata = COALESCE($6, mdata),
    updated_at = NOW()
WHERE
    id = $1
RETURNING
    *;

-- name: AttachmentGetById :one
SELECT
    *
FROM
    public.attachment
WHERE
    id = $1;

-- name: AttachmentListByFileId :many
SELECT
    *
FROM
    public.attachment
WHERE
    file_id = $1
ORDER BY
    created_at DESC;