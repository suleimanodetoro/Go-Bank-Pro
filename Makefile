# Start PostgreSQL container
postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

# Create the "simple_bank" database
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

# Drop the "simple_bank" database
dropdb:
	docker exec -it postgres12 dropdb simple_bank

# Run database migrations (up)
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

# Rollback the last database migration (down)
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

# Generate SQL code using sqlc
sqlc:
	sqlc generate

# Run all Go tests with coverage
test:
	go test -v -cover ./...

# Run the server
server:
	go run main.go

# Generate mocks using mockgen
mock:
	mockgen -source=db/sqlc/store.go \
        -destination=db/mock/store.go \
        -package=mockdb \
        -aux_files=github.com/suleimanodetoro/Go-Bank-Pro/db/sqlc=db/sqlc/querier.go
	go mod tidy

# Tidy up dependencies
tidy:
	go mod tidy

# Declare phony targets
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock tidy
