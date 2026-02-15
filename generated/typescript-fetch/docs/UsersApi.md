# UsersApi

All URIs are relative to *http://localhost:8080/api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createUser**](UsersApi.md#createuseroperation) | **POST** /users | Create user |
| [**getUserByUID**](UsersApi.md#getuserbyuid) | **GET** /user/{uid} | Get user by UID |



## createUser

> UserResponse createUser(createUserRequest)

Create user

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { CreateUserOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new UsersApi();

  const body = {
    // CreateUserRequest
    createUserRequest: ...,
  } satisfies CreateUserOperationRequest;

  try {
    const data = await api.createUser(body);
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
| **createUserRequest** | [CreateUserRequest](CreateUserRequest.md) |  | |

### Return type

[**UserResponse**](UserResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | User created |  -  |
| **400** | API error |  -  |
| **409** | API error |  -  |
| **500** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getUserByUID

> UserResponse getUserByUID(uid)

Get user by UID

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { GetUserByUIDRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // Configure HTTP bearer authorization: bearerAuth
    accessToken: "YOUR BEARER TOKEN",
  });
  const api = new UsersApi(config);

  const body = {
    // string | User UUID
    uid: 38400000-8cf0-11bd-b23e-10b96e4ef00d,
  } satisfies GetUserByUIDRequest;

  try {
    const data = await api.getUserByUID(body);
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
| **uid** | `string` | User UUID | [Defaults to `undefined`] |

### Return type

[**UserResponse**](UserResponse.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | User found |  -  |
| **401** | API error |  -  |
| **404** | API error |  -  |
| **500** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

