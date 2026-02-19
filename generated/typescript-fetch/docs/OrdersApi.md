# OrdersApi

All URIs are relative to *http://localhost:8080/api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createOrder**](OrdersApi.md#createorderoperation) | **POST** /orders | Create order from active cart |
| [**updateOrderStatus**](OrdersApi.md#updateorderstatusoperation) | **PATCH** /orders/{id}/status | Update order status |



## createOrder

> OrderResponse createOrder(createOrderRequest)

Create order from active cart

### Example

```ts
import {
  Configuration,
  OrdersApi,
} from '';
import type { CreateOrderOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new OrdersApi();

  const body = {
    // CreateOrderRequest
    createOrderRequest: ...,
  } satisfies CreateOrderOperationRequest;

  try {
    const data = await api.createOrder(body);
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
| **createOrderRequest** | [CreateOrderRequest](CreateOrderRequest.md) |  | |

### Return type

[**OrderResponse**](OrderResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Order created |  -  |
| **401** | API error |  -  |
| **402** | API error |  -  |
| **422** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## updateOrderStatus

> OrderResponse updateOrderStatus(id, updateOrderStatusRequest)

Update order status

### Example

```ts
import {
  Configuration,
  OrdersApi,
} from '';
import type { UpdateOrderStatusOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new OrdersApi();

  const body = {
    // string | Order UUID
    id: 38400000-8cf0-11bd-b23e-10b96e4ef00d,
    // UpdateOrderStatusRequest
    updateOrderStatusRequest: ...,
  } satisfies UpdateOrderStatusOperationRequest;

  try {
    const data = await api.updateOrderStatus(body);
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
| **id** | `string` | Order UUID | [Defaults to `undefined`] |
| **updateOrderStatusRequest** | [UpdateOrderStatusRequest](UpdateOrderStatusRequest.md) |  | |

### Return type

[**OrderResponse**](OrderResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Order status updated |  -  |
| **401** | API error |  -  |
| **422** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

