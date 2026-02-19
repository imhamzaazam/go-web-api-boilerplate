-- name: CreateSession :one
INSERT INTO "session" (
	"tenant_id"
	, "uid"
	, "user_email"
	, "refresh_token"
	, "user_agent"
	, "client_ip"
	, "is_blocked"
	, "expires_at"
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
	) RETURNING *;

-- name: GetSessionByTenant :one
SELECT *
FROM "session"
WHERE "tenant_id" = $1
	AND "uid" = $2 LIMIT 1;