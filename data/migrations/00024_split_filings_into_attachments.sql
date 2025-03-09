-- +goose Up
-- First create the new tables
CREATE TABLE IF NOT EXISTS public.attachment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID NOT NULL REFERENCES public.file(id),  -- Changed from PRIMARY KEY to UNIQUE
    lang VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    extension VARCHAR NOT NULL,
    hash VARCHAR NOT NULL,
    mdata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.attachment_text_source (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attachment_id UUID NOT NULL REFERENCES public.attachment(id) ON DELETE CASCADE,
    is_original_text BOOLEAN NOT NULL,
    language VARCHAR NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS public.attachment_extras (
    id UUID PRIMARY KEY REFERENCES public.attachment(id),
    extra_obj JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

--  Then migrate the data
-- For this query get each value from file_metadata,  then get the attachemnt_id then copy over the mdata field into the attachment field.
INSERT INTO
    public.attachment (
        file_id,
        lang,
        name,
        extension,
        hash,
        mdata
    )
SELECT
    f.id,
    f.lang,
    f.name,
    f.extension,
    f.hash,
    fm.mdata
FROM
    public.file f
    LEFT JOIN public.file_metadata fm ON f.id = fm.id;

-- Migrate text sources
INSERT INTO
    public.attachment_text_source (
        attachment_id,
        is_original_text,
        language,
        text
    )
SELECT
    a.id,
    fts.is_original_text,
    fts.language,
    fts.text
FROM
    public.file_text_source fts
    JOIN public.attachment a ON fts.file_id = a.file_id;

INSERT INTO
    public.attachment_extras (id, extra_obj)
SELECT
    a.id,
    fextra.extra_obj
FROM
    public.file_extras fextra
    JOIN public.attachment a ON fextra.id = a.file_id;

-- +goose Down
-- Add rollback SQL here if needed
DROP TABLE IF EXISTS public.attachment_extras;

DROP TABLE IF EXISTS public.attachment_text_source;

DROP TABLE IF EXISTS public.attachment;