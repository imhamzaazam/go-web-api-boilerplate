-- name: UpsertInventoryForProduct :one
INSERT INTO "inventory" (
    "tenant_id",
    "product_id",
    "in_stock",
    "reserved"
)
VALUES (
    $1,
    $2,
    $3,
    0
)
ON CONFLICT ("tenant_id", "product_id")
DO UPDATE SET
    "in_stock" = EXCLUDED."in_stock",
    "modified_at" = CURRENT_TIMESTAMP
RETURNING "id", "tenant_id", "product_id", "in_stock", "reserved", "created_at", "modified_at";

-- name: ReserveInventory :one
UPDATE "inventory"
SET
    "reserved" = "reserved" + $3,
    "modified_at" = CURRENT_TIMESTAMP
WHERE "tenant_id" = $1
  AND "product_id" = $2
  AND ("in_stock" - "reserved") >= $3
RETURNING "id", "tenant_id", "product_id", "in_stock", "reserved", "created_at", "modified_at";

-- name: GetInventoryByTenantAndProduct :one
SELECT "id", "tenant_id", "product_id", "in_stock", "reserved", "created_at", "modified_at"
FROM "inventory"
WHERE "tenant_id" = $1
  AND "product_id" = $2
LIMIT 1;
