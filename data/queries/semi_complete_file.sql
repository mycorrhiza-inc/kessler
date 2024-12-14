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
    public.docket_documents.docket_id AS docket_uuid,
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
    public.docket_documents.docket_id AS docket_uuid,
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
    public.file.id = ANY($1 :: UUID []);



-- name: Example Query
-- could you edit the organization name and id, so that it ends up as a ordered list of ids corresponding 
-- to a list of organization names. Since all the other collumns are unique, but a single document can have 
-- up to 30 authors
--
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
    public.docket_documents.docket_id AS docket_uuid,
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
    public.file.id = ANY(ARRAY[
        'b797bc34-2930-4c32-8b4e-0fd0b8878690',
        'd48a3f9e-6baf-4685-b0bf-e4a5480c7bd6',
        '07ddb280-32f3-415b-8680-b559132bd393',
        'f8466f2c-2bd8-49bb-bd9e-0ced5f008d42',
        'df79e1b0-b955-40b3-8ced-4d773be59402',
        '170ceb0e-19d2-45fe-a0ad-35713f037dcb',
        '2234a429-9f61-41cf-85c3-42ebffc4312d',
        '243d2a90-4cbd-40af-9f52-01fdbbda1ecf'
    ]::uuid[]);
