-- name: InsertActivityOrdering :one
INSERT INTO activities_ordering (id, project_id, successor_activity_id, sort_rank)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateActivityOrdering :one
UPDATE activities_ordering
SET 
    sort_rank = ?
WHERE successor_activity_id = ?
RETURNING *;