
# CreateSubscriptionResponse


## Properties

Name | Type
------------ | -------------
`id` | string
`tenantId` | string
`plan` | string
`status` | string
`startsAt` | Date
`endsAt` | Date

## Example

```typescript
import type { CreateSubscriptionResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "tenantId": null,
  "plan": null,
  "status": null,
  "startsAt": null,
  "endsAt": null,
} satisfies CreateSubscriptionResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as CreateSubscriptionResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


