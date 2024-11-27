-- +goose Up

DROP TABLE IF EXISTS public.relation_documents_factions;

DROP TABLE IF EXISTS public.relation_factions_encounters;

DROP TABLE IF EXISTS public.relation_files_events;

DROP TABLE IF EXISTS public.relation_individuals_events;

DROP TABLE IF EXISTS public.relation_organizations_events;

DROP TABLE IF EXISTS public.relation_organizations_factions;

DROP TABLE IF EXISTS public.juristiction_information;

DROP TABLE IF EXISTS public.encounter;

DROP TABLE IF EXISTS public.event;

DROP TABLE IF EXISTS public.faction;

UPDATE
    public.docket_conversations
SET
    name = ''
WHERE
    name IS NULL;

ALTER TABLE
    public.docket_conversations
ALTER COLUMN
    name
SET
    DEFAULT '';

ALTER TABLE
    public.docket_conversations
ALTER COLUMN
    name
SET
    NOT NULL;

UPDATE
    public.docket_conversations
SET
    description = ''
WHERE
    description IS NULL;

ALTER TABLE
    public.docket_conversations
ALTER COLUMN
    description
SET
    DEFAULT '';

ALTER TABLE
    public.docket_conversations
ALTER COLUMN
    description
SET
    NOT NULL;

UPDATE
    public.file
SET
    lang = ''
WHERE
    lang IS NULL;

ALTER TABLE
    public.file
ALTER COLUMN
    lang
SET
    DEFAULT '';

ALTER TABLE
    public.file
ALTER COLUMN
    lang
SET
    NOT NULL;

UPDATE
    public.file
SET
    name = ''
WHERE
    name IS NULL;

ALTER TABLE
    public.file
ALTER COLUMN
    name
SET
    DEFAULT '';

ALTER TABLE
    public.file
ALTER COLUMN
    name
SET
    NOT NULL;

UPDATE
    public.file
SET
    extension = ''
WHERE
    extension IS NULL;

ALTER TABLE
    public.file
ALTER COLUMN
    extension
SET
    DEFAULT '';

ALTER TABLE
    public.file
ALTER COLUMN
    extension
SET
    NOT NULL;

UPDATE
    public.file
SET
    hash = ''
WHERE
    hash IS NULL;

ALTER TABLE
    public.file
ALTER COLUMN
    hash
SET
    DEFAULT '';

ALTER TABLE
    public.file
ALTER COLUMN
    hash
SET
    NOT NULL;

UPDATE
    public.file_text_source
SET
    text = ''
WHERE
    text IS NULL;

ALTER TABLE
    public.file_text_source
ALTER COLUMN
    text
SET
    DEFAULT '';

ALTER TABLE
    public.file_text_source
ALTER COLUMN
    text
SET
    NOT NULL;

UPDATE
    public.organization
SET
    description = ''
WHERE
    description IS NULL;

ALTER TABLE
    public.organization
ALTER COLUMN
    description
SET
    DEFAULT '';

ALTER TABLE
    public.organization
ALTER COLUMN
    description
SET
    NOT NULL;

UPDATE
    public.organization_aliases
SET
    organization_alias = ''
WHERE
    organization_alias IS NULL;

ALTER TABLE
    public.organization_aliases
ALTER COLUMN
    organization_alias
SET
    DEFAULT '';

ALTER TABLE
    public.organization_aliases
ALTER COLUMN
    organization_alias
SET
    NOT NULL;

UPDATE
    public.users
SET
    username = ''
WHERE
    username IS NULL;

ALTER TABLE
    public.users
ALTER COLUMN
    username
SET
    DEFAULT '';

ALTER TABLE
    public.users
ALTER COLUMN
    username
SET
    NOT NULL;

UPDATE
    public.users
SET
    stripe_id = ''
WHERE
    stripe_id IS NULL;

ALTER TABLE
    public.users
ALTER COLUMN
    stripe_id
SET
    DEFAULT '';

ALTER TABLE
    public.users
ALTER COLUMN
    stripe_id
SET
    NOT NULL;

UPDATE
    public.organization
SET
    name = ''
WHERE
    name IS NULL;

ALTER TABLE
    public.organization
ALTER COLUMN
    name
SET
    DEFAULT '';

ALTER TABLE
    public.organization
ALTER COLUMN
    name
SET
    NOT NULL;

UPDATE
    public.docket_conversations
SET
    docket_id = ''
WHERE
    docket_id IS NULL;

ALTER TABLE
    public.docket_conversations
ALTER COLUMN
    docket_id
SET
    DEFAULT '';

ALTER TABLE
    public.docket_conversations
ALTER COLUMN
    docket_id
SET
    NOT NULL;

-- +goose Down
-- FIXME : WRITE THIS SOMETIME IF WE EVER NEED TO REVERT THE DATABASE
