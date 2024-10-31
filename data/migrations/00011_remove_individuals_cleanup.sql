-- +goose Up
ALTER TABLE IF EXISTS public.relation_documents_organizations 
RENAME TO public.relation_documents_organizations_authorship;


ALTER TABLE public.organization
ADD COLUMN is_person bool DEFAULT False;
ALTER TABLE public.relation_documents_organizations_authorship
ADD COLUMN is_primary_author bool DEFAULT True;

DROP TABLE IF EXISTS public.relation_individuals_factions CASCADE;
DROP TABLE IF EXISTS public.relation_individuals_organizations CASCADE;
DROP TABLE IF EXISTS public.relation_individuals_event CASCADE;
DROP TABLE IF EXISTS public.relation_documents_individuals_author CASCADE;
DROP TABLE IF EXISTS public.individual CASCADE;

-- +goose Down

ALTER TABLE public.organization DROP COLUMN is_person;
ALTER TABLE public.relation_documents_organizations_authorship DROP COLUMN is_primary_author;
ALTER TABLE IF EXISTS public.relation_documents_organizations_authorship
RENAME TO public.relation_documents_organizations;








CREATE TABLE public.individual (
    name VARCHAR NOT NULL,
    username VARCHAR,
    chosen_name VARCHAR,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE public.relation_documents_individuals_author (
    document_id UUID NOT NULL,
    individual_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (document_id) REFERENCES public.file(id) ON DELETE CASCADE,
    FOREIGN KEY (individual_id) REFERENCES public.individual(id) ON DELETE CASCADE
);
CREATE TABLE public.relation_individuals_events (
    individual_id UUID NOT NULL,
    event_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (individual_id) REFERENCES public.individual(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES public.event(id) ON DELETE CASCADE
);
CREATE TABLE public.relation_individuals_organizations (
    individual_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (individual_id) REFERENCES public.individual(id) ON DELETE CASCADE,
    FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE
);
CREATE TABLE public.relation_individuals_factions (
    faction_id UUID NOT NULL,
    individual_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (faction_id) REFERENCES public.faction(id) ON DELETE CASCADE,
    FOREIGN KEY (individual_id) REFERENCES public.individual(id) ON DELETE CASCADE
);

