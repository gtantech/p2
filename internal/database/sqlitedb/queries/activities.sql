-- name: InsertActivity :one
INSERT INTO activities (id, project_id, disp_name, duration)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: FindActivityById :one
SELECT * FROM activities WHERE id = ?;

-- name: FindAllActivitiesByProject :many
SELECT * FROM activities WHERE project_id = ?;

-- name: FindAllActivitiesByProjectSorted :many
SELECT a.* , ao.sort_rank
FROM activities a
JOIN activities_ordering ao
    ON a.id = ao.successor_activity_id
    AND a.project_id = ao.project_id
WHERE a.project_id = ?
ORDER BY ao.sort_rank;

-- name: UpdateActivityDuration :one
UPDATE activities
SET 
    duration = ?
WHERE id = ?
RETURNING *;