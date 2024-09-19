-- +migrate Up
CREATE TABLE public.encounter (
    name VARCHAR,
    description VARCHAR,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE public.event (
    date TIMESTAMPTZ,
    name VARCHAR,
    description VARCHAR,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE public.faction (
    name VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
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
    original_text TEXT,
    english_text TEXT,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE public.file_text_source (
    file_id UUID NOT NULL,
    is_original_text BOOLEAN NOT NULL,
    language VARCHAR NOT NULL,
    text TEXT,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    FOREIGN KEY (file_id) REFERENCES public.file(id) ON DELETE CASCADE
);

CREATE TABLE public.individual (
    name VARCHAR NOT NULL,
    username VARCHAR,
    chosen_name VARCHAR,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE public.organization (
    name VARCHAR NOT NULL,
    description VARCHAR,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE public.relation_document_associated_with_organization (
    document_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (document_id) REFERENCES public.file(id) ON DELETE CASCADE,
    FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_document_authored_by_individual (
    document_id UUID NOT NULL,
    individual_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (document_id) REFERENCES public.file(id) ON DELETE CASCADE,
    FOREIGN KEY (individual_id) REFERENCES public.individual(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_documents_in_encounter (
    document_id UUID NOT NULL,
    encounter_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (document_id) REFERENCES public.file(id) ON DELETE CASCADE,
    FOREIGN KEY (encounter_id) REFERENCES public.encounter(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_factions_in_encounter (
    encounter_id UUID NOT NULL,
    faction_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (encounter_id) REFERENCES public.encounter(id) ON DELETE CASCADE,
    FOREIGN KEY (faction_id) REFERENCES public.faction(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_files_associated_with_event (
    file_id UUID NOT NULL,
    event_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (file_id) REFERENCES public.file(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES public.event(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_individuals_associated_with_event (
    individual_id UUID NOT NULL,
    event_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (individual_id) REFERENCES public.individual(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES public.event(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_individuals_currently_associated_organization (
    individual_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (individual_id) REFERENCES public.individual(id) ON DELETE CASCADE,
    FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_individuals_in_faction (
    faction_id UUID NOT NULL,
    individual_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (faction_id) REFERENCES public.faction(id) ON DELETE CASCADE,
    FOREIGN KEY (individual_id) REFERENCES public.individual(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_organizations_associated_with_event (
    organization_id UUID NOT NULL,
    event_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES public.event(id) ON DELETE CASCADE
);

CREATE TABLE public.relation_organizations_in_faction (
    faction_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sa_orm_sentinel INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (faction_id) REFERENCES public.faction(id) ON DELETE CASCADE,
    FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS public.relation_organizations_in_faction;
DROP TABLE IF EXISTS public.relation_organizations_associated_with_event;
DROP TABLE IF EXISTS public.relation_individuals_in_faction;
DROP TABLE IF EXISTS public.relation_individuals_currently_associated_organization;
DROP TABLE IF EXISTS public.relation_individuals_associated_with_event;
DROP TABLE IF EXISTS public.relation_files_associated_with_event;
DROP TABLE IF EXISTS public.relation_factions_in_encounter;
DROP TABLE IF EXISTS public.relation_documents_in_encounter;
DROP TABLE IF EXISTS public.relation_document_authored_by_individual;
DROP TABLE IF EXISTS public.relation_document_associated_with_organization;
DROP TABLE IF EXISTS public.organization;
DROP TABLE IF EXISTS public.individual;
DROP TABLE IF EXISTS public.file_text_source;
DROP TABLE IF EXISTS public.file;
DROP TABLE IF EXISTS public.faction;
DROP TABLE IF EXISTS public.event;
DROP TABLE IF EXISTS public.encounter;
