-- name: FindAllActivitiesByProject :many
SELECT * FROM activities WHERE project_id = ?;

-- name: FindAllActivitiesByNameAndProject :many
SELECT * FROM activities WHERE disp_name LIKE ? AND project_id = ?;

-- name: FindActivityById :one
SELECT * FROM activities WHERE id = ?;

-- name: InsertActivity :one
INSERT INTO activities (id, project_id, disp_name, duration)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: DeleteActivity :exec
DELETE FROM activities WHERE id = ?;

-- name: UpdateActivity :one
UPDATE activities
SET 
    disp_name = ?,
    duration = ?
WHERE id = ?
RETURNING *;