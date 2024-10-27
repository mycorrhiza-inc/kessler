-- +goose Up
CREATE TABLE IF NOT EXISTS public.private_access_controls (
	operator_id UUID NOT NULL,
	operator_table VARCHAR NOT NULL,
	object_id UUID NOT NULL,
	object_table VARCHAR NOT NULL,
	ID UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	created_at TIMESTAMPTZ DEFAULT now(),
	updated_at TIMESTAMPTZ DEFAULT now()
);
CREATE SCHEMA IF NOT EXISTS userfiles;
CREATE TABLE IF NOT EXISTS userfiles.thaumaturgy_api_keys (
	key_name VARCHAR,
	key_blake3_hash VARCHAR PRIMARY KEY NOT NULL,
	ID UUID DEFAULT gen_random_uuid(),
	created_at TIMESTAMPTZ DEFAULT now(),
	updated_at TIMESTAMPTZ DEFAULT now()
);
-- +goose Down
DROP TABLE IF EXISTS usrfiles.private_file;
DROP TABLE IF EXISTS userfiles.private_file_text_source;
DROP TABLE IF EXISTS userfiles.private_access_controls;
DROP TABLE IF EXISTS userfiles.thaumaturgy_api_keys;