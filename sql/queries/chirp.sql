-- name: CreateChirp :one
INSERT INTO
  chirps (id, created_at, updated_at, body, user_id)
VALUES
  (gen_random_uuid (), now(), now(), $1, $2)
RETURNING
  *;

-- name: GetAllChirps :many
SELECT
  *
FROM
  CHIRPS
ORDER BY
  created_at asc;

-- name: GetChirpById :one
SELECT
  *
FROM
  CHIRPS
WHERE
  ID = $1;

-- name: GetChirpsByAuthor :many
SELECT
  *
FROM
  chirps
WHERE
  user_id = $1
ORDER BY
  created_at asc;

-- name: DeleteChirp :exec
DELETE FROM CHIRPS
WHERE
  ID = $1;

-- name: IsChirpExists :one
SELECT
  EXISTS (
    SELECT
      id
    FROM
      CHIRPS
    WHERE
      ID = $1
  );
