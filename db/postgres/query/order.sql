-- name: GetCartByTenantAndID :one
SELECT "id", "tenant_id", "user_uid", "is_active", "created_at", "modified_at"
FROM "cart"
WHERE "tenant_id" = $1
  AND "id" = $2
  AND "is_active" = TRUE
LIMIT 1;

-- name: CountCartItemsByTenantAndCart :one
SELECT COUNT(*)::INT
FROM "cart_item"
WHERE "tenant_id" = $1
  AND "cart_id" = $2
  AND "deleted_at" IS NULL;

-- name: ListCartItemsByTenantAndCart :many
SELECT "id", "tenant_id", "cart_id", "product_id", "quantity", "unit_price", "vat_percent", "note", "prescription_ref", "created_at", "modified_at"
FROM "cart_item"
WHERE "tenant_id" = $1
  AND "cart_id" = $2
  AND "deleted_at" IS NULL
ORDER BY "created_at" ASC;

-- name: CreateOrder :one
INSERT INTO "order" (
    "tenant_id",
    "user_uid",
    "cart_id",
    "status",
    "fulfillment_type",
    "delivery_address_line",
    "delivery_city",
    "delivery_lat",
    "delivery_lng",
    "subtotal",
    "tax",
    "total"
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
    $10,
    $11,
    $12
)
RETURNING "id", "tenant_id", "user_uid", "cart_id", "status", "fulfillment_type", "subtotal", "tax", "total", "created_at", "modified_at";

-- name: CreateOrderItem :one
INSERT INTO "order_item" (
    "tenant_id",
    "order_id",
    "product_id",
    "quantity",
    "unit_price",
    "vat_percent",
    "line_total",
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
    $8,
    $9
)
RETURNING "id", "tenant_id", "order_id", "product_id", "quantity", "unit_price", "vat_percent", "line_total", "note", "prescription_ref", "created_at";

-- name: SetCartInactive :exec
UPDATE "cart"
SET "is_active" = FALSE,
    "modified_at" = CURRENT_TIMESTAMP
WHERE "tenant_id" = $1
  AND "id" = $2;

-- name: GetOrderByTenantAndID :one
SELECT "id", "tenant_id", "user_uid", "cart_id", "status", "fulfillment_type", "subtotal", "tax", "total", ("paid_at" IS NOT NULL)::bool AS "is_paid", "created_at", "modified_at"
FROM "order"
WHERE "tenant_id" = $1
  AND "id" = $2
  AND "deleted_at" IS NULL
LIMIT 1;

-- name: UpdateOrderStatusByTenant :one
UPDATE "order"
SET "status" = $3,
    "cancelled_at" = CASE WHEN $4::BOOL THEN CURRENT_TIMESTAMP ELSE "cancelled_at" END,
    "refunded_at" = CASE WHEN $5::BOOL THEN CURRENT_TIMESTAMP ELSE "refunded_at" END,
    "modified_at" = CURRENT_TIMESTAMP
WHERE "tenant_id" = $1
  AND "id" = $2
  AND "deleted_at" IS NULL
RETURNING "id", "tenant_id", "user_uid", "cart_id", "status", "fulfillment_type", "subtotal", "tax", "total", ("paid_at" IS NOT NULL)::bool AS "is_paid", "created_at", "modified_at";

-- name: MarkOrderPaidByTenant :one
UPDATE "order"
SET "payment_method_id" = $3,
    "paid_at" = CURRENT_TIMESTAMP,
    "modified_at" = CURRENT_TIMESTAMP
WHERE "tenant_id" = $1
  AND "id" = $2
  AND "deleted_at" IS NULL
  AND "paid_at" IS NULL
RETURNING "id", "tenant_id", "user_uid", "cart_id", "status", "fulfillment_type", "subtotal", "tax", "total", TRUE AS "is_paid", "created_at", "modified_at";
