-- name: InsertProject :one
INSERT INTO projects (id, disp_name) VALUES (?, ?)
RETURNING *;