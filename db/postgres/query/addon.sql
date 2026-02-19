-- name: CreateProductAddon :one
INSERT INTO "product_addon" (
    "tenant_id",
    "product_id",
    "name",
    "price"
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING "id", "tenant_id", "product_id", "name", "price", "created_at", "modified_at";
