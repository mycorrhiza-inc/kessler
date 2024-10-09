-- name: CreatePrivateFileTextSource :one
INSERT INTO userfiles.private_file_text_source (
    file_id,
    is_original_text,
    language,
    text,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id;
-- name: ListPrivateTextsOfFileWithLanguage :many
SELECT *
FROM userfiles.private_file_text_source
WHERE file_id = $1 and language = $2;
-- name: ListPrivateTextsOfFile :many
SELECT *
FROM userfiles.private_file_text_source
WHERE file_id = $1;
-- name: ListPrivateTextsOfFileOriginal :many
SELECT *
FROM userfiles.private_file_text_source
WHERE file_id = $1 and is_original_text = true;
-- name: DeletePrivateFileTexts :exec
DELETE FROM userfiles.private_file_text_source
WHERE file_id = $1;
