# migrations
migrate-up:
	@migrate -path=./migrations -database=postgres://veronez:261602317@localhost/greenlight?sslmode=disable up

migrate-down:
	@migrate -path=./migrations -database=postgres://veronez:261602317@localhost/greenlight?sslmode=disable down $(n)

run:
	@go run ./cmd/api
