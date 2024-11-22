-- +goose Up
ALTER TABLE public.docket_conversations 
ADD COLUMN name VARCHAR DEFAULT '';
ALTER TABLE public.docket_conversations
ADD COLUMN description VARCHAR DEFAULT '';
-- +goose Down
ALTER TABLE public.docket_conversations DROP COLUMN name;
ALTER TABLE public.docket_conversations DROP COLUMN description;
