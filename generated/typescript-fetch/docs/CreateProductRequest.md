
# CreateProductRequest


## Properties

Name | Type
------------ | -------------
`name` | string
`sku` | string
`price` | number
`vatPercent` | number
`isPreorder` | boolean
`madeToOrder` | boolean
`requiresPrescription` | boolean
`availableForDelivery` | boolean
`availableForPickup` | boolean

## Example

```typescript
import type { CreateProductRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "name": null,
  "sku": null,
  "price": null,
  "vatPercent": null,
  "isPreorder": null,
  "madeToOrder": null,
  "requiresPrescription": null,
  "availableForDelivery": null,
  "availableForPickup": null,
} satisfies CreateProductRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as CreateProductRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


