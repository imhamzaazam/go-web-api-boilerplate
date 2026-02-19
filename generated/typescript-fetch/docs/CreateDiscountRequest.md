
# CreateDiscountRequest


## Properties

Name | Type
------------ | -------------
`code` | string
`name` | string
`type` | string
`value` | number
`startsAt` | Date
`endsAt` | Date

## Example

```typescript
import type { CreateDiscountRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "code": null,
  "name": null,
  "type": null,
  "value": null,
  "startsAt": null,
  "endsAt": null,
} satisfies CreateDiscountRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as CreateDiscountRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


