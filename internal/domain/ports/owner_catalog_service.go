package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/pgsqlc"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
)

type OwnerCatalogService interface {
	CreateProduct(ctx context.Context, product NewProduct) (pgsqlc.CreateProductRow, *domainerr.DomainError)
	GetProductByTenantAndID(ctx context.Context, tenantID uuid.UUID, productID uuid.UUID) (pgsqlc.GetProductByTenantAndIDRow, *domainerr.DomainError)
	CreateDiscount(ctx context.Context, discount NewDiscount) (pgsqlc.CreateDiscountRow, *domainerr.DomainError)
	CreateProductAddon(ctx context.Context, addon NewProductAddon) (pgsqlc.CreateProductAddonRow, *domainerr.DomainError)
}
