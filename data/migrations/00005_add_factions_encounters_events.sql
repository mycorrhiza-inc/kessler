-- +goose Up
CREATE TABLE public.encounter (
    name VARCHAR,
    description VARCHAR,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE public.event (
    date TIMESTAMPTZ,
    name VARCHAR,
    description VARCHAR,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE public.faction (
    name VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.file (
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
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE public.individual (
    name VARCHAR NOT NULL,
    username VARCHAR,
    chosen_name VARCHAR,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE public.organization (
    name VARCHAR NOT NULL,
    description VARCHAR,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE public.relation_documents_organizations (
    document_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (document_id) REFERENCES public.file(id) ON DELETE CASCADE,
    FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_documents_individuals_author (
    document_id UUID NOT NULL,
    individual_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (document_id) REFERENCES public.file(id) ON DELETE CASCADE,
    FOREIGN KEY (individual_id) REFERENCES public.individual(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_documents_encounters (
    document_id UUID NOT NULL,
    encounter_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (document_id) REFERENCES public.file(id) ON DELETE CASCADE,
    FOREIGN KEY (encounter_id) REFERENCES public.encounter(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_factions_encounters (
    encounter_id UUID NOT NULL,
    faction_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (encounter_id) REFERENCES public.encounter(id) ON DELETE CASCADE,
    FOREIGN KEY (faction_id) REFERENCES public.faction(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_files_events (
    file_id UUID NOT NULL,
    event_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (file_id) REFERENCES public.file(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES public.event(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_individuals_events (
    individual_id UUID NOT NULL,
    event_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (individual_id) REFERENCES public.individual(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES public.event(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_individuals_organizations (
    individual_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (individual_id) REFERENCES public.individual(id) ON DELETE CASCADE,
    FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_individuals_factions (
    faction_id UUID NOT NULL,
    individual_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (faction_id) REFERENCES public.faction(id) ON DELETE CASCADE,
    FOREIGN KEY (individual_id) REFERENCES public.individual(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_organizations_events (
    organization_id UUID NOT NULL,
    event_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES public.event(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_organizations_factions (
    faction_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (faction_id) REFERENCES public.faction(id) ON DELETE CASCADE,
    FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS public.relation_organizations_factions CASCADE;

DROP TABLE IF EXISTS public.relation_organizations_events CASCADE;

DROP TABLE IF EXISTS public.relation_individuals_factions CASCADE;

DROP TABLE IF EXISTS public.relation_individuals_organizations CASCADE;

DROP TABLE IF EXISTS public.relation_individuals_event CASCADE;

DROP TABLE IF EXISTS public.relation_files_events CASCADE;

DROP TABLE IF EXISTS public.relation_factions_encounters CASCADE;

DROP TABLE IF EXISTS public.relation_documents_encounters CASCADE;

DROP TABLE IF EXISTS public.relation_documents_individuals_author CASCADE;

DROP TABLE IF EXISTS public.relation_documents_organizations CASCADE;

DROP TABLE IF EXISTS public.organization CASCADE;

DROP TABLE IF EXISTS public.individual CASCADE;

DROP TABLE IF EXISTS public.file_text_source CASCADE;

DROP TABLE IF EXISTS public.file CASCADE;

DROP TABLE IF EXISTS public.faction CASCADE;

DROP TABLE IF EXISTS public.event CASCADE;

DROP TABLE IF EXISTS public.encounter CASCADE;