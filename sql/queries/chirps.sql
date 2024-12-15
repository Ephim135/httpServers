-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetChirpById :one
SELECT * 
FROM chirps
WHERE id = $1;

-- name: GetChirpsDesc :many
SELECT * 
FROM chirps 
ORDER BY created_at DESC;

-- name: GetChirps :many
SELECT * 
FROM chirps;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;

-- name: GetChirpsByAuthor :many
SELECT *
FROM chirps
WHERE user_id = $1;

-- name: GetChirpsByAuthorDesc :many
SELECT *
FROM chirps
WHERE user_id = $1
ORDER BY created_at DESC;