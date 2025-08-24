mock:
	mockery
migrate:
	migrate -source file://postgres/migrations -database postgres://postgres:postgres@localhost:5432/twitter_clone_development?sslmode=disable up
rollback:
	migrate -source file://postgres/migrations -database postgres://postgres:postgres@localhost:5432/twitter_clone_development?sslmode=disable down
drop:
	migrate -source file://postgres/migrations -database postgres://postgres:postgres@localhost:5432/twitter_clone_development?sslmode=disable drop
migration:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir postgres/migrations $$name

run:
	go run cmd/graphqlserver/main.go

generate:
	go generate ./...
