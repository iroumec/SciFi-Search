-- name: CreateUser :one
INSERT INTO users(username,email,name,middlename,surname,password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: ListUsers :many
SELECT * FROM users ORDER BY username;

-- name: UpdateUser :exec
UPDATE users SET email = $2, name = $3, middlename = $4, surname = $5, password = $6 WHERE username = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE username = $1;

-- name: CreatePreference :one
INSERT INTO preferences(preference) VALUES ($1) RETURNING *;

-- name: ListPreferences :many
SELECT preference FROM preferences ORDER BY preference;

-- name: DeletePreference :exec
DELETE FROM preferences WHERE preference = $1;

-- name: SetPreference :one
INSERT INTO user_preferences(username,preference) VALUES ($1,$2) RETURNING *;

-- name: RemovePreference :exec
DELETE FROM user_preferences WHERE username = $1 AND preference = $2;

-- name: ListPreferencesFromUser :many
SELECT preference FROM user_preferences WHERE username = $1;

-- name: CreateHistoricSearch :one
INSERT INTO historic_searches(username,search_string) VALUES ($1,$2) RETURNING *;

-- name: DeleteHistoricSearch :exec
DELETE FROM historic_searches WHERE historic_search_id = $1;

-- name: ListHistoricSearchesFromUser :many
SELECT search_string FROM historic_searches WHERE username = $1;