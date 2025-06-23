
-- /Users/orchid/mch/kessler/data/migrations/00026_index_attachment_organizations.sql
-- +goose Up

-- Create index for fast lookup of authors by attachment
-- This enables efficient queries like: "Get all authors for attachment X"
CREATE INDEX idx_relation_documents_organizations_authorship_document_id 
ON public.relation_documents_organizations_authorship (document_id);

-- Create index for fast lookup of attachments by author
-- This enables efficient queries like: "Get all attachments by author Y"  
CREATE INDEX idx_relation_documents_organizations_authorship_organization_id
ON public.relation_documents_organizations_authorship (organization_id);

-- Composite index for fast joins and filtering by primary authorship
CREATE INDEX idx_relation_documents_organizations_authorship_composite
ON public.relation_documents_organizations_authorship (document_id, organization_id, is_primary_author);

-- Since attachments now reference files via file_id, we also need an index 
-- to quickly go from attachment to file for author lookups
CREATE INDEX idx_attachment_file_id ON public.attachment (file_id);

-- +goose Down

DROP INDEX IF EXISTS idx_relation_documents_organizations_authorship_document_id;
DROP INDEX IF EXISTS idx_relation_documents_organizations_authorship_organization_id; 
DROP INDEX IF EXISTS idx_relation_documents_organizations_authorship_composite;
DROP INDEX IF EXISTS idx_attachment_file_id;
