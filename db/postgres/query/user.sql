-- name: CreateUser :one
INSERT INTO "user" (
	"tenant_id"
	, "role"
	, "email"
	, "password"
	, "full_name"
	, "is_staff"
	, "is_active"
	, "last_login"
	)
VALUES (
	$1
	, $2
	, $3
	, $4
	, $5
	, $6
	, $7
	, $8
	) RETURNING "uid"
	, "tenant_id"
	, "email"
	, "full_name"
	, "created_at"
	, "modified_at";

-- name: GetUserByTenantAndEmail :one
SELECT
	"id",
	"tenant_id",
	"role",
	"uid",
	"email",
	"password",
	"full_name",
	"is_staff",
	"is_active",
	COALESCE("last_login", CURRENT_TIMESTAMP) AS "last_login",
	"created_at",
	"modified_at"
FROM "user"
WHERE "tenant_id" = $1
	AND "email" = $2
	AND "deleted_at" IS NULL
LIMIT 1;

-- name: GetUserByTenantAndUID :one
SELECT
	"id",
	"tenant_id",
	"role",
	"uid",
	"email",
	"password",
	"full_name",
	"is_staff",
	"is_active",
	COALESCE("last_login", CURRENT_TIMESTAMP) AS "last_login",
	"created_at",
	"modified_at"
FROM "user"
WHERE "tenant_id" = $1
	AND "uid" = $2
	AND "deleted_at" IS NULL
LIMIT 1;
