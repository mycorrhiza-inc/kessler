-- name: CreateUserFileTable :exec
CREATE TABLE IF NOT EXISTS userfiles {
    url VARCHAR,
    doctype VARCHAR,
    lang VARCHAR,
    name VARCHAR,
    source VARCHAR,
    hash VARCHAR,
    mdata VARCHAR,
    stage VARCHAR,
    summary VARCHAR,
    short_summary VARCHAR,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
};
-- name: DeleteUserFileTable :exec
DROP TABLE IF EXISTS userfiles.$1;