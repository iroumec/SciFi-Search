-- name: GetUser :one
SELECT * FROM Users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM Users WHERE username = $1;

-- name: ListUsers :many
SELECT * FROM Users ORDER BY username;

-- name: CreateUser :one
INSERT INTO Users (username, name, email, created_at) VALUES ($1, $2, $3, NOW()) RETURNING *;

-- name: UpdateUser :exec
UPDATE Users SET username = $2, email = $3 WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM Users WHERE id = $1;

-- name: CreateWork :one
INSERT INTO Works (title,content_type_id,unit,saga_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: CreateContentType :one
INSERT INTO ContentTypes (name) VALUES ($1) RETURNING *;

-- name: ReviewWork :one
INSERT INTO Review (user_id,work_id,score,review,when_watched,liked) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetLastReview :one
SELECT * FROM Review WHERE user_id = $1 AND work_id = $2;

-- name: UpdateReview :exec
UPDATE Review SET score = $2, review = $3, when_watched = $4, liked = $5 WHERE id = $1;

-- name: DeleteReview :exec
DELETE FROM Review WHERE id = $1;

-- name: ConsumeWork :one
INSERT INTO ConsumedWorks (user_id,work_id) VALUES ($1,$2) RETURNING *;

-- name: UnconsumeWork :exec
DELETE FROM ConsumedWorks WHERE user_id = $1 AND work_id = $2;

-- name: LikeWork :one
INSERT INTO LikedWorks (user_id,work_id) VALUES ($1,$2) RETURNING *;

-- name: LikeReview :one
INSERT INTO ReviewLike (review_id,user_id) VALUES ($1,$2) RETURNING *;

-- name: UnlikeReview :exec
DELETE FROM ReviewLike WHERE review_id = $1 AND user_id = $2;

-- name: CommentReview :one
INSERT INTO ReviewComment (review_id,user_id,comment) VALUES ($1,$2,$3) RETURNING *;

-- name: DeleteReviewComment :exec
DELETE FROM ReviewComment WHERE id = $1;

-- name: GetConsumedWorksByUser :many
SELECT * FROM Works w WHERE w.id IN (SELECT id_work FROM ConsumedWorks WHERE id_user = $1) ORDER BY (content_type_id,name);