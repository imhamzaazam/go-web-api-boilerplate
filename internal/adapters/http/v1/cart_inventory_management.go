package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httperr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httputils"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
)

type UpsertInventoryRequestDto struct {
	InStock int32 `json:"in_stock" validate:"required,min=0"`
}

type UpsertInventoryResponseDto struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	ProductID uuid.UUID `json:"product_id"`
	InStock   int32     `json:"in_stock"`
	Reserved  int32     `json:"reserved"`
}

func (adapter *HTTPAdapter) upsertInventoryForProduct(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	tenantCtx, tenantErr := adapter.resolveTenantFromRequest(r)
	if tenantErr != nil {
		return tenantErr
	}

	productID, parseErr := uuid.Parse(chi.URLParam(r, "id"))
	if parseErr != nil {
		return domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"id": "The field must be a valid uuid"}, parseErr)
	}

	reqDto, err := httputils.Decode[UpsertInventoryRequestDto](r)
	if err != nil {
		return err
	}

	validationErr := validate.Struct(reqDto)
	if validationErr != nil {
		return httperr.MatchValidationError(validationErr)
	}

	inventory, serviceErr := adapter.userService.UpsertInventoryForProduct(r.Context(), ports.NewInventory{
		TenantID:  tenantCtx.TenantID,
		ProductID: productID,
		InStock:   reqDto.InStock,
	})
	if serviceErr != nil {
		return serviceErr
	}

	return httputils.Encode(w, r, http.StatusOK, UpsertInventoryResponseDto{
		ID:        inventory.ID,
		TenantID:  inventory.TenantID,
		ProductID: inventory.ProductID,
		InStock:   inventory.InStock,
		Reserved:  inventory.Reserved,
	})
}

type CreateCartItemRequestDto struct {
	UserUID         string  `json:"user_uid" validate:"required,uuid"`
	ProductID       string  `json:"product_id" validate:"required,uuid"`
	Quantity        int32   `json:"quantity" validate:"required,min=1"`
	Note            *string `json:"note"`
	PrescriptionRef *string `json:"prescription_ref"`
}

type CreateCartItemResponseDto struct {
	ID        uuid.UUID `json:"id"`
	CartID    uuid.UUID `json:"cart_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int32     `json:"quantity"`
	UnitPrice int64     `json:"unit_price"`
}

func (adapter *HTTPAdapter) createCartItem(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	tenantCtx, tenantErr := adapter.resolveTenantFromRequest(r)
	if tenantErr != nil {
		return tenantErr
	}

	reqDto, err := httputils.Decode[CreateCartItemRequestDto](r)
	if err != nil {
		return err
	}

	validationErr := validate.Struct(reqDto)
	if validationErr != nil {
		return httperr.MatchValidationError(validationErr)
	}

	userUID, userParseErr := uuid.Parse(reqDto.UserUID)
	if userParseErr != nil {
		return domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"user_uid": "The field must be a valid uuid"}, userParseErr)
	}

	productID, productParseErr := uuid.Parse(reqDto.ProductID)
	if productParseErr != nil {
		return domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"product_id": "The field must be a valid uuid"}, productParseErr)
	}

	note := ""
	hasNote := false
	if reqDto.Note != nil {
		note = *reqDto.Note
		hasNote = true
	}

	prescription := ""
	hasPrescription := false
	if reqDto.PrescriptionRef != nil {
		prescription = *reqDto.PrescriptionRef
		hasPrescription = true
	}

	cartItem, serviceErr := adapter.userService.CreateCartItem(r.Context(), ports.NewCartItem{
		TenantID:        tenantCtx.TenantID,
		UserUID:         userUID,
		ProductID:       productID,
		Quantity:        reqDto.Quantity,
		Note:            note,
		PrescriptionRef: prescription,
		HasNote:         hasNote,
		HasPrescription: hasPrescription,
	})
	if serviceErr != nil {
		return serviceErr
	}

	return httputils.Encode(w, r, http.StatusOK, CreateCartItemResponseDto{
		ID:        cartItem.ID,
		CartID:    cartItem.CartID,
		ProductID: cartItem.ProductID,
		Quantity:  cartItem.Quantity,
		UnitPrice: cartItem.UnitPrice,
	})
}
