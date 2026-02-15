-- name: CreateBranch :one
INSERT INTO "branch" (
    "tenant_id",
    "name",
    "code"
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING "id", "tenant_id", "name", "code", "created_at", "modified_at";
