-- name: CreateDiscount :one
INSERT INTO "discount" (
    "tenant_id",
    "product_id",
    "code",
    "name",
    "type",
    "value",
    "starts_at",
    "ends_at"
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
RETURNING "id", "tenant_id", "product_id", "code", "name", "type", "value", "starts_at", "ends_at", "created_at", "modified_at";
