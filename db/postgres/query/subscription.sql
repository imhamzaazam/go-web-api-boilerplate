-- name: CreateSubscription :one
INSERT INTO "subscription" (
    "tenant_id",
    "plan",
    "status",
    "starts_at",
    "ends_at"
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING "id", "tenant_id", "plan", "status", "starts_at", "ends_at", "created_at", "modified_at";
