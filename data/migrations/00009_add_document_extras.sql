-- +goose Up
CREATE TABLE IF NOT EXISTS public.file_extras (
    id UUID PRIMARY KEY REFERENCES public.file(id),
    isPrivate BOOLEAN,
    summary VARCHAR,
    short_summary VARCHAR,
    purpose VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose Down
DROP TABLE if Exists public.file_extras;
