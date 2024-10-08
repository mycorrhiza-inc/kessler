-- +goose Up

-- 
-- -- CREATE TABLE IF NOT EXISTS userfiles.file (
-- --     url VARCHAR,
-- --     doctype VARCHAR,
-- --     lang VARCHAR,
-- --     name VARCHAR,
-- --     source VARCHAR,
-- --     hash VARCHAR,
-- --     mdata VARCHAR,
-- --     stage VARCHAR,
-- --     summary VARCHAR,
-- --     short_summary VARCHAR,
-- --     usergroup_id UUID NOT NULL,
-- --     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
-- --     created_at TIMESTAMPTZ DEFAULT now(),
-- --     updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
-- --     -- when a usergroup is deleted all userfiles associated with it should (This could also be done with a treeshake at some point, also what happens if a user deletes their account, does everyone they shared their docs with loose the doc, sounds leftpad-like.)
-- --     FOREIGN KEY (usergroup_id) REFERENCES public.usergroup(id) ON DELETE CASCADE
-- -- );
--

ALTER TABLE userfiles.file

RENAME TO userfiles.private_file;

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
		updated_at TIMESTAMPTZ DEFAULT now(),
);

-- +goose Down
DROP TABLE IF EXISTS usrfiles.private_file;
DROP TABLE IF EXISTS userfiles.private_file_text_source;
DROP TABLE IF EXISTS userfiles.private_access_controls;
