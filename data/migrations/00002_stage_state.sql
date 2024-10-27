-- +goose Up
CREATE TYPE stage_state AS ENUM ('pending', 'processing', 'completed');
CREATE TABLE IF NOT EXISTS public.stage_log (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	stage_id UUID REFERENCES public.filestage(id),
	status stage_state,
	-- log of the stage state in json
	log jsonb,
	created_at TIMESTAMPTZ DEFAULT now()
);
ALTER TABLE public.filestage
ADD COLUMN status stage_state DEFAULT 'pending';
-- +goose Down
ALTER TABLE public.filestage DROP COLUMN status;
DROP TYPE IF EXISTS stage_state CASCADE;
DROP TABLE IF EXISTS public.stage_log;