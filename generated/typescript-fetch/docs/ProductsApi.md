# ProductsApi

All URIs are relative to *http://localhost:8080/api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createDiscount**](ProductsApi.md#creatediscountoperation) | **POST** /products/{id}/discounts | Create product discount |
| [**createProduct**](ProductsApi.md#createproductoperation) | **POST** /products | Create product |
| [**createProductAddon**](ProductsApi.md#createproductaddonoperation) | **POST** /products/{id}/addons | Create product add-on |



## createDiscount

> CreateDiscountResponse createDiscount(id, createDiscountRequest)

Create product discount

### Example

```ts
import {
  Configuration,
  ProductsApi,
} from '';
import type { CreateDiscountOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ProductsApi();

  const body = {
    // string | Product UUID
    id: 38400000-8cf0-11bd-b23e-10b96e4ef00d,
    // CreateDiscountRequest
    createDiscountRequest: ...,
  } satisfies CreateDiscountOperationRequest;

  try {
    const data = await api.createDiscount(body);
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
| **createDiscountRequest** | [CreateDiscountRequest](CreateDiscountRequest.md) |  | |

### Return type

[**CreateDiscountResponse**](CreateDiscountResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Discount created |  -  |
| **401** | API error |  -  |
| **409** | API error |  -  |
| **422** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## createProduct

> CreateProductResponse createProduct(createProductRequest)

Create product

### Example

```ts
import {
  Configuration,
  ProductsApi,
} from '';
import type { CreateProductOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ProductsApi();

  const body = {
    // CreateProductRequest
    createProductRequest: ...,
  } satisfies CreateProductOperationRequest;

  try {
    const data = await api.createProduct(body);
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
| **createProductRequest** | [CreateProductRequest](CreateProductRequest.md) |  | |

### Return type

[**CreateProductResponse**](CreateProductResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Product created |  -  |
| **409** | API error |  -  |
| **422** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## createProductAddon

> CreateProductAddonResponse createProductAddon(id, createProductAddonRequest)

Create product add-on

### Example

```ts
import {
  Configuration,
  ProductsApi,
} from '';
import type { CreateProductAddonOperationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ProductsApi();

  const body = {
    // string | Product UUID
    id: 38400000-8cf0-11bd-b23e-10b96e4ef00d,
    // CreateProductAddonRequest
    createProductAddonRequest: ...,
  } satisfies CreateProductAddonOperationRequest;

  try {
    const data = await api.createProductAddon(body);
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
| **createProductAddonRequest** | [CreateProductAddonRequest](CreateProductAddonRequest.md) |  | |

### Return type

[**CreateProductAddonResponse**](CreateProductAddonResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Add-on created |  -  |
| **401** | API error |  -  |
| **409** | API error |  -  |
| **422** | API error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

