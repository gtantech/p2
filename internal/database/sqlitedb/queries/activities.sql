-- name: InsertActivity :one
INSERT INTO activities (id, project_id, disp_name, duration)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: FindActivityById :one
SELECT * FROM activities WHERE id = ?;

-- name: FindAllActivitiesByProject :many
SELECT * FROM activities WHERE project_id = ?;