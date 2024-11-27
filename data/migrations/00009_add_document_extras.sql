-- +goose Up
CREATE TABLE IF NOT EXISTS public.file_extras (
    id UUID PRIMARY KEY REFERENCES public.file(id),
    isPrivate BOOLEAN,
    summary VARCHAR,
    short_summary VARCHAR,
    purpose VARCHAR,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS public.file_extras;