-- +goose Up
CREATE TABLE IF NOT EXISTS public.juristiction_information (
    id UUID PRIMARY KEY REFERENCES public.file(id),
    country VARCHAR,
    state VARCHAR,
    municipality VARCHAR,
    agency VARCHAR,
    proceeding_name VARCHAR,
    extra JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- name: JuristictionFileInsert :one
INSERT INTO
    public.juristiction_information (
        id,
        country,
        state,
        municipality,
        agency,
        proceeding_name,
        extra,
        created_at,
        updated_at
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
RETURNING
    id;

-- name: JuristictionFileUpdate :one
UPDATE
    public.juristiction_information
SET
    country = $1,
    state = $2,
    municipality = $3,
    agency = $4,
    proceeding_name = $5,
    extra = $6,
    updated_at = NOW()
WHERE
    id = $7
RETURNING
    id;

-- name: JuristictionFileFetch :many
SELECT
    *
FROM
    public.juristiction_information
WHERE
    id = $1;