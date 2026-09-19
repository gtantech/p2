-- name: InsertDependency :one
INSERT INTO dependencies (id, project_id, relationship, predecessor_activity_id, successor_activity_id) 
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteDependency :exec
DELETE FROM dependencies WHERE id = ?;

-- name: UpdateDependency :one
UPDATE dependencies
SET 
    relationship = ?,
    predecessor_activity_id = ?,
    successor_activity_id = ?
WHERE id = ?
RETURNING *;