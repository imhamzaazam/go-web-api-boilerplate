
# Go Web API Boilerplate

Just writing some golang with things that I like to see in a Web API

## Features

- Hexagonal Architecture (kinda overengineering but ok. Also, just wrote like this to see how it goes)
- Simple routing with chi
- Centralized encoding and decoding
- Centralized error handling
- Versioned HTTP Handler
- SQL type safety with SQLC
- Migrations with golang migrate
- PASETO tokens instead of JWT
- Access and Refresh Tokens
- Tests that uses Testcontainers instead of mocks
- Testing scripts that uses cURL and jq (f* Postman)

## Required dependencies

- jq
- golang-migrate
- docker
- sqlc

## OpenAPI (direct API calling)

- Spec file: `docs/openapi.yaml`
- Validate spec:

```bash
make openapi-check
```

- Start Swagger UI (Try it Out):

```bash
make openapi-ui
```

Open http://localhost:8081 and call endpoints directly from the browser.

- Generate API client SDK:

```bash
make openapi-generate-client
```

Generated client path: `generated/typescript-fetch`
