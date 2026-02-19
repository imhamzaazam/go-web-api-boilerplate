-- name: CreateProduct :one
INSERT INTO "product" (
    "tenant_id",
    "name",
    "sku",
    "price",
    "vat_percent",
    "is_preorder",
    "made_to_order",
    "requires_prescription",
    "available_for_delivery",
    "available_for_pickup"
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10
)
RETURNING "id", "tenant_id", "name", "sku", "price", "vat_percent", "is_preorder", "made_to_order", "requires_prescription", "available_for_delivery", "available_for_pickup", "created_at", "modified_at";

-- name: GetProductByTenantAndID :one
SELECT "id", "tenant_id", "name", "sku", "price", "vat_percent", "is_preorder", "made_to_order", "requires_prescription", "available_for_delivery", "available_for_pickup", "created_at", "modified_at"
FROM "product"
WHERE "tenant_id" = $1
  AND "id" = $2
  AND "deleted_at" IS NULL
LIMIT 1;
