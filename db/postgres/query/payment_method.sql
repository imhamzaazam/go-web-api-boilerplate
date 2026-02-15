-- name: EnsureDefaultCashPaymentMethod :exec
INSERT INTO "payment_method" (
    "tenant_id",
    "type",
    "label",
    "is_default"
)
VALUES (
    $1,
    'cash',
    'Cash',
    TRUE
)
ON CONFLICT DO NOTHING;

-- name: GetDefaultCashPaymentMethodByTenant :one
SELECT "id", "tenant_id", "type", "label", "is_default", "created_at", "modified_at"
FROM "payment_method"
WHERE "tenant_id" = $1
  AND "type" = 'cash'
  AND "deleted_at" IS NULL
ORDER BY "is_default" DESC, "created_at" ASC
LIMIT 1;

-- name: CreatePaymentMethod :one
INSERT INTO "payment_method" (
    "tenant_id",
    "type",
    "label",
    "is_default"
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING "id", "tenant_id", "type", "label", "is_default", "created_at", "modified_at";

-- name: ListPaymentMethodsByTenant :many
SELECT "id", "tenant_id", "type", "label", "is_default", "created_at", "modified_at"
FROM "payment_method"
WHERE "tenant_id" = $1
  AND "deleted_at" IS NULL
ORDER BY "is_default" DESC, "created_at" ASC;

-- name: GetPaymentMethodByTenantAndID :one
SELECT "id", "tenant_id", "type", "label", "is_default", "created_at", "modified_at"
FROM "payment_method"
WHERE "tenant_id" = $1
  AND "id" = $2
  AND "deleted_at" IS NULL
LIMIT 1;
