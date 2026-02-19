BEGIN;

DROP INDEX IF EXISTS "session_tenant_id_idx";
ALTER TABLE "session" DROP CONSTRAINT IF EXISTS "fk_session_user_tenant_email";
ALTER TABLE "session" DROP CONSTRAINT IF EXISTS "fk_session_tenant";
ALTER TABLE "session" DROP COLUMN IF EXISTS "tenant_id";

DROP INDEX IF EXISTS "user_tenant_id_idx";
DROP INDEX IF EXISTS "user_tenant_email_idx";
ALTER TABLE "user" DROP CONSTRAINT IF EXISTS "fk_user_tenant";
ALTER TABLE "user" DROP COLUMN IF EXISTS "tenant_id";

CREATE UNIQUE INDEX "user_email_idx" ON "user" USING BTREE ("email");

ALTER TABLE "session"
ADD CONSTRAINT "fk_user_email" FOREIGN KEY ("user_email") REFERENCES "user" ("email") ON DELETE CASCADE;

DROP TABLE IF EXISTS "tenant";

COMMIT;
