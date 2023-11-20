-- name: CreateGood :one
INSERT INTO goods (
  category,
  model,
  unit,
  amount,
  good_desc
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetGood :one
SELECT * FROM goods
WHERE id = $1 LIMIT 1;

-- name: ListGoods :many
SELECT * FROM goods
WHERE 
    category = $1 OR
    model = $2
ORDER BY id
LIMIT $3
OFFSET $4;

-- name: UpdateGood :one
UPDATE goods
  set unit = $2,
      amount = $3
WHERE id = $1
RETURNING *;

-- name: DeleteGood :exec
DELETE FROM goods
WHERE id = $1;