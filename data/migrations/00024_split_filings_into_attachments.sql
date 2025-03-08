-- +goose Up
-- First create the new tables
CREATE TABLE IF NOT EXISTS public.attachment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID REFERENCES public.file(id),  -- Changed from PRIMARY KEY to UNIQUE
    lang VARCHAR,
    name VARCHAR,
    extension VARCHAR,
    hash VARCHAR,
    mdata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS public.attachment_text_source (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attachment_id UUID REFERENCES public.attachment(id) ON DELETE CASCADE,
    is_original_text BOOLEAN NOT NULL,
    language VARCHAR NOT NULL,
    text TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.attachment_extras (
    id UUID PRIMARY KEY REFERENCES public.attachment(id),
    mdata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Then migrate the data
INSERT INTO public.attachment (
    file_id,
    lang,
    name,
    extension,
    hash,
    created_at,
    updated_at
)
SELECT 
    id,
    lang,
    name,
    extension,
    hash,
    created_at,
    updated_at
FROM public.file;

-- Migrate text sources
INSERT INTO public.attachment_text_source (
    attachment_id,
    is_original_text,
    language,
    text,
    created_at,
    updated_at
)
SELECT 
    a.id,
    fts.is_original_text,  -- assuming all existing texts are original
    fts.language,
    fts.text,
    fts.created_at,
    fts.updated_at
FROM public.file_text_source fts
JOIN public.attachment a ON fts.file_id = a.file_id;


-- +goose Down
-- Add rollback SQL here if needed
DROP TABLE IF EXISTS public.attachment_extras;
DROP TABLE IF EXISTS public.attachment_text_source;
DROP TABLE IF EXISTS public.attachment;
