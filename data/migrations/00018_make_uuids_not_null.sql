-- For the following tables in public
-- table file_text_source
-- row file_id 
-- table stage_log 
-- row file_id 
-- make the uuid not null with a default for the null uuid
-- +goose Up
ALTER TABLE public.file_text_source
    ALTER COLUMN file_id SET NOT NULL,
    ALTER COLUMN file_id SET DEFAULT '00000000-0000-0000-0000-000000000000';

ALTER TABLE public.stage_log
    ALTER COLUMN file_id SET NOT NULL,
    ALTER COLUMN file_id SET DEFAULT '00000000-0000-0000-0000-000000000000';

-- +goose Down
ALTER TABLE public.file_text_source
    ALTER COLUMN file_id DROP NOT NULL,
    ALTER COLUMN file_id DROP DEFAULT;

ALTER TABLE public.stage_log
    ALTER COLUMN file_id DROP NOT NULL,
    ALTER COLUMN file_id DROP DEFAULT;

