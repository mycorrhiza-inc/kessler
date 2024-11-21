
-- name: GetFileWithMetadata :one
SELECT *
FROM public.file
  LEFT JOIN public.file_metadata ON public.file.id = public.file_metadata.id
WHERE public.file.id = $1;

-- name: GetSemiCompleteFile : one 

SELECT 
    f.*,
    fm.*,
    fe.*,
    dd.*,
    org.id as organization_id,
    org.name as organization_name
FROM public.file f
    LEFT JOIN public.file_metadata fm ON f.id = fm.id
    LEFT JOIN public.file_extras fe ON f.id = fe.id
    LEFT JOIN public.docket_documents dd ON f.id = dd.document_id
    LEFT JOIN public.relation_documents_organizations_authorship rdoa ON f.id = rdoa.document_id
    LEFT JOIN public.organizations org ON rdoa.organization_id = org.id
WHERE f.id = $1;
