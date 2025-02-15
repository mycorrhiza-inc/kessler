-- +goose Up
CREATE TABLE IF NOT EXISTS filter_map (
	filter TEXT PRIMARY KEY,
	human_readable TEXT
);
-- +goose Down
DROP TABLE IF EXISTS filter_map;