
# CreateOrderRequest


## Properties

Name | Type
------------ | -------------
`cartId` | string
`paymentMethodId` | string
`fulfillmentType` | string
`location` | [CreateOrderLocation](CreateOrderLocation.md)

## Example

```typescript
import type { CreateOrderRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "cartId": null,
  "paymentMethodId": null,
  "fulfillmentType": null,
  "location": null,
} satisfies CreateOrderRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as CreateOrderRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


