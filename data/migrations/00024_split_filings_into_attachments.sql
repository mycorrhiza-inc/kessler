-- +goose Up

CREATE TABLE IF NOT EXISTS public.attachment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID PRIMARY KEY REFERENCES public.file(id),
    lang VARCHAR,
    name VARCHAR,
    extension VARCHAR,
    hash VARCHAR,
    mdata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS public.attachment_text_source (
    attachment_id UUID REFERENCES public.attachment(id) ON DELETE CASCADE,
    is_original_text BOOLEAN NOT NULL,
    language VARCHAR NOT NULL,
    text TEXT,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.attachment_extras (
    id UUID PRIMARY KEY REFERENCES public.attachment(id),
    isPrivate BOOLEAN,
    mdata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
