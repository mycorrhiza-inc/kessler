-- +goose Up
-- +goose StatementBegin
BEGIN;
-- Create multiselect_values table for many-to-one relationship
CREATE TABLE multiselect_values (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    filter_id UUID NOT NULL REFERENCES filters(id) ON DELETE CASCADE,
    value TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- Ensure no duplicate values for the same filter
    UNIQUE(filter_id, value)
);
CREATE INDEX idx_multiselect_values_filter_id ON multiselect_values(filter_id);
-- Function to validate filter type
CREATE OR REPLACE FUNCTION check_filter_type() RETURNS TRIGGER AS $$ BEGIN IF NOT EXISTS (
        SELECT 1
        FROM filters
        WHERE id = NEW.filter_id
            AND filter_type = 'multiselect'
    ) THEN RAISE EXCEPTION 'filter_id must reference a multiselect filter';
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- Trigger to ensure only multiselect filters can have values
CREATE TRIGGER ensure_multiselect_filter BEFORE
INSERT
    OR
UPDATE ON multiselect_values FOR EACH ROW EXECUTE FUNCTION check_filter_type();
-- Function to update parent filter timestamp
CREATE OR REPLACE FUNCTION update_parent_filter_timestamp() RETURNS TRIGGER AS $$ BEGIN
UPDATE filters
SET updated_at = CURRENT_TIMESTAMP
WHERE id = CASE
        WHEN TG_OP = 'DELETE' THEN OLD.filter_id
        ELSE NEW.filter_id
    END;
RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER update_filter_on_multiselect_change
AFTER
INSERT
    OR
UPDATE
    OR DELETE ON multiselect_values FOR EACH ROW EXECUTE FUNCTION update_parent_filter_timestamp();
COMMIT;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
BEGIN;
-- Drop trigger
DROP TRIGGER IF EXISTS update_filter_on_multiselect_change ON multiselect_values;
-- Drop function
DROP FUNCTION IF EXISTS update_parent_filter_timestamp();
DROP TABLE IF EXISTS multiselect_values;
COMMIT;
-- +goose StatementEnd