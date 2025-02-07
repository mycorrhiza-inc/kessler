-- +goose Up
CREATE MATERIALIZED VIEW testmat AS
SELECT 
	public.file.id,
	JSONB_AGG(
		JSONB_BUILD_OBJECT(
			'id',
			public.organization.id,
			'name',
			public.organization.name,
			'is_person',
			public.organization.is_person
		)
	) FILTER (WHERE public.organization.id IS NOT NULL) AS organizations
FROM public.file
	LEFT JOIN public.relation_documents_organizations_authorship 
		ON public.file.id = public.relation_documents_organizations_authorship.document_id
	LEFT JOIN public.organization 
		ON public.relation_documents_organizations_authorship.organization_id = public.organization.id
GROUP BY file.id;
-- +goose Down
DROP MATERIALIZED VIEW IF EXISTS testmat;
