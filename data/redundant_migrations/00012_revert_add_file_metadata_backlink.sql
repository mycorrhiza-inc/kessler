-- +goose Up
ALTER TABLE public.file DROP COLUMN metadata_id;
-- +goose Down
ALTER TABLE public.file ADD COLUMN metadata_id UUID REFERENCES public.file_metadata(id);
