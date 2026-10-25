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