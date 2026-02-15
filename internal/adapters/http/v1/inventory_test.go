package v1

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/horiondreher/go-web-api-boilerplate/internal/utils"
	"github.com/stretchr/testify/require"
)

func Test_Inventory_Upsert_TenantProtection(t *testing.T) {
	server, err := NewHTTPAdapter(testUserService)
	require.NoError(t, err)

	bakerySlug := "bakery-" + utils.RandomString(6)
	bakeryDomain := bakerySlug + ".localhost"
	bakeryID := createTenantForProductFlow(t, server, bakerySlug, bakeryDomain, "bakery")
	createSubscriptionForProductFlow(t, server, bakeryID, bakeryDomain, "active")
	bakeryAccessToken := createAccessTokenForTenantTest(t, server, bakeryDomain)

	restaurantSlug := "restaurant-" + utils.RandomString(6)
	restaurantDomain := restaurantSlug + ".localhost"
	restaurantID := createTenantForProductFlow(t, server, restaurantSlug, restaurantDomain, "restaurant")
	createSubscriptionForProductFlow(t, server, restaurantID, restaurantDomain, "active")
	restaurantAccessToken := createAccessTokenForTenantTest(t, server, restaurantDomain)

	bakeryProductID := createProductForProductFlow(t, server, bakeryDomain, bakeryAccessToken, `{"name":"Cake","sku":"BK-INV-001","price":999,"vat_percent":15}`, http.StatusCreated)
	restaurantProductID := createProductForProductFlow(t, server, restaurantDomain, restaurantAccessToken, `{"name":"Biryani","sku":"RST-INV-001","price":899,"vat_percent":5}`, http.StatusCreated)

	updateReq, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/inventory/%s", bakeryProductID), bytes.NewBufferString(`{"in_stock":40}`))
	require.NoError(t, err)
	updateReq.Host = bakeryDomain
	setBearerToken(updateReq, bakeryAccessToken)
	updateRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(updateRecorder, updateReq)
	require.Equal(t, http.StatusOK, updateRecorder.Code)

	crossTenantReq, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/inventory/%s", restaurantProductID), bytes.NewBufferString(`{"in_stock":10}`))
	require.NoError(t, err)
	crossTenantReq.Host = bakeryDomain
	setBearerToken(crossTenantReq, bakeryAccessToken)
	crossTenantRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(crossTenantRecorder, crossTenantReq)
	require.Equal(t, http.StatusUnauthorized, crossTenantRecorder.Code)
}
