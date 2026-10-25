-- name: InsertDependency :one
INSERT INTO dependencies (id, project_id, relationship, predecessor_activity_id, successor_activity_id) 
VALUES (?, ?, ?, ?, ?)
RETURNING *;