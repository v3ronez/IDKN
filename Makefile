# migrations
migrate-up: confirm
	@migrate -path=./migrations -database=postgres://veronez:261602317@localhost/greenlight?sslmode=disable up

migrate-down:
	@migrate -path=./migrations -database=postgres://veronez:261602317@localhost/greenlight?sslmode=disable down $(n)

run:
	@go run ./cmd/api

# confirm
confirm:
	@echo 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

# Audit
audit:
	@echo 'Tidying and verifying module dependencies...' go mod tidy
	@go mod verify
	@echo 'Formatting code...'
	@go fmt ./...
	@echo 'Vetting code...'
	@go vet ./...
	@staticcheck ./...
	@echo 'Running tests...'
	@go test -race -vet=off ./...


# Vendoring dependencies in this way basically stores a complete
# copy of the source code for third-party packages in a vendor folder in your project.
vendor:
	@echo 'Tidying and verifying module dependencies...'
	@go mod tidy
	@go mod verify
	@echo 'Vendoring dependencies...'
	@go mod vendor
