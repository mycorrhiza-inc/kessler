
-- name: GetFileWithMetadata :one
SELECT public.file.id,
    public.file.name,
    public.file.extension,
    public.file.lang,
    public.file.verified,
    public.file.hash,
    public.file.isPrivate,
    public.file.created_at,
    public.file.updated_at,
    public.file_metadata.mdata
FROM public.file
    LEFT JOIN public.file_metadata ON public.file.id = public.file_metadata.id
WHERE public.file.id = $1;

-- name: SemiCompleteFileGet :many
SELECT 
    f.id,
    f.name,
    f.extension,
    f.lang,
    f.verified,
    f.hash,
    f.created_at,
    f.updated_at,
    fm.mdata,
    fe.extra_obj,
    dd.docket_id as docket_uuid,
    org.id as organization_id,
    org.name as organization_name
FROM public.file f
    LEFT JOIN public.file_metadata fm ON f.id = fm.id
    LEFT JOIN public.file_extras fe ON f.id = fe.id
    LEFT JOIN public.docket_documents dd ON f.id = dd.document_id
    LEFT JOIN public.relation_documents_organizations_authorship rdoa ON f.id = rdoa.document_id
    LEFT JOIN public.organization org ON rdoa.organization_id = org.id
WHERE f.id = $1;
