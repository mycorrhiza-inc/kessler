-- +goose Up
-- +goose StatementBegin
BEGIN
;

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create filters table
CREATE TABLE filters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    state VARCHAR(255) NOT NULL,
    filter_type VARCHAR(50) NOT NULL DEFAULT 'undefined',
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create filter_dataset_mappings table
CREATE TABLE filter_dataset_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    filter_id UUID NOT NULL REFERENCES filters(id) ON DELETE CASCADE,
    dataset_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(filter_id, dataset_id)
);

-- Create indexes
CREATE INDEX idx_filters_state ON filters(state);

CREATE INDEX idx_filters_filter_type ON filters(filter_type);

CREATE INDEX idx_filters_is_active ON filters(is_active);

CREATE INDEX idx_filter_dataset_mappings_dataset_id ON filter_dataset_mappings(dataset_id);

-- Create updated_at trigger function
CREATE
OR REPLACE FUNCTION update_timestamp() RETURNS TRIGGER AS
$$
BEGIN
NEW.updated_at = CURRENT_TIMESTAMP;

RETURN NEW;

END;

$$
LANGUAGE plpgsql;

-- Create triggers
CREATE TRIGGER update_filters_timestamp BEFORE
UPDATE
    ON filters FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_filter_dataset_mappings_timestamp BEFORE
UPDATE
    ON filter_dataset_mappings FOR EACH ROW EXECUTE FUNCTION update_timestamp();

COMMIT;

-- +goose StatementEnd
-- +goose Down
-- Drop triggers first
DROP TRIGGER IF EXISTS update_filters_timestamp ON filters;

DROP TRIGGER IF EXISTS update_filter_dataset_mappings_timestamp ON filter_dataset_mappings;

-- Drop indexes
DROP INDEX IF EXISTS idx_filters_state;

DROP INDEX IF EXISTS idx_filters_is_active;

DROP INDEX IF EXISTS idx_filter_dataset_mappings_dataset_id;

DROP INDEX IF EXISTS idx_filters_filter_type;

-- Drop tables
DROP TABLE IF EXISTS filter_dataset_mappings;

DROP TABLE IF EXISTS filters;