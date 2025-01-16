-- +goose Up
ALTER TABLE
    public.docket_conversations
ADD
    COLUMN matter_type VARCHAR NOT NULL DEFAULT '';

ALTER TABLE
    public.docket_conversations
ADD
    COLUMN industry_type VARCHAR NOT NULL DEFAULT '';

ALTER TABLE
    public.docket_conversations RENAME COLUMN docket_id TO docket_gov_id;

ALTER TABLE
    public.docket_conversations
ADD
    COLUMN metadata VARCHAR NOT NULL DEFAULT '';

ALTER TABLE
    public.docket_conversations
ADD
    COLUMN extra VARCHAR NOT NULL DEFAULT '';

ALTER TABLE
    public.docket_documents RENAME COLUMN docket_id TO conversation_uuid;

ALTER TABLE
    public.docket_conversations
ADD
    COLUMN date_published TIMESTAMPTZ NOT NULL DEFAULT '1970-01-01 00:00:00+00';

ALTER TABLE
    public.file
ADD
    COLUMN date_published TIMESTAMPTZ NOT NULL DEFAULT '1970-01-01 00:00:00+00';

ALTER TABLE
    public.docket_conversations DROP COLUMN IF EXISTS deleted_at;


-- +goose Down
ALTER TABLE
    public.file DROP COLUMN date_published;

ALTER TABLE
    public.docket_conversations DROP COLUMN date_published;

ALTER TABLE
    public.docket_conversations DROP COLUMN extra;

ALTER TABLE
    public.docket_conversations DROP COLUMN metadata;

ALTER TABLE
    public.docket_documents RENAME COLUMN conversation_uuid TO docket_id;

ALTER TABLE
    public.docket_conversations RENAME COLUMN docket_gov_id TO docket_id;

ALTER TABLE
    public.docket_conversations DROP COLUMN industry_type;

ALTER TABLE
    public.docket_conversations DROP COLUMN matter_type;

ALTER TABLE
    public.docket_conversations ADD COLUMN deleted_at TIMESTAMPTZ;
