-- +goose Up
ALTER TABLE public.file
ADD COLUMN status hash DEFAULT '';
-- +goose Down
ALTER TABLE public.filestage DROP COLUMN status;
