-- name: InsertActivityOrdering :one
INSERT INTO activities_ordering (id, project_id, successor_activity_id, sort_rank)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: FindLastActivityByProjectId :one
SELECT *
FROM activities_ordering
WHERE project_id = ?
ORDER BY sort_rank
LIMIT 1;

-- name: FindNextActivityByProjectIdAndActivityId :one
WITH project_activities AS (
    SELECT *
    FROM activities_ordering
    WHERE project_id = ?
),
current_activity AS (
    SELECT sort_rank
    FROM project_activities
    WHERE successor_activity_id = ?
)
SELECT *
FROM project_activities
WHERE sort_rank > (SELECT sort_rank FROM current_activity)
ORDER BY sort_rank
LIMIT 1;