-- name: SemiCompleteFileQuickwitListGet :many
SELECT
    public.file.id,
    public.file.name,
    public.file.extension,
    public.file.lang,
    public.file.verified,
    public.file.hash,
    public.file.created_at,
    public.file.updated_at,
    public.file.date_published,
    public.file_metadata.mdata,
    public.file_extras.extra_obj,
    public.docket_documents.conversation_uuid,
    public.docket_conversations.docket_gov_id,
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
    ) :: boolean [] AS is_person_list,
    array_agg(
        public.file_text_source.text
        ORDER BY
            public.file_text_source.id
    ) :: VARCHAR [] AS file_texts,
    array_agg(
        public.file_text_source.language
        ORDER BY
            public.file_text_source.id
    ) :: VARCHAR [] AS file_text_languages
FROM
    public.file
    LEFT JOIN public.file_metadata ON public.file.id = public.file_metadata.id
    LEFT JOIN public.file_extras ON public.file.id = public.file_extras.id
    LEFT JOIN public.docket_documents ON public.file.id = public.docket_documents.file_id
    LEFT JOIN public.docket_conversations ON public.docket_documents.conversation_id = public.docket_conversations.id
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
    FILE.date_published,
    file_metadata.mdata,
    file_extras.extra_obj,
    docket_documents.conversation_uuid,
    docket_conversations.docket_gov_id;

-- name: OrganizationCompleteQuickwitListGet :many
SELECT
    public.organization.id,
    public.organization.name,
    public.organization.description,
    public.organization.is_person,
    COUNT(
        public.relation_documents_organizations_authorship.document_id
    ) AS total_documents_authored,
    array_agg(
        public.organization_aliases.organization_alias
        ORDER BY
            public.organization_aliases.organization_alias
    ) :: VARCHAR [] AS organization_aliases
FROM
    public.organization
    LEFT JOIN public.organization_aliases ON public.organization.id = public.organization_aliases.organization_id
    LEFT JOIN public.relation_documents_organizations_authorship ON public.organization.id = public.relation_documents_organizations_authorship.organization_id
GROUP BY
    organization.id,
    organization.name,
    organization.description,
    organization.is_person;
