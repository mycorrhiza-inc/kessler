DROP TABLE IF EXISTS public.juristiction_information;
DROP TABLE IF EXISTS public.encounter;
DROP TABLE IF EXISTS public.event;
DROP TABLE IF EXISTS public.faction;
DROP TABLE IF EXISTS public.relation_documents_factions;
DROP TABLE IF EXISTS public.relation_factions_encounters;
DROP TABLE IF EXISTS public.relation_files_events;
DROP TABLE IF EXISTS public.relation_individuals_events;
DROP TABLE IF EXISTS public.relation_organizations_events;
DROP TABLE IF EXISTS public.relation_organizations_factions;


UPDATE public.docket_conversations SET name="" WHERE name IS NULL;
ALTER TABLE public.docket_conversations ALTER COLUMN name SET NOT NULL;
UPDATE public.docket_conversations SET conversation="" WHERE conversation IS NULL;

