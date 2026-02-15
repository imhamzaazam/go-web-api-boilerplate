# AuthApi

All URIs are relative to *http://localhost:8080/api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**loginUser**](AuthApi.md#loginuser) | **POST** /login | Login user |
| [**renewAccessToken**](AuthApi.md#renewaccesstoken) | **POST** /renew-token | Renew access token |



## loginUser

> LoginResponse loginUser(loginRequest)

Login user

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '';
import type { LoginUserRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new AuthApi();

  const body = {
    // LoginRequest
    loginRequest: ...,
  } satisfies LoginUserRequest;

  try {
    const data = await api.loginUser(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **loginRequest** | [LoginRequest](LoginRequest.md) |  | |

### Return type

[**LoginResponse**](LoginResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Login successful |  -  |
| **400** | API error |  -  |
| **401** | API error |  -  |
| **404** | API error |  -  |
| **500** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## renewAccessToken

> RenewTokenResponse renewAccessToken(renewTokenRequest)

Renew access token

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '';
import type { RenewAccessTokenRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new AuthApi();

  const body = {
    // RenewTokenRequest
    renewTokenRequest: ...,
  } satisfies RenewAccessTokenRequest;

  try {
    const data = await api.renewAccessToken(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **renewTokenRequest** | [RenewTokenRequest](RenewTokenRequest.md) |  | |

### Return type

[**RenewTokenResponse**](RenewTokenResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Access token renewed |  -  |
| **400** | API error |  -  |
| **401** | API error |  -  |
| **404** | API error |  -  |
| **500** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

