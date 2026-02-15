package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/horiondreher/go-web-api-boilerplate/internal/utils"
	"github.com/stretchr/testify/require"
)

type createProductResponseForTest struct {
	ID string `json:"id"`
}

func Test_Product_Discount_Addon_TenantAwareFlow(t *testing.T) {
	server, err := NewHTTPAdapter(testUserService)
	require.NoError(t, err)

	bakerySlug := "bakery-" + utils.RandomString(6)
	bakeryDomain := bakerySlug + ".localhost"
	bakeryID := createTenantForProductFlow(t, server, bakerySlug, bakeryDomain, "bakery")
	createSubscriptionForProductFlow(t, server, bakeryID, bakeryDomain, "active")
	bakeryAccessToken := createAccessTokenForTenantTest(t, server, bakeryDomain)

	productBody := `{"name":"Chocolate Fudge Cake","sku":"BK-CAKE-001","price":1299,"vat_percent":15,"is_preorder":true,"made_to_order":true}`
	productID := createProductForProductFlow(t, server, bakeryDomain, bakeryAccessToken, productBody, http.StatusCreated)

	duplicateBody := `{"name":"Vanilla Cake","sku":"BK-CAKE-001","price":1000,"vat_percent":15}`
	createProductForProductFlow(t, server, bakeryDomain, bakeryAccessToken, duplicateBody, http.StatusConflict)

	discountBody := `{"code":"PAK001","name":"Weekend Cake Offer","type":"percentage","value":10,"starts_at":"2026-02-14T00:00:00Z","ends_at":"2026-02-28T23:59:59Z"}`
	discountReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/products/%s/discounts", productID), bytes.NewBufferString(discountBody))
	require.NoError(t, err)
	discountReq.Host = bakeryDomain
	setBearerToken(discountReq, bakeryAccessToken)
	discountRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(discountRecorder, discountReq)
	require.Equal(t, http.StatusCreated, discountRecorder.Code)

	addonBody := `{"name":"Extra Cream","price":199}`
	addonReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/products/%s/addons", productID), bytes.NewBufferString(addonBody))
	require.NoError(t, err)
	addonReq.Host = bakeryDomain
	setBearerToken(addonReq, bakeryAccessToken)
	addonRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(addonRecorder, addonReq)
	require.Equal(t, http.StatusCreated, addonRecorder.Code)

	pharmacySlug := "pharmacy-" + utils.RandomString(6)
	pharmacyDomain := pharmacySlug + ".localhost"
	pharmacyID := createTenantForProductFlow(t, server, pharmacySlug, pharmacyDomain, "pharmacy")
	createSubscriptionForProductFlow(t, server, pharmacyID, pharmacyDomain, "active")
	pharmacyAccessToken := createAccessTokenForTenantTest(t, server, pharmacyDomain)

	pharmacyProductBody := `{"name":"Amoxicillin 500mg","sku":"PH-DRUG-001","price":499,"vat_percent":12,"requires_prescription":true}`
	pharmacyProductID := createProductForProductFlow(t, server, pharmacyDomain, pharmacyAccessToken, pharmacyProductBody, http.StatusCreated)

	pharmacyAddonReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/products/%s/addons", pharmacyProductID), bytes.NewBufferString(`{"name":"Bundle","price":99}`))
	require.NoError(t, err)
	pharmacyAddonReq.Host = pharmacyDomain
	setBearerToken(pharmacyAddonReq, pharmacyAccessToken)
	pharmacyAddonRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(pharmacyAddonRecorder, pharmacyAddonReq)
	require.Equal(t, http.StatusUnprocessableEntity, pharmacyAddonRecorder.Code)

	restaurantSlug := "restaurant-" + utils.RandomString(6)
	restaurantDomain := restaurantSlug + ".localhost"
	restaurantID := createTenantForProductFlow(t, server, restaurantSlug, restaurantDomain, "restaurant")
	createSubscriptionForProductFlow(t, server, restaurantID, restaurantDomain, "active")
	restaurantAccessToken := createAccessTokenForTenantTest(t, server, restaurantDomain)

	restaurantMadeToOrderBody := `{"name":"Chicken Karahi","sku":"RST-KARAHI-001","price":1599,"vat_percent":5,"made_to_order":true}`
	createProductForProductFlow(t, server, restaurantDomain, restaurantAccessToken, restaurantMadeToOrderBody, http.StatusUnprocessableEntity)
}

func createTenantForProductFlow(t *testing.T, server *HTTPAdapter, slug string, domain string, tenantType string) string {
	t.Helper()
	body := fmt.Sprintf(`{"name":"%s","slug":"%s","type":"%s","domain":"%s"}`, slug, slug, tenantType, domain)
	req, err := http.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewBufferString(body))
	require.NoError(t, err)
	req.Host = "platform.localhost"
	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusCreated, recorder.Code)

	var tenantRes createTenantResponseForTest
	err = json.NewDecoder(recorder.Body).Decode(&tenantRes)
	require.NoError(t, err)
	require.NotEmpty(t, tenantRes.ID)
	return tenantRes.ID
}

func createSubscriptionForProductFlow(t *testing.T, server *HTTPAdapter, tenantID string, host string, status string) {
	t.Helper()
	body := fmt.Sprintf(`{"tenant_id":"%s","plan":"starter","status":"%s"}`, tenantID, status)
	req, err := http.NewRequest(http.MethodPost, "/api/v1/subscriptions", bytes.NewBufferString(body))
	require.NoError(t, err)
	req.Host = host
	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusCreated, recorder.Code)
}

func createProductForProductFlow(t *testing.T, server *HTTPAdapter, host string, accessToken string, body string, expectedCode int) string {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBufferString(body))
	require.NoError(t, err)
	req.Host = host
	setBearerToken(req, accessToken)
	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, req)
	require.Equal(t, expectedCode, recorder.Code)

	if expectedCode != http.StatusCreated {
		return ""
	}

	var productRes createProductResponseForTest
	err = json.NewDecoder(recorder.Body).Decode(&productRes)
	require.NoError(t, err)
	require.NotEmpty(t, productRes.ID)
	return productRes.ID
}
