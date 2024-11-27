-- +goose Up
ALTER TABLE
    public.docket_conversations
ADD
    COLUMN name VARCHAR DEFAULT NULL;

ALTER TABLE
    public.docket_conversations
ADD
    COLUMN description VARCHAR DEFAULT NULL;

CREATE TABLE public.organization_aliases (
    organization_alias VARCHAR,
    organization_id UUID NOT NULL,
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE
);

-- +goose Down
ALTER TABLE
    public.docket_conversations DROP COLUMN name;

ALTER TABLE
    public.docket_conversations DROP COLUMN description;

DROP TABLE public.organization_aliases;