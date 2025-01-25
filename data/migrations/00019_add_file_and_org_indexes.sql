-- +goose Up
CREATE INDEX ON public.file (id);

CREATE INDEX ON public.relation_documents_organizations_authorship (document_id, organization_id);

-- +goose Down
DROP INDEX IF EXISTS public.file_id_idx;

DROP INDEX IF EXISTS public.relation_documents_organizations_authorship_document_id_organization_id_idx;