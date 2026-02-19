package v1

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httperr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httputils"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/pgsqlc"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
)

type CreateOrderLocationDto struct {
	AddressLine string  `json:"address_line"`
	City        string  `json:"city"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
}

type CreateOrderRequestDto struct {
	CartID          string                  `json:"cart_id" validate:"required,uuid"`
	PaymentMethodID string                  `json:"payment_method_id" validate:"required"`
	FulfillmentType string                  `json:"fulfillment_type" validate:"required,oneof=pickup delivery"`
	Location        *CreateOrderLocationDto `json:"location"`
}

type OrderResponseDto struct {
	ID              uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	CartID          uuid.UUID `json:"cart_id"`
	Status          string    `json:"status"`
	FulfillmentType string    `json:"fulfillment_type"`
	Subtotal        int64     `json:"subtotal"`
	Tax             int64     `json:"tax"`
	Total           int64     `json:"total"`
}

func (adapter *HTTPAdapter) createOrder(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	tenantCtx, tenantErr := adapter.resolveTenantFromRequest(r)
	if tenantErr != nil {
		return tenantErr
	}

	reqDto, err := httputils.Decode[CreateOrderRequestDto](r)
	if err != nil {
		return err
	}

	validationErr := validate.Struct(reqDto)
	if validationErr != nil {
		return httperr.MatchValidationError(validationErr)
	}

	cartID, parseErr := uuid.Parse(reqDto.CartID)
	if parseErr != nil {
		return domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"cart_id": "The field must be a valid uuid"}, parseErr)
	}

	var location *ports.NewOrderLocation
	if reqDto.Location != nil {
		location = &ports.NewOrderLocation{
			AddressLine: strings.TrimSpace(reqDto.Location.AddressLine),
			City:        strings.TrimSpace(reqDto.Location.City),
			Lat:         reqDto.Location.Lat,
			Lng:         reqDto.Location.Lng,
		}
	}

	order, serviceErr := adapter.userService.CreateOrder(r.Context(), ports.NewOrder{
		TenantID:         tenantCtx.TenantID,
		CartID:           cartID,
		PaymentMethodRef: reqDto.PaymentMethodID,
		FulfillmentType:  reqDto.FulfillmentType,
		Location:         location,
	})
	if serviceErr != nil {
		return serviceErr
	}

	return httputils.Encode(w, r, http.StatusCreated, OrderResponseDto{
		ID:              order.ID,
		TenantID:        order.TenantID,
		CartID:          order.CartID,
		Status:          string(order.Status),
		FulfillmentType: string(order.FulfillmentType),
		Subtotal:        order.Subtotal,
		Tax:             order.Tax,
		Total:           order.Total,
	})
}

type UpdateOrderStatusRequestDto struct {
	Status string `json:"status" validate:"required,oneof=pending confirmed cancelled out_for_delivery completed refunded"`
}

func (adapter *HTTPAdapter) updateOrderStatus(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	tenantCtx, tenantErr := adapter.resolveTenantFromRequest(r)
	if tenantErr != nil {
		return tenantErr
	}

	orderID, parseErr := uuid.Parse(chi.URLParam(r, "id"))
	if parseErr != nil {
		return domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"id": "The field must be a valid uuid"}, parseErr)
	}

	reqDto, err := httputils.Decode[UpdateOrderStatusRequestDto](r)
	if err != nil {
		return err
	}

	validationErr := validate.Struct(reqDto)
	if validationErr != nil {
		return httperr.MatchValidationError(validationErr)
	}

	order, serviceErr := adapter.userService.UpdateOrderStatus(r.Context(), tenantCtx.TenantID, orderID, reqDto.Status)
	if serviceErr != nil {
		return serviceErr
	}

	return httputils.Encode(w, r, http.StatusOK, OrderResponseDto{
		ID:              order.ID,
		TenantID:        order.TenantID,
		CartID:          order.CartID,
		Status:          string(order.Status),
		FulfillmentType: string(order.FulfillmentType),
		Subtotal:        order.Subtotal,
		Tax:             order.Tax,
		Total:           order.Total,
	})
}

type CreatePaymentMethodRequestDto struct {
	Type      string `json:"type" validate:"required,oneof=cash card"`
	Label     string `json:"label" validate:"required"`
	IsDefault bool   `json:"is_default"`
}

type PaymentMethodResponseDto struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Label     string `json:"label"`
	IsDefault bool   `json:"is_default"`
}

func (adapter *HTTPAdapter) createPaymentMethod(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	tenantCtx, tenantErr := adapter.resolveTenantFromRequest(r)
	if tenantErr != nil {
		return tenantErr
	}

	reqDto, err := httputils.Decode[CreatePaymentMethodRequestDto](r)
	if err != nil {
		return err
	}

	validationErr := validate.Struct(reqDto)
	if validationErr != nil {
		return httperr.MatchValidationError(validationErr)
	}

	createdPaymentMethod, serviceErr := adapter.userService.CreatePaymentMethod(r.Context(), ports.NewPaymentMethod{
		TenantID:  tenantCtx.TenantID,
		Type:      reqDto.Type,
		Label:     reqDto.Label,
		IsDefault: reqDto.IsDefault,
	})
	if serviceErr != nil {
		return serviceErr
	}

	id := createdPaymentMethod.ID.String()
	if createdPaymentMethod.Type == pgsqlc.PaymentMethodTypeCash {
		id = "cash"
	}

	return httputils.Encode(w, r, http.StatusCreated, PaymentMethodResponseDto{
		ID:        id,
		Type:      string(createdPaymentMethod.Type),
		Label:     createdPaymentMethod.Label,
		IsDefault: createdPaymentMethod.IsDefault,
	})
}

func (adapter *HTTPAdapter) listPaymentMethods(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	tenantCtx, tenantErr := adapter.resolveTenantFromRequest(r)
	if tenantErr != nil {
		return tenantErr
	}

	paymentMethods, serviceErr := adapter.userService.ListPaymentMethods(r.Context(), tenantCtx.TenantID)
	if serviceErr != nil {
		return serviceErr
	}

	response := make([]PaymentMethodResponseDto, 0, len(paymentMethods))
	for _, paymentMethod := range paymentMethods {
		id := paymentMethod.ID.String()
		if paymentMethod.Type == pgsqlc.PaymentMethodTypeCash {
			id = "cash"
		}

		response = append(response, PaymentMethodResponseDto{
			ID:        id,
			Type:      string(paymentMethod.Type),
			Label:     paymentMethod.Label,
			IsDefault: paymentMethod.IsDefault,
		})
	}

	return httputils.Encode(w, r, http.StatusOK, response)
}

type PayOrderRequestDto struct {
	PaymentMethodID string `json:"payment_method_id" validate:"required"`
	Amount          int64  `json:"amount" validate:"required,min=1"`
}

type PayOrderResponseDto struct {
	OrderID uuid.UUID `json:"order_id"`
	Status  string    `json:"status"`
}

func (adapter *HTTPAdapter) payOrder(w http.ResponseWriter, r *http.Request) *domainerr.DomainError {
	tenantCtx, tenantErr := adapter.resolveTenantFromRequest(r)
	if tenantErr != nil {
		return tenantErr
	}

	orderID, parseErr := uuid.Parse(chi.URLParam(r, "id"))
	if parseErr != nil {
		return domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"id": "The field must be a valid uuid"}, parseErr)
	}

	reqDto, err := httputils.Decode[PayOrderRequestDto](r)
	if err != nil {
		return err
	}

	validationErr := validate.Struct(reqDto)
	if validationErr != nil {
		return httperr.MatchValidationError(validationErr)
	}

	paidOrder, serviceErr := adapter.userService.PayOrder(r.Context(), ports.PayOrder{
		TenantID:         tenantCtx.TenantID,
		OrderID:          orderID,
		PaymentMethodRef: reqDto.PaymentMethodID,
		Amount:           reqDto.Amount,
	})
	if serviceErr != nil {
		return serviceErr
	}

	return httputils.Encode(w, r, http.StatusOK, PayOrderResponseDto{
		OrderID: paidOrder.ID,
		Status:  string(paidOrder.Status),
	})
}
