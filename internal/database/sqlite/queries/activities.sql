-- name: InsertActivity :one
INSERT INTO activities (activity_id, project_id, disp_name, duration)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateActivityNameByRowId :one
UPDATE activities
SET disp_name = ?
WHERE activity_id = (
    SELECT ao.successor_activity_id
    FROM activities_ordering ao
    WHERE ao.id = ?
      AND ao.project_id = activities.project_id
)
RETURNING *;