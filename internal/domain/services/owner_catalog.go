package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/pgsqlc"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
)

func (service *ServiceManager) CreateProduct(ctx context.Context, product ports.NewProduct) (pgsqlc.CreateProductRow, *domainerr.DomainError) {
	tenant, tenantErr := service.store.GetTenantByID(ctx, product.TenantID)
	if tenantErr != nil {
		return pgsqlc.CreateProductRow{}, domainerr.MatchPostgresError(tenantErr)
	}

	if product.MadeToOrder && tenant.Type != pgsqlc.TenantTypeBakery {
		return pgsqlc.CreateProductRow{}, domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"made_to_order": "only allowed for bakery"}, errors.New("made_to_order only allowed for bakery"))
	}

	vatPercent, vatErr := numericFromFloat64(product.VATPercent)
	if vatErr != nil {
		return pgsqlc.CreateProductRow{}, vatErr
	}

	createdProduct, err := service.store.CreateProduct(ctx, pgsqlc.CreateProductParams{
		TenantID:             product.TenantID,
		Name:                 product.Name,
		Sku:                  product.Sku,
		Price:                product.Price,
		VatPercent:           vatPercent,
		IsPreorder:           product.IsPreorder,
		MadeToOrder:          product.MadeToOrder,
		RequiresPrescription: product.RequiresPrescription,
		AvailableForDelivery: product.AvailableForDelivery,
		AvailableForPickup:   product.AvailableForPickup,
	})
	if err != nil {
		return pgsqlc.CreateProductRow{}, domainerr.MatchPostgresError(err)
	}

	return createdProduct, nil
}

func (service *ServiceManager) GetProductByTenantAndID(ctx context.Context, tenantID uuid.UUID, productID uuid.UUID) (pgsqlc.GetProductByTenantAndIDRow, *domainerr.DomainError) {
	product, err := service.store.GetProductByTenantAndID(ctx, pgsqlc.GetProductByTenantAndIDParams{
		TenantID: tenantID,
		ID:       productID,
	})
	if err != nil {
		return pgsqlc.GetProductByTenantAndIDRow{}, domainerr.MatchPostgresError(err)
	}

	return product, nil
}

func (service *ServiceManager) CreateDiscount(ctx context.Context, discount ports.NewDiscount) (pgsqlc.CreateDiscountRow, *domainerr.DomainError) {
	discountValue, valueErr := numericFromFloat64(discount.Value)
	if valueErr != nil {
		return pgsqlc.CreateDiscountRow{}, valueErr
	}

	createdDiscount, err := service.store.CreateDiscount(ctx, pgsqlc.CreateDiscountParams{
		TenantID:  discount.TenantID,
		ProductID: discount.ProductID,
		Code:      discount.Code,
		Name:      discount.Name,
		Type:      pgsqlc.DiscountType(discount.Type),
		Value:     discountValue,
		StartsAt:  discount.StartsAt,
		EndsAt:    discount.EndsAt,
	})
	if err != nil {
		return pgsqlc.CreateDiscountRow{}, domainerr.MatchPostgresError(err)
	}

	return createdDiscount, nil
}

func (service *ServiceManager) CreateProductAddon(ctx context.Context, addon ports.NewProductAddon) (pgsqlc.CreateProductAddonRow, *domainerr.DomainError) {
	createdAddon, err := service.store.CreateProductAddon(ctx, pgsqlc.CreateProductAddonParams{
		TenantID:  addon.TenantID,
		ProductID: addon.ProductID,
		Name:      addon.Name,
		Price:     addon.Price,
	})
	if err != nil {
		return pgsqlc.CreateProductAddonRow{}, domainerr.MatchPostgresError(err)
	}

	return createdAddon, nil
}
