-- name: CreateUser :one
INSERT INTO users(name, surname, auth_id) VALUES ($1, $2, $3) RETURNING *;

-- name: ListUsers :many
SELECT * FROM users ORDER BY user_id;

-- name: GetUserByID :one
SELECT * FROM users WHERE user_id = $1;

-- name: GetUserByAuthID :one
SELECT * FROM users WHERE auth_id = $1;

-- name: UploadAvatar :exec
UPDATE users SET avatar_url = $2 WHERE user_id = $1;

-- name: DeleteAvatar :exec
UPDATE users SET avatar_url = NULL WHERE user_id = $1;

-- name: UpdateUser :exec
UPDATE users SET name = $2, surname = $3 WHERE user_id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE user_id = $1;

-- name: AddPreference :one
INSERT INTO user_preferences(user_id,preference) VALUES ($1,$2) RETURNING *;

-- name: RemovePreference :exec
DELETE FROM user_preferences WHERE user_id = $1 AND preference = $2;

-- name: RemoveAllPreferenceFromUser :exec
DELETE FROM user_preferences WHERE user_id = $1;

-- name: ListPreferencesFromUser :many
SELECT preference FROM user_preferences WHERE user_id = $1;

-- name: CreateHistoricSearch :one
INSERT INTO historic_searches(user_id,search_string) VALUES ($1,$2) RETURNING *;

-- name: DeleteHistoricSearch :exec
DELETE FROM historic_searches WHERE historic_search_id = $1;

-- name: ListHistoricSearchesFromUser :many
SELECT search_string, search_datetime FROM historic_searches WHERE user_id = $1 ORDER BY search_datetime DESC;

-- name: GetTrendingSearches :many
SELECT search_string, COUNT(*) AS count
FROM historic_searches
-- Búquedas realizadas hoy según UTC.
WHERE search_datetime >= date_trunc('day', NOW() AT TIME ZONE 'UTC')
    AND search_datetime <  date_trunc('day', NOW() AT TIME ZONE 'UTC') + INTERVAL '1 day'
GROUP BY search_string
ORDER BY count DESC
LIMIT $1;

-- name: AddDocument :one
INSERT INTO documents(name, description, first_area, second_area, type, link) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;