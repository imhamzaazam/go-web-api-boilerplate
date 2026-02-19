package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httputils"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/middleware"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/token"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
)

func (adapter *HTTPAdapter) resolveTenantFromRequest(r *http.Request) (ports.TenantContext, *domainerr.DomainError) {
	host := strings.TrimSpace(r.Host)
	host = strings.Split(host, ":")[0]

	if host == "" {
		host = "default.localhost"
	}

	return adapter.userService.ResolveTenantContext(r.Context(), host)
}

func (adapter *HTTPAdapter) tenantClaimsGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantCtx, tenantErr := adapter.resolveTenantFromRequest(r)
		if tenantErr != nil {
			_ = httputils.Encode(w, r, tenantErr.HTTPCode, tenantErr.HTTPErrorBody)
			return
		}

		payload, ok := r.Context().Value(middleware.KeyAuthUser).(*token.Payload)
		if !ok || payload == nil {
			httpErr := domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "missing auth claims", fmt.Errorf("missing auth claims"))
			_ = httputils.Encode(w, r, httpErr.HTTPCode, httpErr.HTTPErrorBody)
			return
		}

		if payload.TenantID != tenantCtx.TenantID || payload.TenantSlug != tenantCtx.TenantSlug {
			httpErr := domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "tenant host and token claims mismatch", fmt.Errorf("tenant host and token claims mismatch"))
			_ = httputils.Encode(w, r, httpErr.HTTPCode, httpErr.HTTPErrorBody)
			return
		}

		next.ServeHTTP(w, r)
	})
}
