BEGIN;

DROP INDEX IF EXISTS "order_item_order_idx";
DROP INDEX IF EXISTS "order_tenant_status_idx";

DROP TABLE IF EXISTS "order_item";
DROP TABLE IF EXISTS "order";
DROP INDEX IF EXISTS "payment_method_default_per_tenant_idx";
DROP TABLE IF EXISTS "payment_method";
DROP INDEX IF EXISTS "inventory_tenant_product_idx";
DROP TABLE IF EXISTS "inventory";
DROP TABLE IF EXISTS "cart_item_addon";
DROP TABLE IF EXISTS "cart_item";
DROP INDEX IF EXISTS "cart_tenant_user_active_idx";
DROP TABLE IF EXISTS "cart";
DROP TABLE IF EXISTS "product_addon";
DROP INDEX IF EXISTS "discount_tenant_code_idx";
DROP TABLE IF EXISTS "discount";
DROP INDEX IF EXISTS "product_tenant_sku_idx";
DROP TABLE IF EXISTS "product";

ALTER TABLE "user" DROP CONSTRAINT IF EXISTS "fk_user_branch";
ALTER TABLE "user"
    DROP COLUMN IF EXISTS "deleted_at",
    DROP COLUMN IF EXISTS "branch_id",
    DROP COLUMN IF EXISTS "role";

DROP INDEX IF EXISTS "subscription_one_open_per_tenant_idx";
DROP TABLE IF EXISTS "subscription";
DROP INDEX IF EXISTS "branch_tenant_code_idx";
DROP TABLE IF EXISTS "branch";

DROP INDEX IF EXISTS "tenant_domain_idx";
ALTER TABLE "tenant"
    DROP COLUMN IF EXISTS "deleted_at",
    DROP COLUMN IF EXISTS "type",
    DROP COLUMN IF EXISTS "domain";

DROP TYPE IF EXISTS payment_method_type;
DROP TYPE IF EXISTS discount_type;
DROP TYPE IF EXISTS fulfillment_type;
DROP TYPE IF EXISTS order_status;
DROP TYPE IF EXISTS user_role;
DROP TYPE IF EXISTS subscription_status;
DROP TYPE IF EXISTS tenant_type;

COMMIT;
