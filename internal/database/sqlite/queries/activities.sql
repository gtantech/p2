-- name: InsertActivity :one
INSERT INTO activities (activity_id, project_id, disp_name, duration)
VALUES (?, ?, ?, ?)
RETURNING *;