-- +goose Up
CREATE TABLE IF NOT EXISTS public.filestage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS public.file (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR,
    extension VARCHAR,
    stage_id UUID REFERENCES public.filestage(id),
    isPrivate BOOLEAN DEFAULT FALSE,
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
CREATE TABLE IF NOT EXISTS public.file_text_source (
    file_id UUID REFERENCES public.file(id) ON DELETE CASCADE,
    is_original_text BOOLEAN NOT NULL,
    language VARCHAR NOT NULL,
    text TEXT,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
-- +goose Down
DROP TABLE if Exists public.file_metadata;
DROP TABLE if Exists public.filestage;
DROP TABLE if Exists public.file;
DROP TABLE if Exists public.file_text_source;