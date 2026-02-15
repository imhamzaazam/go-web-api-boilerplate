# TenantsApi

All URIs are relative to *http://localhost:8080/api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createBranch**](TenantsApi.md#createbranchoperation) | **POST** /branches | Create branch |
| [**createSubscription**](TenantsApi.md#createsubscriptionoperation) | **POST** /subscriptions | Create subscription |
| [**createTenant**](TenantsApi.md#createtenantoperation) | **POST** /tenants | Create tenant |



## createBranch

> CreateBranchResponse createBranch(createBranchRequest)

Create branch

### Example

```ts
import {
  Configuration,
  TenantsApi,
} from '';
import type { CreateBranchOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new TenantsApi();

  const body = {
    // CreateBranchRequest
    createBranchRequest: ...,
  } satisfies CreateBranchOperationRequest;

  try {
    const data = await api.createBranch(body);
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
| **createBranchRequest** | [CreateBranchRequest](CreateBranchRequest.md) |  | |

### Return type

[**CreateBranchResponse**](CreateBranchResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Branch created |  -  |
| **409** | API error |  -  |
| **422** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## createSubscription

> CreateSubscriptionResponse createSubscription(createSubscriptionRequest)

Create subscription

### Example

```ts
import {
  Configuration,
  TenantsApi,
} from '';
import type { CreateSubscriptionOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new TenantsApi();

  const body = {
    // CreateSubscriptionRequest
    createSubscriptionRequest: ...,
  } satisfies CreateSubscriptionOperationRequest;

  try {
    const data = await api.createSubscription(body);
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
| **createSubscriptionRequest** | [CreateSubscriptionRequest](CreateSubscriptionRequest.md) |  | |

### Return type

[**CreateSubscriptionResponse**](CreateSubscriptionResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Subscription created |  -  |
| **409** | API error |  -  |
| **422** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## createTenant

> CreateTenantResponse createTenant(createTenantRequest)

Create tenant

### Example

```ts
import {
  Configuration,
  TenantsApi,
} from '';
import type { CreateTenantOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new TenantsApi();

  const body = {
    // CreateTenantRequest
    createTenantRequest: ...,
  } satisfies CreateTenantOperationRequest;

  try {
    const data = await api.createTenant(body);
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
| **createTenantRequest** | [CreateTenantRequest](CreateTenantRequest.md) |  | |

### Return type

[**CreateTenantResponse**](CreateTenantResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Tenant created |  -  |
| **409** | API error |  -  |
| **422** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

