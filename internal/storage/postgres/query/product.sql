-- name: CreateProduct :one
INSERT INTO products (
    name,
    description,
    price,
    quantity
) VALUES (
    $1, $2, $3, $4
)
RETURNING id;

-- name: GetProductByID :one
SELECT id, name, description, price, quantity
FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT id, name, description, price, quantity
FROM products
ORDER BY id;

-- name: UpdateProduct :exec
UPDATE products
SET
    name = $2,
    description = $3,
    price = $4,
    quantity = $5
WHERE id = $1;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;
