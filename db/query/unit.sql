-- name: CreateUnit :one
INSERT INTO units (
  unit_name,
  unit_value
) VALUES (
  $1, $2
) RETURNING *;

-- name: ListUnits :many
SELECT * FROM units
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUnit :one
UPDATE units
  set unit_name = $2,
      unit_value = $3
WHERE id = $1
RETURNING *;

-- name: DeleteUnit :exec
DELETE FROM units
WHERE id = $1;