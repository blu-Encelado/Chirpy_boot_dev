CREATE EXTENSION pgcrypto;

-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: ResetUser :exec
DELETE FROM users;

-- name: GetUser :one
SELECT id, created_at, updated_at, email, is_chirpy_red FROM users WHERE email = $1;

-- name: UpdateEmailAndPass :exec
UPDATE users SET email = $1, hashed_password = $2 WHERE id = $3;

-- name: UpdateIsChirpyRed :exec
UPDATE users SET is_chirpy_red = $1 WHERE id = $2;