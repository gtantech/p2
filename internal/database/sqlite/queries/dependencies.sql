-- name: InsertDependency :one
INSERT INTO dependencies (id, project_id, relationship, predecessor_activity_id, successor_activity_id) 
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: FindAllPredecessorNamesByProjectId :many
SELECT
    d.id AS dependency_id,
    d.relationship,
    d.successor_activity_id AS successor_activity_id,
    predecessor.activity_id AS predecessor_activity_id,
    predecessor.disp_name AS predecessor_activity_name
FROM dependencies d
JOIN activities predecessor
    ON predecessor.activity_id = d.predecessor_activity_id
   AND predecessor.project_id = d.project_id
WHERE d.project_id = ?;