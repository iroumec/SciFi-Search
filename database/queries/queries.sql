-- name: CreateUser :one
INSERT INTO users(name,surname) VALUES ($1, $2) RETURNING *;

-- name: ListUsers :many
SELECT * FROM users ORDER BY user_id;

-- name: GetUserByID :one
SELECT * FROM users WHERE user_id = $1;

-- name: UpdateUser :exec
UPDATE users SET name = $2, surname = $3 WHERE user_id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE user_id = $1;

-- name: CreatePreference :one
INSERT INTO preferences(preference) VALUES ($1) RETURNING *;

-- name: ListPreferences :many
SELECT preference FROM preferences ORDER BY preference;

-- name: DeletePreference :exec
DELETE FROM preferences WHERE preference = $1;

-- name: SetPreference :one
INSERT INTO user_preferences(user_id,preference) VALUES ($1,$2) RETURNING *;

-- name: RemovePreference :exec
DELETE FROM user_preferences WHERE user_id = $1 AND preference = $2;

-- name: ListPreferencesFromUser :many
SELECT preference FROM user_preferences WHERE user_id = $1;

-- name: CreateHistoricSearch :one
INSERT INTO historic_searches(user_id,search_string) VALUES ($1,$2) RETURNING *;

-- name: DeleteHistoricSearch :exec
DELETE FROM historic_searches WHERE historic_search_id = $1;

-- name: ListHistoricSearchesFromUser :many
SELECT search_string FROM historic_searches WHERE user_id = $1;