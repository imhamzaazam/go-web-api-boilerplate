# PaymentsApi

All URIs are relative to *http://localhost:8080/api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createPaymentMethod**](PaymentsApi.md#createpaymentmethodoperation) | **POST** /payment-methods | Create payment method |
| [**listPaymentMethods**](PaymentsApi.md#listpaymentmethods) | **GET** /payment-methods | List payment methods |
| [**payOrder**](PaymentsApi.md#payorderoperation) | **POST** /orders/{id}/pay | Pay order |



## createPaymentMethod

> PaymentMethodResponse createPaymentMethod(createPaymentMethodRequest)

Create payment method

### Example

```ts
import {
  Configuration,
  PaymentsApi,
} from '';
import type { CreatePaymentMethodOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PaymentsApi();

  const body = {
    // CreatePaymentMethodRequest
    createPaymentMethodRequest: ...,
  } satisfies CreatePaymentMethodOperationRequest;

  try {
    const data = await api.createPaymentMethod(body);
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
| **createPaymentMethodRequest** | [CreatePaymentMethodRequest](CreatePaymentMethodRequest.md) |  | |

### Return type

[**PaymentMethodResponse**](PaymentMethodResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Payment method created |  -  |
| **401** | API error |  -  |
| **422** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listPaymentMethods

> Array&lt;PaymentMethodResponse&gt; listPaymentMethods()

List payment methods

### Example

```ts
import {
  Configuration,
  PaymentsApi,
} from '';
import type { ListPaymentMethodsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PaymentsApi();

  try {
    const data = await api.listPaymentMethods();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**Array&lt;PaymentMethodResponse&gt;**](PaymentMethodResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Payment methods listed |  -  |
| **401** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## payOrder

> PayOrderResponse payOrder(id, payOrderRequest)

Pay order

### Example

```ts
import {
  Configuration,
  PaymentsApi,
} from '';
import type { PayOrderOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PaymentsApi();

  const body = {
    // string | Order UUID
    id: 38400000-8cf0-11bd-b23e-10b96e4ef00d,
    // PayOrderRequest
    payOrderRequest: ...,
  } satisfies PayOrderOperationRequest;

  try {
    const data = await api.payOrder(body);
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
| **payOrderRequest** | [PayOrderRequest](PayOrderRequest.md) |  | |

### Return type

[**PayOrderResponse**](PayOrderResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Order paid |  -  |
| **401** | API error |  -  |
| **402** | API error |  -  |
| **422** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

