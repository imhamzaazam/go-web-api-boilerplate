package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/pgsqlc"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
)

type OwnerCommerceService interface {
	UpsertInventoryForProduct(ctx context.Context, inventory NewInventory) (pgsqlc.Inventory, *domainerr.DomainError)
	CreateCartItem(ctx context.Context, item NewCartItem) (pgsqlc.CreateCartItemRow, *domainerr.DomainError)
	CreateOrder(ctx context.Context, order NewOrder) (pgsqlc.CreateOrderRow, *domainerr.DomainError)
	UpdateOrderStatus(ctx context.Context, tenantID uuid.UUID, orderID uuid.UUID, status string) (pgsqlc.UpdateOrderStatusByTenantRow, *domainerr.DomainError)
	CreatePaymentMethod(ctx context.Context, paymentMethod NewPaymentMethod) (pgsqlc.CreatePaymentMethodRow, *domainerr.DomainError)
	ListPaymentMethods(ctx context.Context, tenantID uuid.UUID) ([]pgsqlc.ListPaymentMethodsByTenantRow, *domainerr.DomainError)
	PayOrder(ctx context.Context, payOrder PayOrder) (pgsqlc.MarkOrderPaidByTenantRow, *domainerr.DomainError)
}
