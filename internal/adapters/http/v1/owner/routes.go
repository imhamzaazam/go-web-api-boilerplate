package owner

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) *domainerr.DomainError

type WrapFunc func(handlerFn HandlerFunc) http.HandlerFunc

type Handlers struct {
	CreateProduct             HandlerFunc
	CreateDiscount            HandlerFunc
	CreateProductAddon        HandlerFunc
	UpsertInventoryForProduct HandlerFunc
	CreateCartItem            HandlerFunc
	CreateOrder               HandlerFunc
	UpdateOrderStatus         HandlerFunc
	CreatePaymentMethod       HandlerFunc
	ListPaymentMethods        HandlerFunc
	PayOrder                  HandlerFunc
}

func RegisterRoutes(router chi.Router, wrap WrapFunc, handlers Handlers) {
	router.Post("/products", wrap(handlers.CreateProduct))
	router.Post("/products/{id}/discounts", wrap(handlers.CreateDiscount))
	router.Post("/products/{id}/addons", wrap(handlers.CreateProductAddon))
	router.Put("/inventory/{id}", wrap(handlers.UpsertInventoryForProduct))
	router.Post("/cart/items", wrap(handlers.CreateCartItem))
	router.Post("/orders", wrap(handlers.CreateOrder))
	router.Patch("/orders/{id}/status", wrap(handlers.UpdateOrderStatus))
	router.Post("/payment-methods", wrap(handlers.CreatePaymentMethod))
	router.Get("/payment-methods", wrap(handlers.ListPaymentMethods))
	router.Post("/orders/{id}/pay", wrap(handlers.PayOrder))
}
