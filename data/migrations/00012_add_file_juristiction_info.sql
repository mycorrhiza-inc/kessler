-- +goose Up
CREATE TABLE IF NOT EXISTS public.juristiction_information (
  id UUID PRIMARY KEY REFERENCES public.file(id),
  country VARCHAR,
  state VARCHAR,
  municipality VARCHAR,
  agency VARCHAR,
  proceeding_name VARCHAR,
	extra JSONB,
	created_at TIMESTAMPTZ DEFAULT now(),
	updated_at TIMESTAMPTZ DEFAULT now()
);
-- +goose Down
DROP TABLE if Exists public.juristiction_information;
