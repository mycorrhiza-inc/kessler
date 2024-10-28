-- +goose Up
ALTER TABLE public.file
ADD COLUMN file hash DEFAULT '';
-- +goose Down
ALTER TABLE public.file DROP COLUMN hash;
