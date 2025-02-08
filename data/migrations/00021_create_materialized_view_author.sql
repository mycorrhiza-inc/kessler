-- +goose Up
CREATE MATERIALIZED VIEW testmat AS
SELECT 
    f.id,
    f.name,
    f.extension,
    f.lang,
    f.verified,
    f.hash,
    f.created_at,
    f.updated_at,
    f.date_published,
    fm.mdata,
    fe.extra_obj,
    dd.conversation_uuid,
    dc.docket_gov_id,
    COALESCE(fts.text, '') AS file_text,
    JSONB_AGG(
        JSONB_BUILD_OBJECT(
            'id', o.id,
            'name', o.name,
            'is_person', o.is_person
        )
    ) FILTER (WHERE o.id IS NOT NULL) AS organizations
FROM public.file AS f
  LEFT JOIN public.file_metadata fm ON f.id = fm.id
  LEFT JOIN public.file_extras fe ON f.id = fe.id
  LEFT JOIN (
      SELECT DISTINCT ON (file_id) file_id, text
      FROM public.file_text_source
      WHERE language = 'en'
      ORDER BY file_id, id
  ) fts ON f.id = fts.file_id
  LEFT JOIN public.docket_documents dd ON f.id = dd.file_id
  LEFT JOIN public.docket_conversations dc ON dd.conversation_uuid = dc.id
  LEFT JOIN public.relation_documents_organizations_authorship rdoa ON f.id = rdoa.document_id
  LEFT JOIN public.organization o ON rdoa.organization_id = o.id
GROUP BY 
    f.id,
    f.name,
    f.extension,
    f.lang,
    f.verified,
    f.hash,
    f.created_at,
    f.updated_at,
    f.date_published,
    fm.mdata,
    fe.extra_obj,
    dd.conversation_uuid,
    dc.docket_gov_id,
    fts.text;
-- +goose Down
DROP MATERIALIZED VIEW IF EXISTS testmat;
