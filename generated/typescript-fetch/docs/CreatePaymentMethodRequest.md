
# CreatePaymentMethodRequest


## Properties

Name | Type
------------ | -------------
`type` | string
`label` | string
`isDefault` | boolean

## Example

```typescript
import type { CreatePaymentMethodRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "type": null,
  "label": null,
  "isDefault": null,
} satisfies CreatePaymentMethodRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as CreatePaymentMethodRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


