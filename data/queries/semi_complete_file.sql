-- name: GetFileWithMetadata :one
SELECT
    public.file.id,
    public.file.name,
    public.file.extension,
    public.file.lang,
    public.file.verified,
    public.file.hash,
    public.file.isPrivate,
    public.file.created_at,
    public.file.updated_at,
    public.file_metadata.mdata
FROM
    public.file
    LEFT JOIN public.file_metadata ON public.file.id = public.file_metadata.id
WHERE
    public.file.id = $1;

-- name: GetFileListWithMetadata :many
SELECT
    public.file.id,
    public.file.name,
    public.file.extension,
    public.file.lang,
    public.file.verified,
    public.file.hash,
    public.file.isPrivate,
    public.file.created_at,
    public.file.updated_at,
    public.file_metadata.mdata
FROM
    public.file
    LEFT JOIN public.file_metadata ON public.file.id = public.file_metadata.id
WHERE
    public.file.id = ANY($1 :: UUID []);

-- name: SemiCompleteFileGet :many
SELECT
    public.file.id,
    public.file.name,
    public.file.extension,
    public.file.lang,
    public.file.verified,
    public.file.hash,
    public.file.created_at,
    public.file.updated_at,
    public.file_metadata.mdata,
    public.file_extras.extra_obj,
    public.docket_documents.conversation_uuid AS docket_uuid,
    public.relation_documents_organizations_authorship.is_primary_author,
    public.organization.id AS organization_id,
    public.organization.name AS organization_name,
    public.organization.is_person
FROM
    public.file
    LEFT JOIN public.file_metadata ON public.file.id = public.file_metadata.id
    LEFT JOIN public.file_extras ON public.file.id = public.file_extras.id
    LEFT JOIN public.docket_documents ON public.file.id = public.docket_documents.file_id
    LEFT JOIN public.relation_documents_organizations_authorship ON public.file.id = public.relation_documents_organizations_authorship.document_id
    LEFT JOIN public.organization ON public.relation_documents_organizations_authorship.organization_id = public.organization.id
WHERE
    public.file.id = $1;

-- name: SemiCompleteFileListGet :many
SELECT
    public.file.id,
    public.file.name,
    public.file.extension,
    public.file.lang,
    public.file.verified,
    public.file.hash,
    public.file.created_at,
    public.file.updated_at,
    public.file_metadata.mdata,
    public.file_extras.extra_obj,
    public.docket_documents.conversation_uuid AS docket_uuid,
    array_agg(
        public.organization.id
        ORDER BY
            public.organization.id
    ) :: uuid [] AS organization_ids,
    array_agg(
        public.organization.name
        ORDER BY
            public.organization.id
    ) :: text [] AS organization_names,
    array_agg(
        public.organization.is_person
        ORDER BY
            public.organization.id
    ) :: boolean [] AS is_person_list
FROM
    public.file
    LEFT JOIN public.file_metadata ON public.file.id = public.file_metadata.id
    LEFT JOIN public.file_extras ON public.file.id = public.file_extras.id
    LEFT JOIN public.docket_documents ON public.file.id = public.docket_documents.file_id
    LEFT JOIN public.relation_documents_organizations_authorship ON public.file.id = public.relation_documents_organizations_authorship.document_id
    LEFT JOIN public.organization ON public.relation_documents_organizations_authorship.organization_id = public.organization.id
WHERE
    public.file.id = ANY($1 :: UUID [])
GROUP BY
    FILE.id,
    FILE.name,
    FILE.extension,
    FILE.lang,
    FILE.verified,
    FILE.hash,
    FILE.created_at,
    FILE.updated_at,
    file_metadata.mdata,
    file_extras.extra_obj,
    docket_documents.conversation_uuid;