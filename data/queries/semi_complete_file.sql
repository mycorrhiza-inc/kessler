
-- name: GetFileWithMetadata :one
SELECT f.*,
  fm.mdata
FROM public.file f
  LEFT JOIN public.file_metadata fm ON public.file.id = public.file_metadata.id
WHERE public.file.id = $1;

-- name: GetSemiCompleteFile :many
SELECT 
    f.*,
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
