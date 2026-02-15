-- name: GetTenantByDomain :one
SELECT "id", "name", "slug", "domain", "type", "created_at", "modified_at"
FROM "tenant"
WHERE "domain" = $1
	AND "deleted_at" IS NULL
LIMIT 1;

-- name: CreateTenant :one
INSERT INTO "tenant" (
	"name",
	"slug",
	"domain",
	"type"
)
VALUES (
	$1,
	$2,
	$3,
	$4
)
RETURNING "id", "name", "slug", "domain", "type", "created_at", "modified_at";

-- name: GetTenantBySlug :one
SELECT "id", "name", "slug", "domain", "type", "created_at", "modified_at"
FROM "tenant"
WHERE "slug" = $1
	AND "deleted_at" IS NULL
LIMIT 1;

-- name: GetTenantByID :one
SELECT "id", "name", "slug", "domain", "type", "created_at", "modified_at"
FROM "tenant"
WHERE "id" = $1
	AND "deleted_at" IS NULL
LIMIT 1;

-- name: GetLatestSubscriptionByTenant :one
SELECT "id", "tenant_id", "plan", "status", "starts_at", "ends_at", "created_at", "modified_at"
FROM "subscription"
WHERE "tenant_id" = $1
	AND "deleted_at" IS NULL
ORDER BY "created_at" DESC
LIMIT 1;
