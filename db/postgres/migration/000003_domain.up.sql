BEGIN;

CREATE TYPE tenant_type AS ENUM ('bakery', 'pharmacy', 'restaurant');
CREATE TYPE subscription_status AS ENUM ('trial', 'active', 'past_due', 'suspended', 'canceled');
CREATE TYPE user_role AS ENUM ('owner', 'admin', 'employee');
CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'cancelled', 'out_for_delivery', 'completed', 'refunded');
CREATE TYPE fulfillment_type AS ENUM ('pickup', 'delivery');
CREATE TYPE discount_type AS ENUM ('percentage');
CREATE TYPE payment_method_type AS ENUM ('cash', 'card');

CREATE TABLE IF NOT EXISTS "tenant" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" VARCHAR(255) NOT NULL,
    "slug" VARCHAR(255) NOT NULL UNIQUE,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modified_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO "tenant" ("name", "slug")
SELECT 'Default Tenant', 'default'
WHERE NOT EXISTS (SELECT 1 FROM "tenant" WHERE "slug" = 'default');

ALTER TABLE "user" ADD COLUMN IF NOT EXISTS "tenant_id" UUID;
ALTER TABLE "session" ADD COLUMN IF NOT EXISTS "tenant_id" UUID;

UPDATE "user"
SET "tenant_id" = (SELECT "id" FROM "tenant" WHERE "slug" = 'default')
WHERE "tenant_id" IS NULL;

UPDATE "session"
SET "tenant_id" = (SELECT "id" FROM "tenant" WHERE "slug" = 'default')
WHERE "tenant_id" IS NULL;

ALTER TABLE "user" ALTER COLUMN "tenant_id" SET NOT NULL;
ALTER TABLE "session" ALTER COLUMN "tenant_id" SET NOT NULL;

ALTER TABLE "user" DROP CONSTRAINT IF EXISTS "fk_user_tenant";
ALTER TABLE "user"
    ADD CONSTRAINT "fk_user_tenant"
    FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id") ON DELETE RESTRICT;

ALTER TABLE "session" DROP CONSTRAINT IF EXISTS "fk_session_tenant";
ALTER TABLE "session"
    ADD CONSTRAINT "fk_session_tenant"
    FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id") ON DELETE RESTRICT;

ALTER TABLE "session" DROP CONSTRAINT IF EXISTS "fk_user_email";
ALTER TABLE "session" DROP CONSTRAINT IF EXISTS "fk_session_user_tenant_email";

DROP INDEX IF EXISTS "user_email_idx";
CREATE UNIQUE INDEX IF NOT EXISTS "user_tenant_email_idx" ON "user" ("tenant_id", "email");
CREATE INDEX IF NOT EXISTS "user_tenant_id_idx" ON "user" ("tenant_id");

ALTER TABLE "session"
    ADD CONSTRAINT "fk_session_user_tenant_email"
    FOREIGN KEY ("tenant_id", "user_email")
    REFERENCES "user" ("tenant_id", "email") ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS "session_tenant_id_idx" ON "session" ("tenant_id");

ALTER TABLE "tenant"
    ADD COLUMN "domain" VARCHAR(255),
    ADD COLUMN "type" tenant_type,
    ADD COLUMN "deleted_at" TIMESTAMPTZ;

UPDATE "tenant"
SET "domain" = "slug" || '.localhost',
    "type" = 'bakery';

ALTER TABLE "tenant"
    ALTER COLUMN "domain" SET NOT NULL,
    ALTER COLUMN "type" SET NOT NULL;

CREATE UNIQUE INDEX "tenant_domain_idx" ON "tenant" USING BTREE ("domain");

CREATE TABLE "branch" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "tenant_id" UUID NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "code" VARCHAR(64) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modified_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMPTZ,
    CONSTRAINT "fk_branch_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id") ON DELETE CASCADE
);

CREATE UNIQUE INDEX "branch_tenant_code_idx" ON "branch" ("tenant_id", "code") WHERE "deleted_at" IS NULL;

CREATE TABLE "subscription" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "tenant_id" UUID NOT NULL,
    "plan" VARCHAR(100) NOT NULL,
    "status" subscription_status NOT NULL,
    "starts_at" TIMESTAMPTZ NOT NULL,
    "ends_at" TIMESTAMPTZ NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modified_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMPTZ,
    CONSTRAINT "fk_subscription_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id") ON DELETE CASCADE,
    CONSTRAINT "chk_subscription_dates" CHECK ("starts_at" < "ends_at")
);

CREATE UNIQUE INDEX "subscription_one_open_per_tenant_idx"
    ON "subscription" ("tenant_id")
    WHERE "deleted_at" IS NULL AND "status" IN ('trial', 'active');

ALTER TABLE "user"
    ADD COLUMN "role" user_role NOT NULL DEFAULT 'employee',
    ADD COLUMN "branch_id" UUID,
    ADD COLUMN "deleted_at" TIMESTAMPTZ;

ALTER TABLE "user"
    ADD CONSTRAINT "fk_user_branch" FOREIGN KEY ("branch_id") REFERENCES "branch" ("id") ON DELETE SET NULL;

CREATE TABLE "product" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "tenant_id" UUID NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "sku" VARCHAR(100) NOT NULL,
    "price" BIGINT NOT NULL,
    "vat_percent" NUMERIC(5,2) NOT NULL,
    "is_preorder" BOOLEAN NOT NULL DEFAULT FALSE,
    "made_to_order" BOOLEAN NOT NULL DEFAULT FALSE,
    "requires_prescription" BOOLEAN NOT NULL DEFAULT FALSE,
    "available_for_delivery" BOOLEAN NOT NULL DEFAULT TRUE,
    "available_for_pickup" BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modified_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMPTZ,
    CONSTRAINT "fk_product_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id") ON DELETE CASCADE,
    CONSTRAINT "chk_product_price" CHECK ("price" >= 0),
    CONSTRAINT "chk_product_vat_percent" CHECK ("vat_percent" >= 0 AND "vat_percent" <= 100)
);

CREATE UNIQUE INDEX "product_tenant_sku_idx" ON "product" ("tenant_id", "sku") WHERE "deleted_at" IS NULL;

CREATE TABLE "discount" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "tenant_id" UUID NOT NULL,
    "product_id" UUID NOT NULL,
    "code" VARCHAR(50) NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "type" discount_type NOT NULL,
    "value" NUMERIC(5,2) NOT NULL,
    "starts_at" TIMESTAMPTZ NOT NULL,
    "ends_at" TIMESTAMPTZ NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modified_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMPTZ,
    CONSTRAINT "fk_discount_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_discount_product" FOREIGN KEY ("product_id") REFERENCES "product" ("id") ON DELETE CASCADE,
    CONSTRAINT "chk_discount_value" CHECK ("value" >= 1 AND "value" <= 100),
    CONSTRAINT "chk_discount_dates" CHECK ("starts_at" < "ends_at")
);

CREATE UNIQUE INDEX "discount_tenant_code_idx" ON "discount" ("tenant_id", "code") WHERE "deleted_at" IS NULL;

CREATE TABLE "product_addon" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "tenant_id" UUID NOT NULL,
    "product_id" UUID NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "price" BIGINT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modified_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMPTZ,
    CONSTRAINT "fk_addon_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_addon_product" FOREIGN KEY ("product_id") REFERENCES "product" ("id") ON DELETE CASCADE,
    CONSTRAINT "chk_addon_price" CHECK ("price" >= 0)
);

CREATE TABLE "cart" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "tenant_id" UUID NOT NULL,
    "user_uid" UUID NOT NULL,
    "is_active" BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modified_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "fk_cart_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_cart_user_uid" FOREIGN KEY ("user_uid") REFERENCES "user" ("uid") ON DELETE CASCADE
);

CREATE UNIQUE INDEX "cart_tenant_user_active_idx" ON "cart" ("tenant_id", "user_uid") WHERE "is_active" = TRUE;

CREATE TABLE "cart_item" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "tenant_id" UUID NOT NULL,
    "cart_id" UUID NOT NULL,
    "product_id" UUID NOT NULL,
    "quantity" INTEGER NOT NULL,
    "unit_price" BIGINT NOT NULL,
    "vat_percent" NUMERIC(5,2) NOT NULL,
    "note" VARCHAR(500),
    "prescription_ref" VARCHAR(255),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modified_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMPTZ,
    CONSTRAINT "fk_cart_item_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_cart_item_cart" FOREIGN KEY ("cart_id") REFERENCES "cart" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_cart_item_product" FOREIGN KEY ("product_id") REFERENCES "product" ("id") ON DELETE RESTRICT,
    CONSTRAINT "chk_cart_item_quantity" CHECK ("quantity" > 0),
    CONSTRAINT "chk_cart_item_unit_price" CHECK ("unit_price" >= 0),
    CONSTRAINT "chk_cart_item_vat_percent" CHECK ("vat_percent" >= 0 AND "vat_percent" <= 100)
);

CREATE TABLE "cart_item_addon" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "tenant_id" UUID NOT NULL,
    "cart_item_id" UUID NOT NULL,
    "addon_id" UUID NOT NULL,
    "quantity" INTEGER NOT NULL DEFAULT 1,
    CONSTRAINT "fk_cart_item_addon_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_cart_item_addon_cart_item" FOREIGN KEY ("cart_item_id") REFERENCES "cart_item" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_cart_item_addon_addon" FOREIGN KEY ("addon_id") REFERENCES "product_addon" ("id") ON DELETE RESTRICT,
    CONSTRAINT "chk_cart_item_addon_quantity" CHECK ("quantity" > 0)
);

CREATE TABLE "inventory" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "tenant_id" UUID NOT NULL,
    "product_id" UUID NOT NULL,
    "in_stock" INTEGER NOT NULL DEFAULT 0,
    "reserved" INTEGER NOT NULL DEFAULT 0,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modified_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "fk_inventory_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_inventory_product" FOREIGN KEY ("product_id") REFERENCES "product" ("id") ON DELETE CASCADE,
    CONSTRAINT "chk_inventory_in_stock" CHECK ("in_stock" >= 0),
    CONSTRAINT "chk_inventory_reserved" CHECK ("reserved" >= 0)
);

CREATE UNIQUE INDEX "inventory_tenant_product_idx" ON "inventory" ("tenant_id", "product_id");

CREATE TABLE "payment_method" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "tenant_id" UUID NOT NULL,
    "type" payment_method_type NOT NULL,
    "label" VARCHAR(255) NOT NULL,
    "is_default" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modified_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMPTZ,
    CONSTRAINT "fk_payment_method_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id") ON DELETE CASCADE
);

CREATE UNIQUE INDEX "payment_method_default_per_tenant_idx"
    ON "payment_method" ("tenant_id")
    WHERE "is_default" = TRUE AND "deleted_at" IS NULL;

CREATE TABLE "order" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "tenant_id" UUID NOT NULL,
    "user_uid" UUID NOT NULL,
    "branch_id" UUID,
    "cart_id" UUID,
    "status" order_status NOT NULL DEFAULT 'pending',
    "fulfillment_type" fulfillment_type NOT NULL DEFAULT 'pickup',
    "delivery_address_line" VARCHAR(255),
    "delivery_city" VARCHAR(100),
    "delivery_lat" NUMERIC(10,7),
    "delivery_lng" NUMERIC(10,7),
    "subtotal" BIGINT NOT NULL,
    "tax" BIGINT NOT NULL,
    "total" BIGINT NOT NULL,
    "payment_method_id" UUID,
    "paid_at" TIMESTAMPTZ,
    "cancelled_at" TIMESTAMPTZ,
    "refunded_at" TIMESTAMPTZ,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "modified_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMPTZ,
    CONSTRAINT "fk_order_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_order_user_uid" FOREIGN KEY ("user_uid") REFERENCES "user" ("uid") ON DELETE RESTRICT,
    CONSTRAINT "fk_order_branch" FOREIGN KEY ("branch_id") REFERENCES "branch" ("id") ON DELETE SET NULL,
    CONSTRAINT "fk_order_cart" FOREIGN KEY ("cart_id") REFERENCES "cart" ("id") ON DELETE SET NULL,
    CONSTRAINT "fk_order_payment_method" FOREIGN KEY ("payment_method_id") REFERENCES "payment_method" ("id") ON DELETE SET NULL,
    CONSTRAINT "chk_order_delivery_location"
        CHECK (
            ("fulfillment_type" = 'pickup')
            OR (
                "fulfillment_type" = 'delivery'
                AND "delivery_address_line" IS NOT NULL
                AND "delivery_city" IS NOT NULL
                AND "delivery_lat" IS NOT NULL
                AND "delivery_lng" IS NOT NULL
            )
        ),
    CONSTRAINT "chk_order_subtotal" CHECK ("subtotal" >= 0),
    CONSTRAINT "chk_order_tax" CHECK ("tax" >= 0),
    CONSTRAINT "chk_order_total" CHECK ("total" >= 0)
);

CREATE TABLE "order_item" (
    "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "tenant_id" UUID NOT NULL,
    "order_id" UUID NOT NULL,
    "product_id" UUID NOT NULL,
    "quantity" INTEGER NOT NULL,
    "unit_price" BIGINT NOT NULL,
    "vat_percent" NUMERIC(5,2) NOT NULL,
    "line_total" BIGINT NOT NULL,
    "note" VARCHAR(500),
    "prescription_ref" VARCHAR(255),
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "fk_order_item_tenant" FOREIGN KEY ("tenant_id") REFERENCES "tenant" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_order_item_order" FOREIGN KEY ("order_id") REFERENCES "order" ("id") ON DELETE CASCADE,
    CONSTRAINT "fk_order_item_product" FOREIGN KEY ("product_id") REFERENCES "product" ("id") ON DELETE RESTRICT,
    CONSTRAINT "chk_order_item_quantity" CHECK ("quantity" > 0),
    CONSTRAINT "chk_order_item_unit_price" CHECK ("unit_price" >= 0),
    CONSTRAINT "chk_order_item_line_total" CHECK ("line_total" >= 0)
);

CREATE INDEX "order_tenant_status_idx" ON "order" ("tenant_id", "status");
CREATE INDEX "order_item_order_idx" ON "order_item" ("order_id");

COMMIT;
