-- name: CreateCategory :one
INSERT INTO categories (
  category_name,
  section_name
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = $1 LIMIT 1;

-- name: ListCategories :many
SELECT * FROM categories
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateCategory :one
UPDATE categories
  set category_name = $2,
      section_name = $3
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;