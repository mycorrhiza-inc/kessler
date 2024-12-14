-- name: CreateFile :one
INSERT INTO
    public.file (
        id,
        extension,
        lang,
        name,
        isPrivate,
        hash,
        verified,
        created_at,
        updated_at
    )
VALUES
    (
        gen_random_uuid(),
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        NOW(),
        NOW()
    )
RETURNING
    id;

-- name: GetFile :one
SELECT
    *
FROM
    public.file
WHERE
    id = $1;

-- name: HashGetFileID :many
SELECT
    id
FROM
    public.file
WHERE
    public.file.hash = $1;

-- name: UpdateFile :exec
UPDATE
    public.file
SET
    extension = $1,
    lang = $2,
    name = $3,
    isPrivate = $4,
    hash = $5,
    verified = $6,
    updated_at = NOW()
WHERE
    id = $7;

-- name: ReadFile :one
SELECT
    *
FROM
    public.file
WHERE
    id = $1;

-- name: FilesList :many
SELECT
    *
FROM
    public.file
ORDER BY
    updated_at DESC;

-- name: FilesListUnverified :many
SELECT
    *
FROM
    public.file
WHERE
    verified = false
ORDER BY
    RANDOM()
LIMIT
    $1;

-- name: DeleteFile :exec
DELETE FROM
    public.file
WHERE
    id = $1;

-- name: FileVerifiedUpdate :one
UPDATE
    public.file
SET
    verified = $1,
    updated_at = NOW()
WHERE
    public.file.id = $2
RETURNING
    id;

-- name: StageLogAdd :one
-- used to log the state of a file processing stage and update filestage status
INSERT INTO
    public.stage_log (file_id, STATUS, log)
VALUES
    ($1, $2, $3)
RETURNING
    id,
    file_id,
    STATUS;

-- name: StageLogFileGetLatest :one
SELECT
    *
FROM
    public.stage_log
WHERE
    file_id = $1
ORDER BY
    created_at DESC
LIMIT
    1;

-- name: InsertMetadata :one
INSERT INTO
    public.file_metadata (
        id,
        isPrivate,
        mdata,
        created_at,
        updated_at
    )
VALUES
    ($1, $2, $3, NOW(), NOW())
RETURNING
    id;

-- name: UpdateMetadata :one
UPDATE
    public.file_metadata
SET
    isPrivate = $1,
    mdata = $2,
    updated_at = NOW()
WHERE
    id = $3
RETURNING
    id;

-- name: FetchMetadata :one
SELECT
    *
FROM
    public.file_metadata
WHERE
    id = $1;

-- name: FetchMetadataList :many
SELECT
    *
FROM
    public.file_metadata
WHERE
    id = ANY($1 :: UUID []);

-- name: ExtrasFileCreate :one
INSERT INTO
    public.file_extras (
        id,
        isPrivate,
        extra_obj,
        created_at,
        updated_at
    )
VALUES
    ($1, $2, $3, NOW(), NOW())
RETURNING
    id;

-- name: ExtrasFileUpdate :one
UPDATE
    public.file_extras
SET
    isPrivate = $1,
    extra_obj = $2,
    updated_at = NOW()
WHERE
    id = $3
RETURNING
    id;

-- name: ExtrasFileFetch :one
SELECT
    *
FROM
    public.file_extras
WHERE
    id = $1;