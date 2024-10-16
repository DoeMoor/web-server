-- name: CreateChirp :one
INSERT INTO chirps(id, created_at, updated_at, body, user_id)
VALUES (
  gen_random_uuid(),
  now(),
  now(),
  $1,
  $2
)
RETURNING *;

-- name: GetAllChirps :many
SELECT
	*
FROM
	CHIRPS
ORDER BY
	CREATED_AT ASC;