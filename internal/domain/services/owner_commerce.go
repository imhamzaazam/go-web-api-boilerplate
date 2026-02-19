package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/pgsqlc"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (service *ServiceManager) UpsertInventoryForProduct(ctx context.Context, inventory ports.NewInventory) (pgsqlc.Inventory, *domainerr.DomainError) {
	_, productErr := service.GetProductByTenantAndID(ctx, inventory.TenantID, inventory.ProductID)
	if productErr != nil {
		return pgsqlc.Inventory{}, domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "invalid tenant product", productErr)
	}

	updatedInventory, err := service.store.UpsertInventoryForProduct(ctx, pgsqlc.UpsertInventoryForProductParams{
		TenantID:  inventory.TenantID,
		ProductID: inventory.ProductID,
		InStock:   inventory.InStock,
	})
	if err != nil {
		return pgsqlc.Inventory{}, domainerr.MatchPostgresError(err)
	}

	return updatedInventory, nil
}

func (service *ServiceManager) CreateCartItem(ctx context.Context, item ports.NewCartItem) (pgsqlc.CreateCartItemRow, *domainerr.DomainError) {
	product, productErr := service.GetProductByTenantAndID(ctx, item.TenantID, item.ProductID)
	if productErr != nil {
		return pgsqlc.CreateCartItemRow{}, domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "invalid tenant product", productErr)
	}

	_, reserveErr := service.store.ReserveInventory(ctx, pgsqlc.ReserveInventoryParams{
		TenantID:  item.TenantID,
		ProductID: item.ProductID,
		Reserved:  item.Quantity,
	})
	if reserveErr != nil {
		if errors.Is(reserveErr, pgx.ErrNoRows) {
			return pgsqlc.CreateCartItemRow{}, domainerr.NewDomainError(http.StatusConflict, domainerr.QueryError, "insufficient stock", reserveErr)
		}
		return pgsqlc.CreateCartItemRow{}, domainerr.MatchPostgresError(reserveErr)
	}

	cart, cartErr := service.getOrCreateActiveCart(ctx, item.TenantID, item.UserUID)
	if cartErr != nil {
		return pgsqlc.CreateCartItemRow{}, cartErr
	}

	note := pgtype.Text{String: item.Note, Valid: item.HasNote}
	prescription := pgtype.Text{String: item.PrescriptionRef, Valid: item.HasPrescription}

	createdItem, err := service.store.CreateCartItem(ctx, pgsqlc.CreateCartItemParams{
		TenantID:        item.TenantID,
		CartID:          cart.ID,
		ProductID:       item.ProductID,
		Quantity:        item.Quantity,
		UnitPrice:       product.Price,
		VatPercent:      product.VatPercent,
		Note:            note,
		PrescriptionRef: prescription,
	})
	if err != nil {
		return pgsqlc.CreateCartItemRow{}, domainerr.MatchPostgresError(err)
	}

	return createdItem, nil
}

func (service *ServiceManager) getOrCreateActiveCart(ctx context.Context, tenantID uuid.UUID, userUID uuid.UUID) (pgsqlc.Cart, *domainerr.DomainError) {
	cart, err := service.store.GetActiveCartByTenantAndUser(ctx, pgsqlc.GetActiveCartByTenantAndUserParams{
		TenantID: tenantID,
		UserUid:  userUID,
	})
	if err == nil {
		return cart, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return pgsqlc.Cart{}, domainerr.MatchPostgresError(err)
	}

	createdCart, createErr := service.store.CreateCart(ctx, pgsqlc.CreateCartParams{
		TenantID: tenantID,
		UserUid:  userUID,
	})
	if createErr != nil {
		return pgsqlc.Cart{}, domainerr.MatchPostgresError(createErr)
	}

	return createdCart, nil
}
