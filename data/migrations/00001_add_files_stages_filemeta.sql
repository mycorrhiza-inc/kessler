-- +goose Up
CREATE TABLE IF NOT EXISTS public.filestage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.file (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lang VARCHAR,
    name VARCHAR,
    extension VARCHAR,
    stage_id UUID REFERENCES public.filestage(id),
    isPrivate BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.file_metadata (
    id UUID PRIMARY KEY REFERENCES public.file(id),
    isPrivate BOOLEAN,
    mdata jsonb,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.file_text_source (
    file_id UUID REFERENCES public.file(id) ON DELETE CASCADE,
    is_original_text BOOLEAN NOT NULL,
    language VARCHAR NOT NULL,
    text TEXT,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS public.file_metadata;

DROP TABLE IF EXISTS public.filestage;

DROP TABLE IF EXISTS public.file;

DROP TABLE IF EXISTS public.file_text_source;