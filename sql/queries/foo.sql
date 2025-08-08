-- name: CreateFoo :one
INSERT INTO foos (
  message
) VALUES (
  $1
)
RETURNING *;

-- name: GetFoo :one
SELECT * FROM foos
WHERE id = $1 LIMIT 1;
