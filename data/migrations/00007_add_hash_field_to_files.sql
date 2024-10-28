-- +goose Up
ALTER TABLE public.file
ADD COLUMN hash VARCHAR;
-- +goose Down
ALTER TABLE public.file DROP COLUMN hash;
