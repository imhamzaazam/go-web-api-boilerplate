
# UserResponse


## Properties

Name | Type
------------ | -------------
`uid` | string
`fullName` | string
`email` | string

## Example

```typescript
import type { UserResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "uid": eb01f6d3-b964-4bdc-be66-a40c95378480,
  "fullName": John Doe,
  "email": john.doe@example.com,
} satisfies UserResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as UserResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


