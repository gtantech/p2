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

-- name: FindAllPredecessorNamesByProjectId :many
SELECT
    d.id AS dependency_id,
    d.relationship,
    d.successor_activity_id AS successor_activity_id,
    predecessor.id AS predecessor_activity_id,
    predecessor.disp_name AS predecessor_activity_name
FROM dependencies d
JOIN activities predecessor
    ON predecessor.id = d.predecessor_activity_id
   AND predecessor.project_id = d.project_id
WHERE d.project_id = ?;

-- name: FindAllDependenciesBySuccessorId :many
SELECT *
FROM dependencies
WHERE successor_activity_id = ?;

-- name: UpdatePredecessorById :one
UPDATE dependencies
SET 
    predecessor_activity_id = ?
WHERE id = ?
RETURNING *;

-- name: FindAllPredecessorActivitiesBySuccessorActivityId :many
SELECT
    d.id AS dependency_id,
    d.relationship,
    predecessor.*
FROM dependencies d
JOIN activities predecessor
    ON predecessor.id = d.predecessor_activity_id
   AND predecessor.project_id = d.project_id
WHERE d.successor_activity_id = ?;

-- name: DeleteDependencyById :exec
DELETE FROM dependencies
WHERE id = ?;