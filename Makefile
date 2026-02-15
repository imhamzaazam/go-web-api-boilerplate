run:
	go run cmd/http/main.go

build:
	@go build -o bin/go-boilerplate cmd/http/main.go

test:
	@go test -v ./...

sqlc:
	sqlc generate

migrateup:
	migrate -path db/postgres/migration -database "$(DB_SOURCE)" -verbose up

migrateup1:
	migrate -path db/postgres/migration -database "$(DB_SOURCE)" -verbose up 1

migratedown:
	migrate -path db/postgres/migration -database "$(DB_SOURCE)" -verbose down

migratedown1:
	migrate -path db/postgres/migration -database "$(DB_SOURCE)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/postgres/migration -seq $(name)

openapi-ui:
	docker run --rm -p 8081:8080 -e SWAGGER_JSON=/spec/openapi.yaml -v $(CURDIR)/docs:/spec swaggerapi/swagger-ui

openapi-check:
	docker run --rm -v $(CURDIR):/local openapitools/openapi-generator-cli validate -i /local/docs/openapi.yaml

openapi-generate-client:
	docker run --rm -v $(CURDIR):/local openapitools/openapi-generator-cli generate -i /local/docs/openapi.yaml -g typescript-fetch -o /local/generated/typescript-fetch

.PHONY:
	run build test sqlc migrateup migrateup1 migratedown migratedown1 new_migration openapi-ui openapi-check openapi-generate-client