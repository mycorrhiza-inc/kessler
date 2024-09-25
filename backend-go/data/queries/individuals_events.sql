-- name: CreateIndividualsAssociatedWithEvent :one
INSERT INTO public.relation_individuals_events (
		individual_id,
		event_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING id;
-- name: GetIndividualsAssociatedWithEvent :one
SELECT individual_id
FROM public.relation_individuals_events
WHERE event_id = $1;
-- name: GetIndividualEventId :one
SELECT id
FROM public.relation_individuals_events
WHERE individual_id = $1
	AND event_id = $2;
-- name: GetEventsAssociatedWithIndividual :one
SELECT event_id
FROM public.relation_individuals_events
WHERE individual_id = $1;
-- name: DeleteIndividualEventConnection :exec
DELETE FROM public.relation_individuals_events
WHERE individual_id = $1
	AND event_id = $2;