
-- name: GetFileWithMetadata :one
SELECT *
FROM public.file
  LEFT JOIN public.file_metadata ON public.file.id = public.file_metadata.id
WHERE public.file.id = $1;

-- name: GetSemiCompleteFile : one 

SELECT *
FROM public.file
  LEFT JOIN public.file_metadata ON public.file.id = public.file_metadata.id
  LEFT JOIN public.file_extras ON public.file.id = public.file_extras.id
  LEFT JOIN public.docket_documents on public.file.id = public.docket_documents.document_id
-- I want to get a list of all the organizations that were authors on this docket, by getting all the organization ids associated with the document_id in the table relation_documents_organizations_authorship, then take that foreign key to public.organizations and return a list of that organization info, including the name and id, along with the rest of this schema
  LEFT JOIN public.relation_documents_organizations_authorship ON public.file.id = public.relation_documents_organizations_authorship.document_id
WHERE public.file.id = $1;
