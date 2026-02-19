# InventoryApi

All URIs are relative to *http://localhost:8080/api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**upsertInventoryForProduct**](InventoryApi.md#upsertinventoryforproduct) | **PUT** /inventory/{id} | Set inventory stock for product |



## upsertInventoryForProduct

> UpsertInventoryResponse upsertInventoryForProduct(id, upsertInventoryRequest)

Set inventory stock for product

### Example

```ts
import {
  Configuration,
  InventoryApi,
} from '';
import type { UpsertInventoryForProductRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new InventoryApi();

  const body = {
    // string | Product UUID
    id: 38400000-8cf0-11bd-b23e-10b96e4ef00d,
    // UpsertInventoryRequest
    upsertInventoryRequest: ...,
  } satisfies UpsertInventoryForProductRequest;

  try {
    const data = await api.upsertInventoryForProduct(body);
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
| **id** | `string` | Product UUID | [Defaults to `undefined`] |
| **upsertInventoryRequest** | [UpsertInventoryRequest](UpsertInventoryRequest.md) |  | |

### Return type

[**UpsertInventoryResponse**](UpsertInventoryResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Inventory updated |  -  |
| **401** | API error |  -  |
| **422** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

