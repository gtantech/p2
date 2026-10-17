-- name: InsertActivityOrdering :one
INSERT INTO activities_ordering (id, project_id, successor_activity_id, sort_rank)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: FindLastActivityByProjectId :one
SELECT *
FROM activities_ordering
WHERE project_id = ?
ORDER BY sort_rank
LIMIT 1;