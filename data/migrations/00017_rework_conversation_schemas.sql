-- +goose Up
ALTER TABLE
    public.docket_conversations
ADD
    COLUMN matter_type VARCHAR NOT NULL DEFAULT '';

ALTER TABLE
    public.docket_conversations
ADD
    COLUMN matter_subtype VARCHAR NOT NULL DEFAULT '';

ALTER TABLE
    public.docket_conversations
ADD
    COLUMN industry_type VARCHAR NOT NULL DEFAULT '';

ALTER TABLE
    public.docket_conversations RENAME COLUMN docket_id TO docket_gov_id;

ALTER TABLE
    public.docket_documents RENAME COLUMN docket_id TO conversation_uuid;

-- +goose Down
ALTER TABLE
    public.docket_documents RENAME COLUMN conversation_uuid TO docket_id;

ALTER TABLE
    public.docket_conversations RENAME COLUMN docket_gov_id TO docket_id;

ALTER TABLE
    public.docket_conversations DROP COLUMN industry_type;

ALTER TABLE
    public.docket_conversations DROP COLUMN matter_subtype;

ALTER TABLE
    public.docket_conversations DROP COLUMN matter_type;