-- +goose Up
CREATE TABLE public.encounter (
    name VARCHAR,
    description VARCHAR,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE public.event (
    date TIMESTAMPTZ,
    name VARCHAR,
    description VARCHAR,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE public.faction (
    name VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE public.file (
    url VARCHAR,
    doctype VARCHAR,
    lang VARCHAR,
    name VARCHAR,
    source VARCHAR,
    hash VARCHAR,
    mdata VARCHAR,
    stage VARCHAR,
    summary VARCHAR,
    short_summary VARCHAR,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE public.file_text_source (
    file_id UUID NOT NULL,
    is_original_text BOOLEAN NOT NULL,
    language VARCHAR NOT NULL,
    text TEXT,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    FOREIGN KEY (file_id) REFERENCES public.file(id) ON DELETE CASCADE
);
CREATE TABLE public.individual (
    name VARCHAR NOT NULL,
    username VARCHAR,
    chosen_name VARCHAR,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE public.organization (
    name VARCHAR NOT NULL,
    description VARCHAR,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE public.relation_document_organization (
    document_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (document_id) REFERENCES public.file(id) ON DELETE CASCADE,
    FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE
);
CREATE TABLE public.relation_document_individual_author (
    document_id UUID NOT NULL,
    individual_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (document_id) REFERENCES public.file(id) ON DELETE CASCADE,
    FOREIGN KEY (individual_id) REFERENCES public.individual(id) ON DELETE CASCADE
);
CREATE TABLE public.relation_document_encounter (
    document_id UUID NOT NULL,
    encounter_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (document_id) REFERENCES public.file(id) ON DELETE CASCADE,
    FOREIGN KEY (encounter_id) REFERENCES public.encounter(id) ON DELETE CASCADE
);
CREATE TABLE public.relation_faction_encounter (
    encounter_id UUID NOT NULL,
    faction_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (encounter_id) REFERENCES public.encounter(id) ON DELETE CASCADE,
    FOREIGN KEY (faction_id) REFERENCES public.faction(id) ON DELETE CASCADE
);
CREATE TABLE public.relation_file_event (
    file_id UUID NOT NULL,
    event_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (file_id) REFERENCES public.file(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES public.event(id) ON DELETE CASCADE
);
CREATE TABLE public.relation_individual_event (
    individual_id UUID NOT NULL,
    event_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (individual_id) REFERENCES public.individual(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES public.event(id) ON DELETE CASCADE
);
CREATE TABLE public.relation_individual_organization (
    individual_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (individual_id) REFERENCES public.individual(id) ON DELETE CASCADE,
    FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE
);
CREATE TABLE public.relation_individual_faction (
    faction_id UUID NOT NULL,
    individual_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (faction_id) REFERENCES public.faction(id) ON DELETE CASCADE,
    FOREIGN KEY (individual_id) REFERENCES public.individual(id) ON DELETE CASCADE
);
CREATE TABLE public.relation_organization_event (
    organization_id UUID NOT NULL,
    event_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES public.event(id) ON DELETE CASCADE
);
CREATE TABLE public.relation_organization_faction (
    faction_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (faction_id) REFERENCES public.faction(id) ON DELETE CASCADE,
    FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE IF EXISTS public.relation_organization_faction;
DROP TABLE IF EXISTS public.relation_organization_event;
DROP TABLE IF EXISTS public.relation_individual_faction;
DROP TABLE IF EXISTS public.relation_individual_organization;
DROP TABLE IF EXISTS public.relation_individual_event;
DROP TABLE IF EXISTS public.relation_file_event;
DROP TABLE IF EXISTS public.relation_faction_encounter;
DROP TABLE IF EXISTS public.relation_document_encounter;
DROP TABLE IF EXISTS public.relation_document_individual_author;
DROP TABLE IF EXISTS public.relation_document_organization;
DROP TABLE IF EXISTS public.organization;
DROP TABLE IF EXISTS public.individual;
DROP TABLE IF EXISTS public.file_text_source;
DROP TABLE IF EXISTS public.file;
DROP TABLE IF EXISTS public.faction;
DROP TABLE IF EXISTS public.event;
DROP TABLE IF EXISTS public.encounter;