-- +goose Up
ALTER TABLE
    public.stage_log DROP COLUMN stage_id;

ALTER TABLE
    public.file DROP COLUMN stage_id;

DROP TABLE IF EXISTS public.filestage;

ALTER TABLE
    public.stage_log
ADD
    COLUMN file_id UUID REFERENCES public.file(id);

-- +goose Down
CREATE TABLE IF NOT EXISTS public.filestage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE
    public.stage_log DROP COLUMN file_id;

ALTER TABLE
    public.stage_log
ADD
    COLUMN stage_id UUID REFERENCES public.filestage(id);

ALTER TABLE
    public.file
ADD
    COLUMN stage_id UUID REFERENCES public.filestage(id);