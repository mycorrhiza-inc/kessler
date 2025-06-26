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

-- name: AttachmentListByHash :many
SELECT
    *
FROM
    public.attachment
WHERE
    hash = $1
ORDER BY
    created_at DESC;

-- name: GetAllSearchAttachments :many
SELECT
	a.id AS id,
	a.name AS name,
	a.created_at,
	fm.mdata,
	ats.text
FROM
	public.attachment AS a
	LEFT JOIN public.attachment_text_source AS ats
		ON ats.attachment_id = a.id
	LEFT JOIN public.file AS f
		ON f.id = a.file_id
	LEFT JOIN public.file_metadata AS fm
		ON fm.id = f.id
WHERE ats.text IS NOT NULL AND ats.text != '';

-- name: GetSearchAttachmentById :one
SELECT
	a.id AS id,
	a.name AS name,
	a.created_at,
	fm.mdata,
	ats.text
FROM
	public.attachment AS a
	LEFT JOIN public.attachment_text_source AS ats
		ON ats.attachment_id = a.id
	LEFT JOIN public.file AS f
		ON f.id = a.file_id
	LEFT JOIN public.file_metadata AS fm
		ON fm.id = f.id
WHERE a.id = $1;

-- NEW OPTIMIZED QUERIES FOR AUTHOR DATA

-- name: GetAttachmentWithAuthors :one
SELECT
    a.id,
    a.file_id,
    a.lang,
    a.name,
    a.extension,
    a.hash,
    a.mdata,
    a.created_at,
    a.updated_at,
    fm.mdata as file_mdata,
    ats.text,
    COALESCE(
        JSON_AGG(
            JSON_BUILD_OBJECT(
                'author_id', o.id::text,
                'author_name', o.name,
                'is_person', COALESCE(o.is_person, false),
                'is_primary_author', COALESCE(rdoa.is_primary_author, false),
                'description', o.description
            )
        ) FILTER (WHERE o.id IS NOT NULL),
        '[]'::json
    )::text as authors_json
FROM
    public.attachment AS a
    LEFT JOIN public.file AS f ON f.id = a.file_id
    LEFT JOIN public.file_metadata AS fm ON fm.id = f.id
    LEFT JOIN public.attachment_text_source AS ats ON ats.attachment_id = a.id
    LEFT JOIN public.relation_documents_organizations_authorship AS rdoa ON rdoa.document_id = a.file_id
    LEFT JOIN public.organization AS o ON o.id = rdoa.organization_id
WHERE 
    a.id = $1
GROUP BY 
    a.id, a.file_id, a.lang, a.name, a.extension, a.hash, a.mdata, 
    a.created_at, a.updated_at, fm.mdata, ats.text;

-- name: GetAttachmentAuthorsOnly :many
SELECT
    o.id as author_id,
    o.name as author_name,
    COALESCE(o.is_person, false) as is_person,
    COALESCE(rdoa.is_primary_author, false) as is_primary_author,
    o.description as author_description
FROM
    public.attachment AS a
    INNER JOIN public.relation_documents_organizations_authorship AS rdoa 
        ON rdoa.document_id = a.file_id
    INNER JOIN public.organization AS o 
        ON o.id = rdoa.organization_id
WHERE 
    a.id = $1
ORDER BY 
    rdoa.is_primary_author DESC NULLS LAST,
    o.name ASC;

-- name: GetAttachmentsByAuthor :many
SELECT
    a.id,
    a.name,
    a.extension,
    a.created_at,
    COALESCE(rdoa.is_primary_author, false) as is_primary_author,
    CASE WHEN ats.text IS NOT NULL AND ats.text != '' THEN true ELSE false END as has_text
FROM
    public.attachment AS a
    INNER JOIN public.relation_documents_organizations_authorship AS rdoa 
        ON rdoa.document_id = a.file_id
    LEFT JOIN public.attachment_text_source AS ats 
        ON ats.attachment_id = a.id
WHERE 
    rdoa.organization_id = $1
ORDER BY 
    rdoa.is_primary_author DESC NULLS LAST,
    a.created_at DESC
LIMIT $2;

-- name: GetSearchAttachmentsWithAuthors :many
SELECT
    a.id,
    a.name,
    a.created_at,
    fm.mdata,
    ats.text,
    COALESCE(
        JSON_AGG(
            JSON_BUILD_OBJECT(
                'author_id', o.id::text,
                'author_name', o.name,
                'is_person', COALESCE(o.is_person, false),
                'is_primary_author', COALESCE(rdoa.is_primary_author, false)
            )
        ) FILTER (WHERE o.id IS NOT NULL),
        '[]'::json
    )::text as authors_json
FROM
    public.attachment AS a
    LEFT JOIN public.attachment_text_source AS ats ON ats.attachment_id = a.id
    LEFT JOIN public.file AS f ON f.id = a.file_id
    LEFT JOIN public.file_metadata AS fm ON fm.id = f.id
    LEFT JOIN public.relation_documents_organizations_authorship AS rdoa ON rdoa.document_id = a.file_id
    LEFT JOIN public.organization AS o ON o.id = rdoa.organization_id
WHERE 
    ats.text IS NOT NULL AND ats.text != ''
GROUP BY 
    a.id, a.name, a.created_at, fm.mdata, ats.text
ORDER BY 
    a.created_at DESC;

-- name: GetAttachmentSearchStats :one
SELECT 
	COUNT(*) as total_count,
	COUNT(ats.text) as with_text_count,
	COUNT(*) - COUNT(ats.text) as without_text_count
FROM public.attachment AS a
LEFT JOIN public.attachment_text_source AS ats ON ats.attachment_id = a.id;

-- name: GetAttachmentExtensionStats :many
SELECT 
	LOWER(a.extension) as extension,
	COUNT(*) as count
FROM public.attachment AS a
WHERE a.extension IS NOT NULL AND a.extension != ''
GROUP BY LOWER(a.extension)
ORDER BY count DESC
LIMIT $1;

-- name: GetAttachmentsByExtension :many
SELECT
	a.id,
	a.name,
	a.extension,
	a.created_at,
	CASE WHEN ats.text IS NOT NULL AND ats.text != '' THEN true ELSE false END as has_text
FROM public.attachment AS a
LEFT JOIN public.attachment_text_source AS ats ON ats.attachment_id = a.id
WHERE LOWER(a.extension) = LOWER($1)
ORDER BY a.created_at DESC
LIMIT $2;

-- name: GetAttachmentsByDateRange :many
SELECT
	a.id,
	a.name,
	a.extension,
	a.created_at,
	CASE WHEN ats.text IS NOT NULL AND ats.text != '' THEN true ELSE false END as has_text,
	LENGTH(ats.text) as text_length
FROM public.attachment AS a
LEFT JOIN public.attachment_text_source AS ats ON ats.attachment_id = a.id
WHERE a.created_at >= $1 AND a.created_at <= $2
ORDER BY a.created_at DESC
LIMIT $3;

-- name: GetAttachmentsNeedingReindex :many
SELECT DISTINCT a.id
FROM public.attachment AS a
JOIN public.attachment_text_source AS ats ON ats.attachment_id = a.id
WHERE ats.text IS NOT NULL AND ats.text != '';

-- name: GetAttachmentDateRange :one
SELECT 
	MIN(a.created_at) as earliest_date,
	MAX(a.created_at) as latest_date
FROM public.attachment AS a
WHERE a.created_at IS NOT NULL;

-- BATCH QUERIES FOR PERFORMANCE

-- name: GetMultipleAttachmentsWithAuthors :many
SELECT
    a.id,
    a.name,
    a.created_at,
    fm.mdata,
    ats.text,
    COALESCE(
        JSON_AGG(
            JSON_BUILD_OBJECT(
                'author_id', o.id::text,
                'author_name', o.name,
                'is_person', COALESCE(o.is_person, false),
                'is_primary_author', COALESCE(rdoa.is_primary_author, false)
            )
        ) FILTER (WHERE o.id IS NOT NULL),
        '[]'::json
    )::text as authors_json
FROM
    public.attachment AS a
    LEFT JOIN public.attachment_text_source AS ats ON ats.attachment_id = a.id
    LEFT JOIN public.file AS f ON f.id = a.file_id
    LEFT JOIN public.file_metadata AS fm ON fm.id = f.id
    LEFT JOIN public.relation_documents_organizations_authorship AS rdoa ON rdoa.document_id = a.file_id
    LEFT JOIN public.organization AS o ON o.id = rdoa.organization_id
WHERE 
    a.id = ANY($1::uuid[])
GROUP BY 
    a.id, a.name, a.created_at, fm.mdata, ats.text
ORDER BY 
    a.created_at DESC;
