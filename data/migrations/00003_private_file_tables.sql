-- +goose Up


ALTER TABLE userfiles.file RENAME TO private_file;

CREATE TABLE IF NOT EXISTS userfiles.private_file_text_source (
    file_id UUID NOT NULL,
    is_original_text BOOLEAN NOT NULL,
    language VARCHAR NOT NULL,
    text TEXT,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    FOREIGN KEY (file_id) REFERENCES userfiles.private_file(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS userfiles.private_access_controls (
		operator_id UUID NOT NULL,
    operator_table VARCHAR NOT NULL,
		object_id UUID NOT NULL,
    object_table VARCHAR NOT NULL,
    ID UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS usrfiles.private_file;
DROP TABLE IF EXISTS userfiles.private_file_text_source;
DROP TABLE IF EXISTS userfiles.private_access_controls;
