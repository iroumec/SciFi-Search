-- name: GetUser :one
SELECT * FROM Users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM Users ORDER BY username;

-- name: CreateUser :one
INSERT INTO Users (username, name, email, created_at) VALUES ($1, $2, $3, NOW()) RETURNING *;

-- name: UpdateUser :exec
UPDATE Users SET username = $2, email = $3 WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM Users WHERE id = $1;

-- name: CreateWork :exec
INSERT INTO Works (title,content_type_id,unit,saga_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: ReviewWork :exec
INSERT INTO Review (user_id,work_id,score,review,when_watched,liked) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;
