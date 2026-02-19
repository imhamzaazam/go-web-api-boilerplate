
# CreateCartItemResponse


## Properties

Name | Type
------------ | -------------
`id` | string
`cartId` | string
`productId` | string
`quantity` | number
`unitPrice` | number

## Example

```typescript
import type { CreateCartItemResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "cartId": null,
  "productId": null,
  "quantity": null,
  "unitPrice": null,
} satisfies CreateCartItemResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as CreateCartItemResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


