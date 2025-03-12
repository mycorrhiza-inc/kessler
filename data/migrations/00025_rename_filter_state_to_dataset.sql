-- +goose Up
-- +goose StatementBegin
BEGIN;

-- Rename the state column to dataset
ALTER TABLE filters RENAME COLUMN state TO dataset;

-- Rename the index to match the new column name
DROP INDEX IF EXISTS idx_filters_state;
CREATE INDEX idx_filters_dataset ON filters(dataset);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;

-- Revert the column name change
ALTER TABLE filters RENAME COLUMN dataset TO state;

-- Revert the index name change
DROP INDEX IF EXISTS idx_filters_dataset;
CREATE INDEX idx_filters_state ON filters(state);

COMMIT;
-- +goose StatementEnd