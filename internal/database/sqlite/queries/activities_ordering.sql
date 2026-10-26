-- name: InsertActivityOrdering :one
INSERT INTO activities_ordering (id, project_id, successor_activity_id, sort_rank)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: FindAllActivitiesByProjectSorted :many
SELECT a.* , ao.sort_rank, ao.id AS row_id
FROM activities a
JOIN activities_ordering ao
    ON a.activity_id = ao.successor_activity_id
    AND a.project_id = ao.project_id
WHERE a.project_id = ?
ORDER BY ao.sort_rank;

-- name: FindCurrentAndNextActivityByRowId :many
WITH current_activity AS (
    SELECT *
    FROM activities_ordering
    WHERE activities_ordering.id = ?
), project_activities AS (
    SELECT *
    FROM activities_ordering
    WHERE project_id =  (SELECT project_id FROM current_activity)
)
SELECT *
FROM project_activities
WHERE sort_rank > (SELECT sort_rank FROM current_activity)
   OR id = (SELECT id FROM current_activity)
ORDER BY sort_rank
LIMIT 2;