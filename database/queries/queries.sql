-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY username;

-- name: CreateUser :one
INSERT INTO users (username, name, email, password, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING *;

-- name: UpdateUser :exec
UPDATE users SET username = $2, email = $3 WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: CreateNew :one
INSERT INTO noticias (titulo, contenido, publicada_en, tiempo_lectura_estimado) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: ListNews: many
SELECT * FROM noticias ORDER BY publicada_en LIMIT 5 OFFSET $1;

-- name: CreateWork :one
INSERT INTO works (title, content_type_id, unit,saga_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: CreateContentType :one
INSERT INTO content_types (name) VALUES ($1) RETURNING *;

-- name: ReviewWork :one
INSERT INTO review (user_id, work_id, score, review, watched_at, liked) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetLastReview :one
SELECT * FROM review WHERE user_id = $1 AND work_id = $2;

-- name: UpdateReview :exec
UPDATE review SET score = $2, review = $3, watched_at = $4, liked = $5 WHERE id = $1;

-- name: DeleteReview :exec
DELETE FROM review WHERE id = $1;

-- name: ConsumeWork :one
INSERT INTO consumed_works (user_id,work_id) VALUES ($1,$2) RETURNING *;

-- name: UnconsumeWork :exec
DELETE FROM consumed_works WHERE user_id = $1 AND work_id = $2;

-- name: LikeWork :one
INSERT INTO liked_works (user_id,work_id) VALUES ($1,$2) RETURNING *;

-- name: LikeReview :one
INSERT INTO review_like (review_id,user_id) VALUES ($1,$2) RETURNING *;

-- name: UnlikeReview :exec
DELETE FROM review_like WHERE review_id = $1 AND user_id = $2;

-- name: CommentReview :one
INSERT INTO review_comment (review_id, user_id,comment) VALUES ($1,$2,$3) RETURNING *;

-- name: DeleteReviewComment :exec
DELETE FROM review_comment WHERE id = $1;

-- name: GetConsumedWorksByUser :many
SELECT * FROM works w WHERE w.id IN (SELECT id_work FROM consumed_works WHERE user_id = $1) ORDER BY (content_type_id,name);

-- name: FollowUser :one
INSERT INTO user_follows (follower_id, followed_id) VALUES ($1,$2) RETURNING *;

-- name: UnfollowUser :exec
DELETE FROM user_follows WHERE follower_id = $1 AND followed_id = $2;

-- name: GetNumberOfFollowers :one
SELECT COUNT(*) FROM user_follows WHERE followed_id = $1;

-- name: GetNumberOfFollowings :one
SELECT COUNT(*) FROM user_follows WHERE follower_id = $1;

-- name: AddWorkToFavourites :one
INSERT INTO user_favourites (user_id, work_id) VALUES ($1,$2) RETURNING *;

-- name: RemoveWorkFromFavourites :exec
DELETE FROM user_favourites WHERE user_id = $1 AND work_id = $2;

-- name: GetNumberOfFavouritesFromUser :one
SELECT COUNT(*) FROM user_favourites WHERE user_id = $1;

-- name: GetNumberOfFavouritesFromWork :one
SELECT COUNT(*) FROM user_favourites WHERE work_id = $1;