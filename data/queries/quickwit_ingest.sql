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

-- name: ConversationCompleteQuickwitListGet :many
SELECT
    public.docket_conversations.id,
    public.docket_conversations.docket_gov_id,
    public.docket_conversations.state,
    public.docket_conversations.name,
    public.docket_conversations.description,
    public.docket_conversations.matter_type,
    public.docket_conversations.industry_type,
    public.docket_conversations.metadata,
    public.docket_conversations.extra,
    public.docket_conversations.date_published,
    public.docket_conversations.created_at,
    public.docket_conversations.updated_at,
    COUNT(public.docket_documents.file_id) AS total_documents
FROM
    public.docket_conversations
    LEFT JOIN public.docket_documents ON public.docket_conversations.id = public.docket_documents.conversation_uuid
GROUP BY
    docket_conversations.id,
    docket_conversations.docket_gov_id,
    docket_conversations.state,
    docket_conversations.name,
    docket_conversations.description,
    docket_conversations.matter_type,
    docket_conversations.industry_type,
    docket_conversations.metadata,
    docket_conversations.extra,
    docket_conversations.date_published,
    docket_conversations.created_at,
    docket_conversations.updated_at;
