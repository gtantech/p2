-- name: FindAllActivitiesByProject :many
SELECT * FROM activities WHERE project_id = ?;

-- name: FindAllActivitiesByProjectAndClosestName :many
SELECT * FROM activities WHERE project_id = ? AND disp_name LIKE ?;

-- name: FindAllActivitiesByProjectAndName :many
SELECT * FROM activities WHERE project_id = ? AND disp_name = ?;

-- name: FindActivityById :one
SELECT * FROM activities WHERE id = ?;

-- name: InsertActivity :one
INSERT INTO activities (id, project_id, disp_name, duration)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: FindAllActivitiesByProjectSorted :many
SELECT a.* , ao.sort_rank
FROM activities a
JOIN activities_ordering ao
    ON a.id = ao.successor_activity_id
    AND a.project_id = ao.project_id
WHERE a.project_id = ?
ORDER BY ao.sort_rank;