-- name: CreateIndividualsAssociatedWithEvent :one
INSERT INTO public.relation_individual_event (
		individual_id,
		event_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id;
-- name: GetIndividualsAssociatedWithEvent :one
SELECT individual_id
FROM public.relation_individual_event
WHERE event_id = $1;
-- name: GetIndividualEventId :one
SELECT id
FROM public.relation_individual_event
WHERE individual_id = $1
	AND event_id = $2;
-- name: GetEventsAssociatedWithIndividual :one
SELECT event_id
FROM public.relation_individual_event
WHERE individual_id = $1;
-- name: DeleteIndividualEventConnection :one
DELETE FROM public.relation_individual_event
WHERE individual_id = $1
	AND event_id = $2;