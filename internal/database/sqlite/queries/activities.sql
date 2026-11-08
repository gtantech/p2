-- name: InsertActivity :one
INSERT INTO activities (activity_id, project_id, disp_name, duration)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateActivityNameByRowId :one
UPDATE activities
SET disp_name = ?
WHERE id = (
    SELECT a.id
    FROM activities a
    JOIN activities_ordering ao
        ON a.activity_id = ao.successor_activity_id
        AND a.project_id = ao.project_id
    WHERE ao.id = ?
)
RETURNING *;