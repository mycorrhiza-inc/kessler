-- name: CreateFactionsInEncounter :one
INSERT INTO public.relation_factions_encounters (
		encounter_id,
		faction_id,
		created_at,
		updated_at
	)
VALUES ($1, $2, NOW(), NOW())
RETURNING *;
-- name: ListFactionsInEncounter :one
SELECT faction_id
FROM public.relation_factions_encounters
WHERE encounter_id = $1;
-- name: ListEncountersThatIncludeFaction :many
SELECT encounter_id
FROM public.relation_factions_encounters
WHERE faction_id = $1;
-- name: UpdateFactionInEncounter :one
UPDATE public.relation_factions_encounters
SET encounter_id = $1,
	faction_id = $2,
	updated_at = NOW()
WHERE id = $3
RETURNING *;
-- name: DeleteEncounterFaction :exec
DELETE FROM public.relation_factions_encounters
WHERE faction_id = $1;
