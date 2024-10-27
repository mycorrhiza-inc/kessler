-- +goose Up

CREATE TABLE IF NOT EXISTS public.filestage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.file_metadata (
    id UUID PRIMARY KEY REFERENCES public.file(id),
    isPrivate BOOLEAN,
    mdata jsonb,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- create the metadata table
CREATE TABLE IF NOT EXISTS public.file (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR,
    extension VARCHAR,
    stage_id UUID REFERENCES public.filestage(id),
    isPrivate BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);



ALTER TABLE public.file
RENAME COLUMN doctype TO extension;




-- +goose Down
DROP TABLE if Exists public.file_metadata;
DROP TABLE if Exists public.filestage;
DROP TABLE if Exists public.file;

