
# CreateDiscountResponse


## Properties

Name | Type
------------ | -------------
`id` | string
`productId` | string
`code` | string
`name` | string
`type` | string
`value` | number
`startsAt` | Date
`endsAt` | Date

## Example

```typescript
import type { CreateDiscountResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "productId": null,
  "code": null,
  "name": null,
  "type": null,
  "value": null,
  "startsAt": null,
  "endsAt": null,
} satisfies CreateDiscountResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as CreateDiscountResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


