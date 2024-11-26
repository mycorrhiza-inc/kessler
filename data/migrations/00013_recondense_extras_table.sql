-- +goose Up
ALTER TABLE
    public.file_extras DROP COLUMN summary;

ALTER TABLE
    public.file_extras DROP COLUMN short_summary;

ALTER TABLE
    public.file_extras DROP COLUMN purpose;

ALTER TABLE
    public.file_extras
ADD
    COLUMN extra_obj JSONB;

-- +goose Down
ALTER TABLE
    public.file_extras DROP COLUMN extra_obj;

ALTER TABLE
    public.file_extras
ADD
    COLUMN purpose TEXT;

ALTER TABLE
    public.file_extras
ADD
    COLUMN short_summary TEXT;

ALTER TABLE
    public.file_extras
ADD
    COLUMN summary TEXT;