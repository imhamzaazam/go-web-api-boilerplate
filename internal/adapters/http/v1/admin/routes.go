package admin

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) *domainerr.DomainError

type WrapFunc func(handlerFn HandlerFunc) http.HandlerFunc

type Handlers struct {
	CreateTenant       HandlerFunc
	CreateBranch       HandlerFunc
	CreateSubscription HandlerFunc
}

func RegisterRoutes(router chi.Router, wrap WrapFunc, handlers Handlers) {
	router.Post("/tenants", wrap(handlers.CreateTenant))
	router.Post("/branches", wrap(handlers.CreateBranch))
	router.Post("/subscriptions", wrap(handlers.CreateSubscription))
}
