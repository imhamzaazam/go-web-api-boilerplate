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

type createUserResponseForCartTest struct {
	UID uuid.UUID `json:"uid"`
}

func Test_Cart_CreateItem_TenantAwareFlow(t *testing.T) {
	server, err := NewHTTPAdapter(testUserService)
	require.NoError(t, err)

	bakerySlug := "bakery-" + utils.RandomString(6)
	bakeryDomain := bakerySlug + ".localhost"
	bakeryID := createTenantForProductFlow(t, server, bakerySlug, bakeryDomain, "bakery")
	createSubscriptionForProductFlow(t, server, bakeryID, bakeryDomain, "active")
	bakeryAccessToken := createAccessTokenForTenantTest(t, server, bakeryDomain)

	productBody := `{"name":"Chocolate Fudge Cake","sku":"BK-CART-001","price":1299,"vat_percent":15}`
	productID := createProductForProductFlow(t, server, bakeryDomain, bakeryAccessToken, productBody, http.StatusCreated)

	inventoryReq, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/inventory/%s", productID), bytes.NewBufferString(`{"in_stock":5}`))
	require.NoError(t, err)
	inventoryReq.Host = bakeryDomain
	setBearerToken(inventoryReq, bakeryAccessToken)
	inventoryRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(inventoryRecorder, inventoryReq)
	require.Equal(t, http.StatusOK, inventoryRecorder.Code)

	userUID := createUserForTenantCartTest(t, server, bakeryDomain)

	addReq, err := http.NewRequest(http.MethodPost, "/api/v1/cart/items", bytes.NewBufferString(fmt.Sprintf(`{"user_uid":"%s","product_id":"%s","quantity":2,"note":"Happy birthday"}`, userUID, productID)))
	require.NoError(t, err)
	addReq.Host = bakeryDomain
	setBearerToken(addReq, bakeryAccessToken)
	addRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(addRecorder, addReq)
	require.Equal(t, http.StatusOK, addRecorder.Code)

	insufficientReq, err := http.NewRequest(http.MethodPost, "/api/v1/cart/items", bytes.NewBufferString(fmt.Sprintf(`{"user_uid":"%s","product_id":"%s","quantity":10}`, userUID, productID)))
	require.NoError(t, err)
	insufficientReq.Host = bakeryDomain
	setBearerToken(insufficientReq, bakeryAccessToken)
	insufficientRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(insufficientRecorder, insufficientReq)
	require.Equal(t, http.StatusConflict, insufficientRecorder.Code)
}

func createUserForTenantCartTest(t *testing.T, server *HTTPAdapter, host string) uuid.UUID {
	t.Helper()
	body := fmt.Sprintf(`{"full_name":"%s","email":"%s","password":"%s"}`, utils.RandomString(8), utils.RandomEmail(), utils.RandomString(10))
	req, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(body))
	require.NoError(t, err)
	req.Host = host
	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, req)
	require.Equal(t, http.StatusCreated, recorder.Code)

	var userRes createUserResponseForCartTest
	err = json.NewDecoder(recorder.Body).Decode(&userRes)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, userRes.UID)

	return userRes.UID
}
