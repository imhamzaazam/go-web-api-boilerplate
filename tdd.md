# TDD Process Instructions

Follow strict TDD process:
1) Write failing test(s) first.
2) Implement minimum code to pass.
3) Refactor with tests green.
4) Keep scope minimal to this document.
5) Work in the order defined below.

---

## Project Context

- Architecture: Hexagonal
- API base: `/api/v1`
- Tenancy: shared DB, strict tenant isolation
- Auth: JWT bearer token
- Cross-tenant policy: `401 Unauthorized`
- Subscription gate policy: `402 Payment Required`

---

## Core Rules

- All business tables include `tenant_id`.
- Tenant is resolved from request host/subdomain.
- Tenant domain/subdomain is globally unique.
- Same user email may exist in multiple tenants.
- Never trust tenant from request body if token exists.
- Soft-delete means `deleted_at IS NOT NULL`.
- Soft-delete scope (v1): `tenant`, `subscription`, `product`.

---

## Domain Targets (v1)

- Tenant (`bakery`, `pharmacy`, `restaurant`)
- Branch
- Subscription
- Product
- Discount
- Product Add-on
- Cart
- Cart Item
- Inventory
- Order
- Order Item
- Payment Method

---

## Enums

### SubscriptionStatus
- `trial`
- `active`
- `past_due`
- `suspended`
- `canceled`

### OrderStatus
- `pending`
- `confirmed`
- `cancelled`
- `out_for_delivery`
- `completed`
- `refunded`

### Roles
- `owner`
- `admin`
- `employee`

### DiscountType (v1)
- `percentage`

---

## Permissions

- Owner:
  - manage employees
  - manage branches
  - manage subscription lifecycle
  - view sales/revenue reports
  - full business access

- Admin:
  - full business access
  - view sales/revenue reports
  - cannot cancel/delete subscription

- Employee:
  - product/cart/order/payment operations only
  - cannot manage employees/branches/subscription/reports

---

## Access Policies

### Cross-Tenant
- Accessing other tenant data returns `401`.

### Subscription Gate
- Allowed statuses for business operations: `trial`, `active`.
- Blocked statuses: `past_due`, `suspended`, `canceled`.
- Blocked response:
```json
{ "code": "subscription/inactive", "errors": "Subscription is not active" }
```

### Tenant Resolution
- Tenant is resolved by host/subdomain (example: `bakery-a.localhost`).
- Unknown host/subdomain returns `401`.
- Login authenticates in resolved tenant context only.

---

## Pricing Rules

- `price` in minor units (integer).
- `vat_percent` is percentage (integer or decimal).
- Discount type is percentage only in v1.
- Discount `value` range: `1..100`.
- Discount applies to base price only.
- VAT is computed after discount on discounted base price.

---

## Order Transition Rules

Allowed transitions:
- `pending -> confirmed`
- `confirmed -> out_for_delivery`
- `out_for_delivery -> completed`
- `pending -> cancelled`
- `confirmed -> cancelled`
- `confirmed|out_for_delivery|completed -> refunded` (only if paid)

Any other transition returns `422`.

---

## Vertical Rules

### Bakery
- add-ons allowed
- preorder/made-to-order allowed
- stock required unless made-to-order

### Pharmacy
- product has `requires_prescription`
- add-ons not allowed
- `prescription_ref` in cart item is optional and defaults `null`
- no oversell

### Restaurant
- add-ons allowed
- supports delivery/pickup flags
- cart item notes allowed

Restaurant flow coverage (inspired by live ordering UX):
- user must select order type first: `delivery` or `pickup`
- for delivery, location is mandatory
- support current-location flow (`lat`, `lng`) and manual address fallback
- changing location may change available branch/menu
- card payment is available in checkout alongside default cash
- primary city example for v1 tests: `Karachi`
- restaurant sample catalog style: Indo-Chinese mains, rice/noodles, sides, beverages

---

## API Contracts

Error format:
```json
{ "code": "string", "errors": "string|object|array" }
```

### POST /tenants
Request:
```json
{ "name": "Continental Bakery", "slug": "continental-bakery", "type": "bakery", "domain": "bakery-a.localhost" }
```
Success `201`:
```json
{ "id": "uuid", "name": "Continental Bakery", "slug": "continental-bakery", "type": "bakery", "domain": "bakery-a.localhost", "created_at": "2026-02-14T10:00:00Z" }
```
Failure: `422`, `409` duplicate slug/domain

### POST /branches
Request:
```json
{ "tenant_id": "uuid", "name": "Main Branch", "code": "BR001" }
```
Success `201`
Failure: `401`, `409`, `422`

### POST /subscriptions
Request:
```json
{ "tenant_id": "uuid", "plan": "starter", "status": "trial", "starts_at": "2026-02-14T00:00:00Z", "ends_at": "2026-03-14T00:00:00Z" }
```
Success `201`
Failure: `404` tenant not found, `422` invalid status, `409` active subscription conflict

### POST /login
Request:
```json
{ "email": "john@x.com", "password": "secret" }
```
Success `200` with token claims including:
- `tenant_id`
- `tenant_slug`
- `subscription_status`
Failure: `401` invalid credentials, `401` unknown host tenant

### POST /products
Request:
```json
{ "name": "Chocolate Cake", "sku": "CAKE-001", "price": 1299, "vat_percent": 15 }
```
Success `201`
Failure: `409` duplicate sku in tenant, `402` inactive subscription, `401` cross-tenant

### POST /products/{id}/discounts
Request:
```json
{ "code": "PAK001", "name": "Ramadan Offer", "type": "percentage", "value": 10, "starts_at": "2026-02-14T00:00:00Z", "ends_at": "2026-02-28T23:59:59Z" }
```
Success `201`
Failure: `422` invalid range/type, `409` duplicate code, `401` cross-tenant

### POST /products/{id}/addons
Request:
```json
{ "name": "Extra Cream", "price": 199 }
```
Success `201`
Failure: `422` pharmacy tenant, `401` cross-tenant

### POST /cart/items
Request:
```json
{ "product_id": "uuid", "quantity": 2, "addon_ids": ["uuid"], "prescription_ref": null, "note": "Less sugar" }
```
Success `200`
Failure: `409` out of stock, `402` inactive subscription, `401` cross-tenant

### POST /orders
Request:
```json
{
  "cart_id": "uuid",
  "payment_method_id": "uuid",
  "fulfillment_type": "delivery",
  "location": {
    "address_line": "221B Main Street",
    "city": "Karachi",
    "lat": 24.8607,
    "lng": 67.0011
  }
}
```
Success `201` (status `pending`)
Failure: `422` empty cart, `422` missing location when `fulfillment_type=delivery`, `422` invalid order type, `401` cross-tenant, `402` inactive subscription

Additional restaurant checks:
- `pickup` order must not require location
- `delivery` order requires valid `lat/lng`
- location outside service area returns `422`

### POST /orders/{id}/pay
Request:
```json
{ "payment_method_id": "uuid", "amount": 2797 }
```
Success `200` (status `paid`)
Failure: `422` invalid transition, `401` wrong tenant payment method, `402` inactive subscription

### POST /payment-methods
Request:
```json
{ "type": "card", "label": "Visa **** 4242", "is_default": false }
```
Success `201`
Failure: `422`, `401`

### GET /payment-methods
Success `200` includes default cash:
```json
[{ "id": "cash", "type": "cash", "label": "Cash", "is_default": true }]
```

---

## Test Plan (Strict Order)

1. Tenant + branch + subscription creation linkage
2. Login token claims with host-based tenant resolution
3. Subscription gate allow/block behavior
4. Cross-tenant `401` behavior
5. Product CRUD (tenant scope)
6. Discount rules
7. Add-on rules
8. Cart/cart-item rules
9. Inventory reserve/deduct/insufficient stock
10. Order lifecycle + transition validation
11. Payment methods + pay order
12. Soft-delete visibility/blocking

---

## Failure Test Matrix (Must Include)

For every endpoint, include failure tests for:
- invalid payload (`422`)
- unauthorized (`401`)
- cross-tenant (`401`)
- subscription-inactive (`402`) where business endpoint
- duplicate conflict (`409`) where unique constraints exist
- missing resource (`404`) where applicable

---

## Test Definition Template

### Requirement: <name>
- Tenant Type: bakery|pharmacy|restaurant|all
- Given:
- When:
- Then:
- Expected HTTP code:
- Expected DB state:

### Minimum scenarios
- happy path
- validation failure
- authorization failure
- cross-tenant failure
- subscription gate failure
- soft-delete behavior

---

## Folder Rules

- Migrations: `db/postgres/migration`
- Queries: `db/postgres/query`
- SQLC generated: `internal/adapters/pgsqlc`
- Integration tests: `internal/adapters/http/v1`

---

## Done Criteria

- All scoped tests green.
- No cross-tenant leak.
- Subscription gate enforced.
- Soft-delete enforced for tenant/subscription/product.
- OpenAPI updated and `make openapi-check` passes.
