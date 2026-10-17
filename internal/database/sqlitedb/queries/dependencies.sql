-- name: InsertDependency :one
INSERT INTO dependencies (id, project_id, relationship, predecessor_activity_id, successor_activity_id) 
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: FindAllPredecessorNamesBySuccessor :many
SELECT
    d.id AS dependency_id,
    d.relationship,
    predecessor.id AS predecessor_activity_id,
    predecessor.disp_name AS predecessor_activity_name
FROM dependencies d
JOIN activities predecessor
    ON predecessor.id = d.predecessor_activity_id
   AND predecessor.project_id = d.project_id
WHERE d.successor_activity_id = ?;