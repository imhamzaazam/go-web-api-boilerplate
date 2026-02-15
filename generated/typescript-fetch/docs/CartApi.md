# CartApi

All URIs are relative to *http://localhost:8080/api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createCartItem**](CartApi.md#createcartitemoperation) | **POST** /cart/items | Add item to active cart |



## createCartItem

> CreateCartItemResponse createCartItem(createCartItemRequest)

Add item to active cart

### Example

```ts
import {
  Configuration,
  CartApi,
} from '';
import type { CreateCartItemOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new CartApi();

  const body = {
    // CreateCartItemRequest
    createCartItemRequest: ...,
  } satisfies CreateCartItemOperationRequest;

  try {
    const data = await api.createCartItem(body);
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
| **createCartItemRequest** | [CreateCartItemRequest](CreateCartItemRequest.md) |  | |

### Return type

[**CreateCartItemResponse**](CreateCartItemResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Cart item created |  -  |
| **401** | API error |  -  |
| **409** | API error |  -  |
| **422** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

