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

type createOrderResponseForTest struct {
	ID    uuid.UUID `json:"id"`
	Total int64     `json:"total"`
}

type createCartItemResponseForOrderTest struct {
	CartID uuid.UUID `json:"cart_id"`
}

func Test_Order_CreateAndTransition_Flow(t *testing.T) {
	server, err := NewHTTPAdapter(testUserService)
	require.NoError(t, err)

	restaurantSlug := "restaurant-" + utils.RandomString(6)
	restaurantDomain := restaurantSlug + ".localhost"
	restaurantID := createTenantForProductFlow(t, server, restaurantSlug, restaurantDomain, "restaurant")
	createSubscriptionForProductFlow(t, server, restaurantID, restaurantDomain, "active")
	restaurantAccessToken := createAccessTokenForTenantTest(t, server, restaurantDomain)

	productID := createProductForProductFlow(t, server, restaurantDomain, restaurantAccessToken, `{"name":"Chicken Chowmein","sku":"RST-ORD-001","price":1200,"vat_percent":10}`, http.StatusCreated)
	seedInventoryForOrderTest(t, server, restaurantDomain, restaurantAccessToken, productID, 20)

	userUID := createUserForTenantCartTest(t, server, restaurantDomain)
	cartID := addItemAndGetCartIDForOrderTest(t, server, restaurantDomain, restaurantAccessToken, userUID, productID, 2)

	orderID, _ := createOrderForOrderTest(t, server, restaurantDomain, restaurantAccessToken, fmt.Sprintf(`{"cart_id":"%s","payment_method_id":"cash","fulfillment_type":"delivery","location":{"address_line":"Clifton Block 5","city":"Karachi","lat":24.8138,"lng":67.0306}}`, cartID), http.StatusCreated)

	patchOrderStatusForOrderTest(t, server, restaurantDomain, restaurantAccessToken, orderID, "confirmed", http.StatusOK)
	patchOrderStatusForOrderTest(t, server, restaurantDomain, restaurantAccessToken, orderID, "out_for_delivery", http.StatusOK)
	patchOrderStatusForOrderTest(t, server, restaurantDomain, restaurantAccessToken, orderID, "completed", http.StatusOK)
	patchOrderStatusForOrderTest(t, server, restaurantDomain, restaurantAccessToken, orderID, "completed", http.StatusUnprocessableEntity)

	secondCartID := addItemAndGetCartIDForOrderTest(t, server, restaurantDomain, restaurantAccessToken, userUID, productID, 1)
	secondOrderID, _ := createOrderForOrderTest(t, server, restaurantDomain, restaurantAccessToken, fmt.Sprintf(`{"cart_id":"%s","payment_method_id":"cash","fulfillment_type":"pickup"}`, secondCartID), http.StatusCreated)
	patchOrderStatusForOrderTest(t, server, restaurantDomain, restaurantAccessToken, secondOrderID, "refunded", http.StatusUnprocessableEntity)
}

func Test_Order_Create_ValidationAndTenantIsolation(t *testing.T) {
	server, err := NewHTTPAdapter(testUserService)
	require.NoError(t, err)

	restaurantSlug := "restaurant-" + utils.RandomString(6)
	restaurantDomain := restaurantSlug + ".localhost"
	restaurantID := createTenantForProductFlow(t, server, restaurantSlug, restaurantDomain, "restaurant")
	createSubscriptionForProductFlow(t, server, restaurantID, restaurantDomain, "active")
	restaurantAccessToken := createAccessTokenForTenantTest(t, server, restaurantDomain)

	bakerySlug := "bakery-" + utils.RandomString(6)
	bakeryDomain := bakerySlug + ".localhost"
	bakeryID := createTenantForProductFlow(t, server, bakerySlug, bakeryDomain, "bakery")
	createSubscriptionForProductFlow(t, server, bakeryID, bakeryDomain, "active")
	bakeryAccessToken := createAccessTokenForTenantTest(t, server, bakeryDomain)

	productID := createProductForProductFlow(t, server, restaurantDomain, restaurantAccessToken, `{"name":"Zinger Burger","sku":"RST-ORD-002","price":900,"vat_percent":8}`, http.StatusCreated)
	seedInventoryForOrderTest(t, server, restaurantDomain, restaurantAccessToken, productID, 10)

	userUID := createUserForTenantCartTest(t, server, restaurantDomain)
	cartID := addItemAndGetCartIDForOrderTest(t, server, restaurantDomain, restaurantAccessToken, userUID, productID, 1)

	_, _ = createOrderForOrderTest(t, server, restaurantDomain, restaurantAccessToken, fmt.Sprintf(`{"cart_id":"%s","payment_method_id":"cash","fulfillment_type":"dinein"}`, cartID), http.StatusUnprocessableEntity)
	_, _ = createOrderForOrderTest(t, server, restaurantDomain, restaurantAccessToken, fmt.Sprintf(`{"cart_id":"%s","payment_method_id":"cash","fulfillment_type":"delivery"}`, cartID), http.StatusUnprocessableEntity)
	_, _ = createOrderForOrderTest(t, server, restaurantDomain, restaurantAccessToken, fmt.Sprintf(`{"cart_id":"%s","payment_method_id":"cash","fulfillment_type":"delivery","location":{"address_line":"Unknown Area","city":"Karachi","lat":25.8,"lng":66.5}}`, cartID), http.StatusUnprocessableEntity)
	_, _ = createOrderForOrderTest(t, server, bakeryDomain, bakeryAccessToken, fmt.Sprintf(`{"cart_id":"%s","payment_method_id":"cash","fulfillment_type":"pickup"}`, cartID), http.StatusUnauthorized)
}

func seedInventoryForOrderTest(t *testing.T, server *HTTPAdapter, host string, accessToken string, productID string, inStock int32) {
	t.Helper()
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/inventory/%s", productID), bytes.NewBufferString(fmt.Sprintf(`{"in_stock":%d}`, inStock)))
	require.NoError(t, err)
	req.Host = host
	setBearerToken(req, accessToken)
	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func addItemAndGetCartIDForOrderTest(t *testing.T, server *HTTPAdapter, host string, accessToken string, userUID uuid.UUID, productID string, quantity int32) uuid.UUID {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/cart/items", bytes.NewBufferString(fmt.Sprintf(`{"user_uid":"%s","product_id":"%s","quantity":%d}`, userUID, productID, quantity)))
	require.NoError(t, err)
	req.Host = host
	setBearerToken(req, accessToken)
	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusOK, recorder.Code)

	var response createCartItemResponseForOrderTest
	err = json.NewDecoder(recorder.Body).Decode(&response)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, response.CartID)

	return response.CartID
}

func createOrderForOrderTest(t *testing.T, server *HTTPAdapter, host string, accessToken string, body string, expectedCode int) (uuid.UUID, int64) {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBufferString(body))
	require.NoError(t, err)
	req.Host = host
	setBearerToken(req, accessToken)
	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, req)
	require.Equal(t, expectedCode, recorder.Code)

	if expectedCode != http.StatusCreated {
		return uuid.Nil, 0
	}

	var response createOrderResponseForTest
	err = json.NewDecoder(recorder.Body).Decode(&response)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, response.ID)

	return response.ID, response.Total
}

func patchOrderStatusForOrderTest(t *testing.T, server *HTTPAdapter, host string, accessToken string, orderID uuid.UUID, status string, expectedCode int) {
	t.Helper()
	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/orders/%s/status", orderID), bytes.NewBufferString(fmt.Sprintf(`{"status":"%s"}`, status)))
	require.NoError(t, err)
	req.Host = host
	setBearerToken(req, accessToken)
	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, req)
	require.Equal(t, expectedCode, recorder.Code)
}
