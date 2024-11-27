-- +goose Up
CREATE TABLE public.docket_conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    docket_id VARCHAR(255) NOT NULL,
    state VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.docket_documents (
    docket_id UUID REFERENCES public.docket_conversations(id),
    file_id UUID NOT NULL REFERENCES public.file(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (docket_id, file_id)
);

-- +goose Down
DROP TABLE docket_conversations;

DROP TABLE docket_documents;