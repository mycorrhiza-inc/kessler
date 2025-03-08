-- name: AttachmentTextCreate :one
INSERT INTO
    public.attachment_text_source (
        attachment_id,
        is_original_text,
        language,
        text,
        created_at,
        updated_at
    )
VALUES
    ($1, $2, $3, $4, NOW(), NOW())
RETURNING
    id;

-- name: AttachmentTextList :many
SELECT
    *
FROM
    public.attachment_text_source
WHERE
    attachment_id = $1;

-- name: AttachmentTextListByLanguage :many
SELECT
    *
FROM
    public.attachment_text_source
WHERE
    attachment_id = $1
    AND language = $2;

-- name: AttachmentTextListOriginal :many
SELECT
    *
FROM
    public.attachment_text_source
WHERE
    attachment_id = $1
    AND is_original_text = TRUE;

-- name: AttachmentTextDelete :exec
DELETE FROM
    public.attachment_text_source
WHERE
    attachment_id = $1;

-- name: AttachmentTextListByFileId :many
SELECT
    ats.*
FROM
    public.attachment_text_source ats
    JOIN public.attachment a ON a.id = ats.attachment_id
WHERE
    a.file_id = $1;

-- name: AttachmentTextListByFileIdAndLanguage :many
SELECT
    ats.*
FROM
    public.attachment_text_source ats
    JOIN public.attachment a ON a.id = ats.attachment_id
WHERE
    a.file_id = $1
    AND ats.language = $2;
