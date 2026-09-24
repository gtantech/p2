-- name: FindAllDependenciesByProject :many
SELECT * FROM dependencies WHERE project_id = ?;

-- name: FindAllDependenciesByPredecessor :many
SELECT * FROM dependencies WHERE predecessor_activity_id = ?;

-- name: FindAllDependenciesBySuccessor :many
SELECT * FROM dependencies WHERE successor_activity_id = ?;

-- name: FindAllPredecessorNamesBySuccessor :many
SELECT
    d.id,
    d.relationship,
    predecessor.disp_name AS predecessor_activity_name
FROM dependencies d
JOIN activities predecessor
    ON predecessor.id = d.predecessor_activity_id
   AND predecessor.project_id = d.project_id
WHERE d.successor_activity_id = ?;


-- name: FindDependencyById :one
SELECT * FROM dependencies WHERE id = ?;

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