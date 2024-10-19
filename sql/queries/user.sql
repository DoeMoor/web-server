-- name: CreateUser :one 
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid (),
    now(),
    now(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: UpdateUser :one
UPDATE users
SET updated_at = now(),
    email = $2,
    hashed_password = $3
WHERE id = $1
RETURNING *;

-- name: IsUserIdExists :one
SELECT EXISTS (
    SELECT id
    FROM users
    WHERE id = $1
);