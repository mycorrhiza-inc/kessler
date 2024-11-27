-- +goose Up
CREATE TABLE IF NOT EXISTS public.juristiction_information (
    id UUID PRIMARY KEY REFERENCES public.file(id),
    country VARCHAR,
    state VARCHAR,
    municipality VARCHAR,
    agency VARCHAR,
    proceeding_name VARCHAR,
    extra JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS public.juristiction_information;