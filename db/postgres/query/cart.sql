-- name: GetActiveCartByTenantAndUser :one
SELECT "id", "tenant_id", "user_uid", "is_active", "created_at", "modified_at"
FROM "cart"
WHERE "tenant_id" = $1
  AND "user_uid" = $2
  AND "is_active" = TRUE
LIMIT 1;

-- name: CreateCart :one
INSERT INTO "cart" (
    "tenant_id",
    "user_uid",
    "is_active"
)
VALUES (
    $1,
    $2,
    TRUE
)
RETURNING "id", "tenant_id", "user_uid", "is_active", "created_at", "modified_at";

-- name: CreateCartItem :one
INSERT INTO "cart_item" (
    "tenant_id",
    "cart_id",
    "product_id",
    "quantity",
    "unit_price",
    "vat_percent",
    "note",
    "prescription_ref"
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING "id", "tenant_id", "cart_id", "product_id", "quantity", "unit_price", "vat_percent", "note", "prescription_ref", "created_at", "modified_at";
