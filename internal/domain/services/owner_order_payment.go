package services

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/pgsqlc"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (service *ServiceManager) CreateOrder(ctx context.Context, order ports.NewOrder) (pgsqlc.CreateOrderRow, *domainerr.DomainError) {
	fulfillmentType := pgsqlc.FulfillmentType(order.FulfillmentType)
	if fulfillmentType != pgsqlc.FulfillmentTypePickup && fulfillmentType != pgsqlc.FulfillmentTypeDelivery {
		return pgsqlc.CreateOrderRow{}, domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"fulfillment_type": "must be pickup or delivery"}, errors.New("invalid fulfillment type"))
	}

	tenant, tenantErr := service.store.GetTenantByID(ctx, order.TenantID)
	if tenantErr != nil {
		return pgsqlc.CreateOrderRow{}, domainerr.MatchPostgresError(tenantErr)
	}

	deliveryAddress := pgtype.Text{}
	deliveryCity := pgtype.Text{}
	deliveryLat := pgtype.Numeric{}
	deliveryLng := pgtype.Numeric{}
	if fulfillmentType == pgsqlc.FulfillmentTypeDelivery {
		if order.Location == nil || strings.TrimSpace(order.Location.AddressLine) == "" || strings.TrimSpace(order.Location.City) == "" {
			return pgsqlc.CreateOrderRow{}, domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"location": "delivery location is required"}, errors.New("missing delivery location"))
		}
		if order.Location.Lat < -90 || order.Location.Lat > 90 || order.Location.Lng < -180 || order.Location.Lng > 180 {
			return pgsqlc.CreateOrderRow{}, domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"location": "invalid coordinates"}, errors.New("invalid coordinates"))
		}
		if tenant.Type == pgsqlc.TenantTypeRestaurant {
			if strings.ToLower(strings.TrimSpace(order.Location.City)) != "karachi" || order.Location.Lat < 24.6 || order.Location.Lat > 25.2 || order.Location.Lng < 66.8 || order.Location.Lng > 67.5 {
				return pgsqlc.CreateOrderRow{}, domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"location": "outside service area"}, errors.New("outside service area"))
			}
		}

		latNumeric, latErr := numericFromFloat64(order.Location.Lat)
		if latErr != nil {
			return pgsqlc.CreateOrderRow{}, latErr
		}
		lngNumeric, lngErr := numericFromFloat64(order.Location.Lng)
		if lngErr != nil {
			return pgsqlc.CreateOrderRow{}, lngErr
		}

		deliveryAddress = pgtype.Text{String: strings.TrimSpace(order.Location.AddressLine), Valid: true}
		deliveryCity = pgtype.Text{String: strings.TrimSpace(order.Location.City), Valid: true}
		deliveryLat = latNumeric
		deliveryLng = lngNumeric
	}

	cart, cartErr := service.store.GetCartByTenantAndID(ctx, pgsqlc.GetCartByTenantAndIDParams{TenantID: order.TenantID, ID: order.CartID})
	if cartErr != nil {
		return pgsqlc.CreateOrderRow{}, domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "invalid tenant cart", cartErr)
	}

	itemCount, itemCountErr := service.store.CountCartItemsByTenantAndCart(ctx, pgsqlc.CountCartItemsByTenantAndCartParams{TenantID: order.TenantID, CartID: order.CartID})
	if itemCountErr != nil {
		return pgsqlc.CreateOrderRow{}, domainerr.MatchPostgresError(itemCountErr)
	}
	if itemCount == 0 {
		return pgsqlc.CreateOrderRow{}, domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"cart_id": "cart is empty"}, errors.New("empty cart"))
	}

	cartItems, cartItemsErr := service.store.ListCartItemsByTenantAndCart(ctx, pgsqlc.ListCartItemsByTenantAndCartParams{TenantID: order.TenantID, CartID: order.CartID})
	if cartItemsErr != nil {
		return pgsqlc.CreateOrderRow{}, domainerr.MatchPostgresError(cartItemsErr)
	}

	subtotal := int64(0)
	tax := int64(0)
	for _, cartItem := range cartItems {
		lineSubtotal := cartItem.UnitPrice * int64(cartItem.Quantity)
		subtotal += lineSubtotal

		vatPercent, vatErr := numericToFloat64InService(cartItem.VatPercent)
		if vatErr != nil {
			return pgsqlc.CreateOrderRow{}, vatErr
		}
		lineTax := int64(math.Round(float64(lineSubtotal) * vatPercent / 100))
		tax += lineTax
	}
	total := subtotal + tax

	createdOrder, createOrderErr := service.store.CreateOrder(ctx, pgsqlc.CreateOrderParams{
		TenantID:            order.TenantID,
		UserUid:             cart.UserUid,
		CartID:              order.CartID,
		Status:              pgsqlc.OrderStatusPending,
		FulfillmentType:     fulfillmentType,
		DeliveryAddressLine: deliveryAddress,
		DeliveryCity:        deliveryCity,
		DeliveryLat:         deliveryLat,
		DeliveryLng:         deliveryLng,
		Subtotal:            subtotal,
		Tax:                 tax,
		Total:               total,
	})
	if createOrderErr != nil {
		return pgsqlc.CreateOrderRow{}, domainerr.MatchPostgresError(createOrderErr)
	}

	for _, cartItem := range cartItems {
		lineTotal := cartItem.UnitPrice * int64(cartItem.Quantity)
		_, createOrderItemErr := service.store.CreateOrderItem(ctx, pgsqlc.CreateOrderItemParams{
			TenantID:        order.TenantID,
			OrderID:         createdOrder.ID,
			ProductID:       cartItem.ProductID,
			Quantity:        cartItem.Quantity,
			UnitPrice:       cartItem.UnitPrice,
			VatPercent:      cartItem.VatPercent,
			LineTotal:       lineTotal,
			Note:            cartItem.Note,
			PrescriptionRef: cartItem.PrescriptionRef,
		})
		if createOrderItemErr != nil {
			return pgsqlc.CreateOrderRow{}, domainerr.MatchPostgresError(createOrderItemErr)
		}
	}

	if setCartErr := service.store.SetCartInactive(ctx, pgsqlc.SetCartInactiveParams{TenantID: order.TenantID, ID: order.CartID}); setCartErr != nil {
		return pgsqlc.CreateOrderRow{}, domainerr.MatchPostgresError(setCartErr)
	}

	return createdOrder, nil
}

func (service *ServiceManager) UpdateOrderStatus(ctx context.Context, tenantID uuid.UUID, orderID uuid.UUID, status string) (pgsqlc.UpdateOrderStatusByTenantRow, *domainerr.DomainError) {
	targetStatus := pgsqlc.OrderStatus(status)
	if targetStatus != pgsqlc.OrderStatusPending && targetStatus != pgsqlc.OrderStatusConfirmed && targetStatus != pgsqlc.OrderStatusCancelled && targetStatus != pgsqlc.OrderStatusOutForDelivery && targetStatus != pgsqlc.OrderStatusCompleted && targetStatus != pgsqlc.OrderStatusRefunded {
		return pgsqlc.UpdateOrderStatusByTenantRow{}, domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"status": "invalid order status"}, errors.New("invalid order status"))
	}

	currentOrder, orderErr := service.store.GetOrderByTenantAndID(ctx, pgsqlc.GetOrderByTenantAndIDParams{TenantID: tenantID, ID: orderID})
	if orderErr != nil {
		return pgsqlc.UpdateOrderStatusByTenantRow{}, domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "invalid tenant order", orderErr)
	}

	if !isValidStatusTransition(currentOrder.Status, targetStatus) {
		return pgsqlc.UpdateOrderStatusByTenantRow{}, domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"status": "invalid status transition"}, errors.New("invalid status transition"))
	}

	if targetStatus == pgsqlc.OrderStatusRefunded && !currentOrder.IsPaid {
		return pgsqlc.UpdateOrderStatusByTenantRow{}, domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"status": "refunded allowed only for paid orders"}, errors.New("order is not paid"))
	}

	updatedOrder, updateErr := service.store.UpdateOrderStatusByTenant(ctx, pgsqlc.UpdateOrderStatusByTenantParams{
		TenantID: tenantID,
		ID:       orderID,
		Status:   targetStatus,
		Column4:  targetStatus == pgsqlc.OrderStatusCancelled,
		Column5:  targetStatus == pgsqlc.OrderStatusRefunded,
	})
	if updateErr != nil {
		return pgsqlc.UpdateOrderStatusByTenantRow{}, domainerr.MatchPostgresError(updateErr)
	}

	return updatedOrder, nil
}

func (service *ServiceManager) CreatePaymentMethod(ctx context.Context, paymentMethod ports.NewPaymentMethod) (pgsqlc.CreatePaymentMethodRow, *domainerr.DomainError) {
	if paymentMethod.Type != "cash" && paymentMethod.Type != "card" {
		return pgsqlc.CreatePaymentMethodRow{}, domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"type": "must be cash or card"}, errors.New("invalid payment method type"))
	}

	if err := service.store.EnsureDefaultCashPaymentMethod(ctx, paymentMethod.TenantID); err != nil {
		return pgsqlc.CreatePaymentMethodRow{}, domainerr.MatchPostgresError(err)
	}

	createdPaymentMethod, err := service.store.CreatePaymentMethod(ctx, pgsqlc.CreatePaymentMethodParams{
		TenantID:  paymentMethod.TenantID,
		Type:      pgsqlc.PaymentMethodType(paymentMethod.Type),
		Label:     paymentMethod.Label,
		IsDefault: paymentMethod.IsDefault,
	})
	if err != nil {
		return pgsqlc.CreatePaymentMethodRow{}, domainerr.MatchPostgresError(err)
	}

	return createdPaymentMethod, nil
}

func (service *ServiceManager) ListPaymentMethods(ctx context.Context, tenantID uuid.UUID) ([]pgsqlc.ListPaymentMethodsByTenantRow, *domainerr.DomainError) {
	if err := service.store.EnsureDefaultCashPaymentMethod(ctx, tenantID); err != nil {
		return nil, domainerr.MatchPostgresError(err)
	}

	paymentMethods, err := service.store.ListPaymentMethodsByTenant(ctx, tenantID)
	if err != nil {
		return nil, domainerr.MatchPostgresError(err)
	}

	return paymentMethods, nil
}

func (service *ServiceManager) PayOrder(ctx context.Context, payOrder ports.PayOrder) (pgsqlc.MarkOrderPaidByTenantRow, *domainerr.DomainError) {
	if payOrder.Amount <= 0 {
		return pgsqlc.MarkOrderPaidByTenantRow{}, domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"amount": "must be greater than 0"}, errors.New("invalid amount"))
	}

	order, orderErr := service.store.GetOrderByTenantAndID(ctx, pgsqlc.GetOrderByTenantAndIDParams{TenantID: payOrder.TenantID, ID: payOrder.OrderID})
	if orderErr != nil {
		return pgsqlc.MarkOrderPaidByTenantRow{}, domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "invalid tenant order", orderErr)
	}

	if order.IsPaid {
		return pgsqlc.MarkOrderPaidByTenantRow{}, domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"order": "order already paid"}, errors.New("order already paid"))
	}

	if payOrder.Amount != order.Total {
		return pgsqlc.MarkOrderPaidByTenantRow{}, domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"amount": "must match order total"}, errors.New("amount mismatch"))
	}

	paymentMethodID, paymentMethodErr := service.resolvePaymentMethodID(ctx, payOrder.TenantID, payOrder.PaymentMethodRef)
	if paymentMethodErr != nil {
		return pgsqlc.MarkOrderPaidByTenantRow{}, paymentMethodErr
	}

	paidOrder, markPaidErr := service.store.MarkOrderPaidByTenant(ctx, pgsqlc.MarkOrderPaidByTenantParams{
		TenantID:        payOrder.TenantID,
		ID:              payOrder.OrderID,
		PaymentMethodID: paymentMethodID,
	})
	if markPaidErr != nil {
		if errors.Is(markPaidErr, pgx.ErrNoRows) {
			return pgsqlc.MarkOrderPaidByTenantRow{}, domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"order": "order already paid"}, markPaidErr)
		}
		return pgsqlc.MarkOrderPaidByTenantRow{}, domainerr.MatchPostgresError(markPaidErr)
	}

	return paidOrder, nil
}

func (service *ServiceManager) resolvePaymentMethodID(ctx context.Context, tenantID uuid.UUID, paymentMethodRef string) (uuid.UUID, *domainerr.DomainError) {
	if paymentMethodRef == "cash" {
		if err := service.store.EnsureDefaultCashPaymentMethod(ctx, tenantID); err != nil {
			return uuid.Nil, domainerr.MatchPostgresError(err)
		}
		cashMethod, err := service.store.GetDefaultCashPaymentMethodByTenant(ctx, tenantID)
		if err != nil {
			return uuid.Nil, domainerr.MatchPostgresError(err)
		}
		return cashMethod.ID, nil
	}

	parsedID, parseErr := uuid.Parse(paymentMethodRef)
	if parseErr != nil {
		return uuid.Nil, domainerr.NewDomainError(http.StatusUnprocessableEntity, domainerr.ValidationError, map[string]string{"payment_method_id": "must be cash or a valid uuid"}, parseErr)
	}

	paymentMethod, paymentMethodErr := service.store.GetPaymentMethodByTenantAndID(ctx, pgsqlc.GetPaymentMethodByTenantAndIDParams{TenantID: tenantID, ID: parsedID})
	if paymentMethodErr != nil {
		return uuid.Nil, domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "invalid tenant payment method", paymentMethodErr)
	}

	return paymentMethod.ID, nil
}

func numericToFloat64InService(value pgtype.Numeric) (float64, *domainerr.DomainError) {
	floatVal, err := value.Float64Value()
	if err != nil {
		return 0, domainerr.NewDomainError(http.StatusInternalServerError, domainerr.InternalError, "invalid numeric value", err)
	}
	if !floatVal.Valid {
		return 0, domainerr.NewDomainError(http.StatusInternalServerError, domainerr.InternalError, "invalid numeric value", errors.New("numeric value is not valid"))
	}
	return floatVal.Float64, nil
}

func isValidStatusTransition(current pgsqlc.OrderStatus, target pgsqlc.OrderStatus) bool {
	if current == target {
		return false
	}

	switch current {
	case pgsqlc.OrderStatusPending:
		return target == pgsqlc.OrderStatusConfirmed || target == pgsqlc.OrderStatusCancelled
	case pgsqlc.OrderStatusConfirmed:
		return target == pgsqlc.OrderStatusOutForDelivery || target == pgsqlc.OrderStatusCancelled || target == pgsqlc.OrderStatusRefunded
	case pgsqlc.OrderStatusOutForDelivery:
		return target == pgsqlc.OrderStatusCompleted || target == pgsqlc.OrderStatusRefunded
	case pgsqlc.OrderStatusCompleted:
		return target == pgsqlc.OrderStatusRefunded
	default:
		return false
	}
}
