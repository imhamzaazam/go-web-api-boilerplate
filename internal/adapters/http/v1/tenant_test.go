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

type loginResponseForTenantTest struct {
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

type createTenantResponseForTest struct {
	ID string `json:"id"`
}

func Test_Tenant_Create_Success(t *testing.T) {
	server, err := NewHTTPAdapter(testUserService)
	require.NoError(t, err)

	slug := "tenant-" + utils.RandomString(6)
	domain := slug + ".localhost"
	body := fmt.Sprintf(`{"name":"%s","slug":"%s","type":"bakery","domain":"%s"}`, slug, slug, domain)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewBufferString(body))
	require.NoError(t, err)
	req.Host = "platform.localhost"

	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusCreated, recorder.Code)
}

func Test_Tenant_Isolation_CrossTenantAccessDenied(t *testing.T) {
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

	productID := createProductForProductFlow(t, server, restaurantDomain, restaurantAccessToken, `{"name":"Chicken Karahi","sku":"RST-CRS-001","price":1599,"vat_percent":5}`, http.StatusCreated)

	crossTenantReq, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/inventory/%s", productID), bytes.NewBufferString(`{"in_stock":5}`))
	require.NoError(t, err)
	crossTenantReq.Host = bakeryDomain
	setBearerToken(crossTenantReq, bakeryAccessToken)

	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, crossTenantReq)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func Test_Tenant_Subscription_Branch_TDDMatrix(t *testing.T) {
	server, err := NewHTTPAdapter(testUserService)
	require.NoError(t, err)

	slug := "tenant-" + utils.RandomString(6)
	domain := slug + ".localhost"
	tenantBody := fmt.Sprintf(`{"name":"%s","slug":"%s","type":"bakery","domain":"%s"}`, slug, slug, domain)

	tenantReq, err := http.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewBufferString(tenantBody))
	require.NoError(t, err)
	tenantReq.Host = "platform.localhost"

	tenantRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(tenantRecorder, tenantReq)
	require.Equal(t, http.StatusCreated, tenantRecorder.Code)

	var tenantRes createTenantResponseForTest
	err = json.NewDecoder(tenantRecorder.Body).Decode(&tenantRes)
	require.NoError(t, err)
	require.NotEmpty(t, tenantRes.ID)

	duplicateReq, err := http.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewBufferString(fmt.Sprintf(`{"name":"dup-%s","slug":"dup-%s","type":"bakery","domain":"%s"}`, slug, slug, domain)))
	require.NoError(t, err)
	duplicateReq.Host = "platform.localhost"

	duplicateRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(duplicateRecorder, duplicateReq)
	require.Equal(t, http.StatusConflict, duplicateRecorder.Code)

	branchReq, err := http.NewRequest(http.MethodPost, "/api/v1/branches", bytes.NewBufferString(fmt.Sprintf(`{"tenant_id":"%s","name":"Main Branch","code":"BR001"}`, tenantRes.ID)))
	require.NoError(t, err)
	branchReq.Host = domain

	branchRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(branchRecorder, branchReq)
	require.Equal(t, http.StatusCreated, branchRecorder.Code)

	branchDuplicateReq, err := http.NewRequest(http.MethodPost, "/api/v1/branches", bytes.NewBufferString(fmt.Sprintf(`{"tenant_id":"%s","name":"Second Branch","code":"BR001"}`, tenantRes.ID)))
	require.NoError(t, err)
	branchDuplicateReq.Host = domain

	branchDuplicateRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(branchDuplicateRecorder, branchDuplicateReq)
	require.Equal(t, http.StatusConflict, branchDuplicateRecorder.Code)

	subscriptionReq, err := http.NewRequest(http.MethodPost, "/api/v1/subscriptions", bytes.NewBufferString(fmt.Sprintf(`{"tenant_id":"%s","plan":"starter","status":"active"}`, tenantRes.ID)))
	require.NoError(t, err)
	subscriptionReq.Host = domain

	subscriptionRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(subscriptionRecorder, subscriptionReq)
	require.Equal(t, http.StatusCreated, subscriptionRecorder.Code)

	invalidSubscriptionReq, err := http.NewRequest(http.MethodPost, "/api/v1/subscriptions", bytes.NewBufferString(fmt.Sprintf(`{"tenant_id":"%s","plan":"starter","status":"unknown"}`, tenantRes.ID)))
	require.NoError(t, err)
	invalidSubscriptionReq.Host = domain

	invalidSubscriptionRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(invalidSubscriptionRecorder, invalidSubscriptionReq)
	require.Equal(t, http.StatusUnprocessableEntity, invalidSubscriptionRecorder.Code)
}

func Test_Login_HostTenantResolution_TDDMatrix(t *testing.T) {
	server, err := NewHTTPAdapter(testUserService)
	require.NoError(t, err)

	bakerySlug := "bakery-" + utils.RandomString(6)
	bakeryDomain := bakerySlug + ".localhost"
	bakeryID := createTenantForProductFlow(t, server, bakerySlug, bakeryDomain, "bakery")
	createSubscriptionForProductFlow(t, server, bakeryID, bakeryDomain, "active")

	email := utils.RandomEmail()
	password := utils.RandomString(10)
	fullName := utils.RandomString(8)

	createBody := fmt.Sprintf(`{"full_name":"%s","email":"%s","password":"%s"}`, fullName, email, password)
	createReq, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(createBody))
	require.NoError(t, err)
	createReq.Host = bakeryDomain
	createRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(createRecorder, createReq)
	require.Equal(t, http.StatusCreated, createRecorder.Code)

	loginBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	t.Run("login success by host tenant", func(t *testing.T) {
		req, reqErr := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(loginBody))
		require.NoError(t, reqErr)
		req.Host = bakeryDomain
		recorder := httptest.NewRecorder()
		server.router.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("login unknown host tenant failure", func(t *testing.T) {
		req, reqErr := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(loginBody))
		require.NoError(t, reqErr)
		req.Host = "unknown.localhost"
		recorder := httptest.NewRecorder()
		server.router.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("login invalid password failure", func(t *testing.T) {
		req, reqErr := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(fmt.Sprintf(`{"email":"%s","password":"bad-password"}`, email)))
		require.NoError(t, reqErr)
		req.Host = bakeryDomain
		recorder := httptest.NewRecorder()
		server.router.ServeHTTP(recorder, req)
		require.Equal(t, http.StatusUnauthorized, recorder.Code)
	})
}

func Test_Login_TokenContainsTenantID(t *testing.T) {
	server, err := NewHTTPAdapter(testUserService)
	require.NoError(t, err)

	bakerySlug := "bakery-" + utils.RandomString(6)
	bakeryDomain := bakerySlug + ".localhost"
	bakeryID := createTenantForProductFlow(t, server, bakerySlug, bakeryDomain, "bakery")
	createSubscriptionForProductFlow(t, server, bakeryID, bakeryDomain, "active")

	fullName := utils.RandomString(8)
	email := utils.RandomEmail()
	password := utils.RandomString(10)

	createBody := fmt.Sprintf(`{"full_name":"%s","email":"%s","password":"%s"}`, fullName, email, password)
	createReq, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(createBody))
	require.NoError(t, err)
	createReq.Host = bakeryDomain

	createRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(createRecorder, createReq)
	require.Equal(t, http.StatusCreated, createRecorder.Code)

	loginBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)
	loginReq, err := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(loginBody))
	require.NoError(t, err)
	loginReq.Host = bakeryDomain

	loginRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(loginRecorder, loginReq)
	require.Equal(t, http.StatusOK, loginRecorder.Code)

	var loginRes loginResponseForTenantTest
	err = json.NewDecoder(loginRecorder.Body).Decode(&loginRes)
	require.NoError(t, err)
	require.NotEmpty(t, loginRes.AccessToken)
	require.Equal(t, email, loginRes.Email)

	payload, verifyErr := server.tokenMaker.VerifyToken(loginRes.AccessToken)
	require.Nil(t, verifyErr)
	require.Equal(t, email, payload.Email)
	require.Equal(t, "employee", payload.Role)
	require.Equal(t, bakerySlug, payload.TenantSlug)
	require.Equal(t, uuid.MustParse(bakeryID), payload.TenantID)
	require.Equal(t, "active", payload.SubscriptionStatus)
}
