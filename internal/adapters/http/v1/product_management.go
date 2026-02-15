package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httperr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httputils"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateProductRequestDto struct {
	Name                 string  `json:"name" validate:"required"`
	SKU                  string  `json:"sku" validate:"required"`
	Price                int64   `json:"price" validate:"required,min=0"`
	VATPercent           float64 `json:"vat_percent" validate:"required,min=0,max=100"`
	IsPreorder           bool    `json:"is_preorder"`
	MadeToOrder          bool    `json:"made_to_order"`
	RequiresPrescription bool    `json:"requires_prescription"`
	AvailableForDelivery bool    `json:"available_for_delivery"`
	AvailableForPickup   bool    `json:"available_for_pickup"`
}

type CreateProductResponseDto struct {
	ID       uuid.UUID `json:"id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Name     string    `json:"name"`
	SKU      string    `json:"sku"`
	Price    int64     `json:"price"`
}

func (adapter *HTTPAdapter) createProduct(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	tenantCtx, tenantErr := adapter.resolveTenantFromRequest(r)
	if tenantErr != nil {
		return tenantErr
	}

	reqDto, err := httputils.Decode[CreateProductRequestDto](r)
	if err != nil {
		return err
	}

	validationErr := validate.Struct(reqDto)
	if validationErr != nil {
		return httperr.MatchValidationError(validationErr)
	}

	if reqDto.MadeToOrder && tenantCtx.TenantType != "bakery" {
		return domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"made_to_order": "only allowed for bakery"}, errors.New("made_to_order only allowed for bakery"))
	}

	product, serviceErr := adapter.userService.CreateProduct(r.Context(), ports.NewProduct{
		TenantID:             tenantCtx.TenantID,
		Name:                 reqDto.Name,
		Sku:                  reqDto.SKU,
		Price:                reqDto.Price,
		VATPercent:           reqDto.VATPercent,
		IsPreorder:           reqDto.IsPreorder,
		MadeToOrder:          reqDto.MadeToOrder,
		RequiresPrescription: reqDto.RequiresPrescription,
		AvailableForDelivery: reqDto.AvailableForDelivery,
		AvailableForPickup:   reqDto.AvailableForPickup,
	})
	if serviceErr != nil {
		return serviceErr
	}

	return httputils.Encode(w, r, http.StatusCreated, CreateProductResponseDto{
		ID:       product.ID,
		TenantID: product.TenantID,
		Name:     product.Name,
		SKU:      product.Sku,
		Price:    product.Price,
	})
}

type CreateDiscountRequestDto struct {
	Code     string    `json:"code" validate:"required"`
	Name     string    `json:"name" validate:"required"`
	Type     string    `json:"type" validate:"required,oneof=percentage"`
	Value    float64   `json:"value" validate:"required,min=1,max=100"`
	StartsAt time.Time `json:"starts_at" validate:"required"`
	EndsAt   time.Time `json:"ends_at" validate:"required"`
}

type CreateDiscountResponseDto struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
}

func (adapter *HTTPAdapter) createDiscount(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	tenantCtx, tenantErr := adapter.resolveTenantFromRequest(r)
	if tenantErr != nil {
		return tenantErr
	}

	productID, parseErr := uuid.Parse(chi.URLParam(r, "id"))
	if parseErr != nil {
		return domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"id": "The field must be a valid uuid"}, parseErr)
	}

	reqDto, err := httputils.Decode[CreateDiscountRequestDto](r)
	if err != nil {
		return err
	}

	validationErr := validate.Struct(reqDto)
	if validationErr != nil {
		return httperr.MatchValidationError(validationErr)
	}

	if !reqDto.StartsAt.Before(reqDto.EndsAt) {
		return domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"ends_at": "must be after starts_at"}, errors.New("invalid discount dates"))
	}

	_, productErr := adapter.userService.GetProductByTenantAndID(r.Context(), tenantCtx.TenantID, productID)
	if productErr != nil {
		return domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "invalid tenant product", productErr)
	}

	discount, serviceErr := adapter.userService.CreateDiscount(r.Context(), ports.NewDiscount{
		TenantID:  tenantCtx.TenantID,
		ProductID: productID,
		Code:      reqDto.Code,
		Name:      reqDto.Name,
		Type:      reqDto.Type,
		Value:     reqDto.Value,
		StartsAt:  reqDto.StartsAt,
		EndsAt:    reqDto.EndsAt,
	})
	if serviceErr != nil {
		return serviceErr
	}

	value, valueErr := numericToFloat64(discount.Value)
	if valueErr != nil {
		return valueErr
	}

	return httputils.Encode(w, r, http.StatusCreated, CreateDiscountResponseDto{
		ID:        discount.ID,
		ProductID: discount.ProductID,
		Code:      discount.Code,
		Name:      discount.Name,
		Type:      string(discount.Type),
		Value:     value,
		StartsAt:  discount.StartsAt,
		EndsAt:    discount.EndsAt,
	})
}

type CreateProductAddonRequestDto struct {
	Name  string `json:"name" validate:"required"`
	Price int64  `json:"price" validate:"required,min=0"`
}

type CreateProductAddonResponseDto struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
}

func (adapter *HTTPAdapter) createProductAddon(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	tenantCtx, tenantErr := adapter.resolveTenantFromRequest(r)
	if tenantErr != nil {
		return tenantErr
	}

	if tenantCtx.TenantType == "pharmacy" {
		return domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"addons": "addons are not allowed for pharmacy"}, errors.New("addons are not allowed for pharmacy"))
	}

	productID, parseErr := uuid.Parse(chi.URLParam(r, "id"))
	if parseErr != nil {
		return domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"id": "The field must be a valid uuid"}, parseErr)
	}

	reqDto, err := httputils.Decode[CreateProductAddonRequestDto](r)
	if err != nil {
		return err
	}

	validationErr := validate.Struct(reqDto)
	if validationErr != nil {
		return httperr.MatchValidationError(validationErr)
	}

	_, productErr := adapter.userService.GetProductByTenantAndID(r.Context(), tenantCtx.TenantID, productID)
	if productErr != nil {
		return domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "invalid tenant product", productErr)
	}

	addon, serviceErr := adapter.userService.CreateProductAddon(r.Context(), ports.NewProductAddon{
		TenantID:  tenantCtx.TenantID,
		ProductID: productID,
		Name:      reqDto.Name,
		Price:     reqDto.Price,
	})
	if serviceErr != nil {
		return serviceErr
	}

	return httputils.Encode(w, r, http.StatusCreated, CreateProductAddonResponseDto{
		ID:        addon.ID,
		ProductID: addon.ProductID,
		Name:      addon.Name,
		Price:     addon.Price,
	})
}

func numericToFloat64(value pgtype.Numeric) (float64, *domainerr.DomainError) {
	floatVal, err := value.Float64Value()
	if err != nil {
		return 0, domainerr.NewDomainError(http.StatusInternalServerError, domainerr.InternalError, "invalid numeric value", err)
	}
	if !floatVal.Valid {
		return 0, domainerr.NewDomainError(http.StatusInternalServerError, domainerr.InternalError, "invalid numeric value", errors.New("numeric value is not valid"))
	}

	return floatVal.Float64, nil
}
