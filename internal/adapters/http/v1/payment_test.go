package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/horiondreher/go-web-api-boilerplate/internal/utils"
	"github.com/stretchr/testify/require"
)

type createPaymentMethodResponseForTest struct {
	ID string `json:"id"`
}

type listPaymentMethodResponseForTest struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func Test_Payment_Methods_And_PayOrder_Flow(t *testing.T) {
	server, err := NewHTTPAdapter(testUserService)
	require.NoError(t, err)

	bakerySlug := "bakery-" + utils.RandomString(6)
	bakeryDomain := bakerySlug + ".localhost"
	bakeryID := createTenantForProductFlow(t, server, bakerySlug, bakeryDomain, "bakery")
	createSubscriptionForProductFlow(t, server, bakeryID, bakeryDomain, "active")
	bakeryAccessToken := createAccessTokenForTenantTest(t, server, bakeryDomain)

	productID := createProductForProductFlow(t, server, bakeryDomain, bakeryAccessToken, `{"name":"Brownie","sku":"BK-PAY-001","price":1299,"vat_percent":1}`, http.StatusCreated)
	seedInventoryForOrderTest(t, server, bakeryDomain, bakeryAccessToken, productID, 10)
	userUID := createUserForTenantCartTest(t, server, bakeryDomain)
	cartID := addItemAndGetCartIDForOrderTest(t, server, bakeryDomain, bakeryAccessToken, userUID, productID, 1)
	orderID, total := createOrderForOrderTest(t, server, bakeryDomain, bakeryAccessToken, fmt.Sprintf(`{"cart_id":"%s","payment_method_id":"cash","fulfillment_type":"pickup"}`, cartID), http.StatusCreated)

	cardMethodID := createPaymentMethodForPaymentTest(t, server, bakeryDomain, bakeryAccessToken, `{"type":"card","label":"Visa **** 4242","is_default":false}`, http.StatusCreated)
	listPaymentMethodsForPaymentTest(t, server, bakeryDomain, bakeryAccessToken)

	payOrderForPaymentTest(t, server, bakeryDomain, bakeryAccessToken, orderID, fmt.Sprintf(`{"payment_method_id":"%s","amount":%d}`, cardMethodID, total), http.StatusOK)
	payOrderForPaymentTest(t, server, bakeryDomain, bakeryAccessToken, orderID, fmt.Sprintf(`{"payment_method_id":"cash","amount":%d}`, total), http.StatusUnprocessableEntity)
}

func Test_Payment_WrongTenantMethod_And_InvalidAmount(t *testing.T) {
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

	productID := createProductForProductFlow(t, server, bakeryDomain, bakeryAccessToken, `{"name":"Croissant","sku":"BK-PAY-002","price":999,"vat_percent":1}`, http.StatusCreated)
	seedInventoryForOrderTest(t, server, bakeryDomain, bakeryAccessToken, productID, 10)
	userUID := createUserForTenantCartTest(t, server, bakeryDomain)
	cartID := addItemAndGetCartIDForOrderTest(t, server, bakeryDomain, bakeryAccessToken, userUID, productID, 1)
	orderID, total := createOrderForOrderTest(t, server, bakeryDomain, bakeryAccessToken, fmt.Sprintf(`{"cart_id":"%s","payment_method_id":"cash","fulfillment_type":"pickup"}`, cartID), http.StatusCreated)

	otherTenantCardID := createPaymentMethodForPaymentTest(t, server, restaurantDomain, restaurantAccessToken, `{"type":"card","label":"Mastercard **** 1234","is_default":false}`, http.StatusCreated)

	payOrderForPaymentTest(t, server, bakeryDomain, bakeryAccessToken, orderID, `{"payment_method_id":"cash","amount":0}`, http.StatusUnprocessableEntity)
	payOrderForPaymentTest(t, server, bakeryDomain, bakeryAccessToken, orderID, fmt.Sprintf(`{"payment_method_id":"%s","amount":%d}`, otherTenantCardID, total), http.StatusUnauthorized)
}

func createPaymentMethodForPaymentTest(t *testing.T, server *HTTPAdapter, host string, accessToken string, body string, expectedCode int) string {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/payment-methods", bytes.NewBufferString(body))
	require.NoError(t, err)
	req.Host = host
	setBearerToken(req, accessToken)
	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, req)
	require.Equal(t, expectedCode, recorder.Code)

	if expectedCode != http.StatusCreated {
		return ""
	}

	var response createPaymentMethodResponseForTest
	err = json.NewDecoder(recorder.Body).Decode(&response)
	require.NoError(t, err)
	require.NotEmpty(t, response.ID)
	return response.ID
}

func listPaymentMethodsForPaymentTest(t *testing.T, server *HTTPAdapter, host string, accessToken string) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "/api/v1/payment-methods", nil)
	require.NoError(t, err)
	req.Host = host
	setBearerToken(req, accessToken)
	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)

	var response []listPaymentMethodResponseForTest
	err = json.NewDecoder(recorder.Body).Decode(&response)
	require.NoError(t, err)
	require.NotEmpty(t, response)

	hasCash := false
	for _, paymentMethod := range response {
		if paymentMethod.ID == "cash" && paymentMethod.Type == "cash" {
			hasCash = true
			break
		}
	}
	require.True(t, hasCash)
}

func payOrderForPaymentTest(t *testing.T, server *HTTPAdapter, host string, accessToken string, orderID uuid.UUID, body string, expectedCode int) {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/orders/%s/pay", orderID), bytes.NewBufferString(body))
	require.NoError(t, err)
	req.Host = host
	setBearerToken(req, accessToken)
	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, req)
	require.Equal(t, expectedCode, recorder.Code)
}
