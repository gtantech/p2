-- name: FindProjectById :one
SELECT * FROM projects WHERE id = ?;

-- name: InsertProject :one
INSERT INTO projects (id, disp_name) VALUES (?, ?)
RETURNING *;