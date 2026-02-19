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

type loginResponseForAuthHelperTest struct {
	AccessToken string `json:"access_token"`
}

func createAccessTokenForTenantTest(t *testing.T, server *HTTPAdapter, host string) string {
	t.Helper()

	password := utils.RandomString(10)
	createBody := fmt.Sprintf(`{"full_name":"%s","email":"%s","password":"%s"}`, utils.RandomString(8), utils.RandomEmail(), password)

	createReq, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(createBody))
	require.NoError(t, err)
	createReq.Host = host

	createRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(createRecorder, createReq)
	require.Equal(t, http.StatusCreated, createRecorder.Code)

	var createdUser CreateUserResponseDto
	err = json.NewDecoder(createRecorder.Body).Decode(&createdUser)
	require.NoError(t, err)
	require.NotEmpty(t, createdUser.Email)

	loginBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, createdUser.Email, password)
	loginReq, err := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBufferString(loginBody))
	require.NoError(t, err)
	loginReq.Host = host

	loginRecorder := httptest.NewRecorder()
	server.router.ServeHTTP(loginRecorder, loginReq)
	require.Equal(t, http.StatusOK, loginRecorder.Code)

	var loginRes loginResponseForAuthHelperTest
	err = json.NewDecoder(loginRecorder.Body).Decode(&loginRes)
	require.NoError(t, err)
	require.NotEmpty(t, loginRes.AccessToken)

	return loginRes.AccessToken
}

func setBearerToken(req *http.Request, accessToken string) {
	req.Header.Set("Authorization", "Bearer "+accessToken)
}
