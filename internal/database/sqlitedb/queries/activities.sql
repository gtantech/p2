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