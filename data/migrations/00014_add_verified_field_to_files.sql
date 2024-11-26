-- +goose Up
ALTER TABLE
    public.file
ADD
    COLUMN verified bool DEFAULT False;

-- +goose Down
ALTER TABLE
    public.file DROP COLUMN verified;