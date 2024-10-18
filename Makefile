# Start PostgreSQL container
# Runs a new PostgreSQL 12 container with name "postgres12"
postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

# Create the "simple_bank" database
# Executes inside the PostgreSQL container to create the "simple_bank" database
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

# Drop the "simple_bank" database
# Executes inside the PostgreSQL container to drop the "simple_bank" database
dropdb:
	docker exec -it postgres12 dropdb simple_bank

# Run database migrations (up)
# Applies all migrations to bring the database to the latest version
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

# Apply only the next migration (up)
migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

# Rollback the last database migration (down)
# Rolls back all applied migrations
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

# Rollback only the last applied migration (down)
migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

# Generate SQL code using sqlc
# Generates Go code for SQL queries defined in "db/migration"
sqlc:
	sqlc generate

# Run all Go tests with coverage
# Runs all tests with verbose output and displays code coverage
test:
	go test -v -cover ./...

# Run the server
# Starts the main server application
server:
	go run main.go

# Generate mocks using mockgen
# Generates mock implementations for store interfaces, used for testing
mock:
	mockgen -source=db/sqlc/store.go \
        -destination=db/mock/store.go \
        -package=mockdb \
        -aux_files=github.com/suleimanodetoro/Go-Bank-Pro/db/sqlc=db/sqlc/querier.go
	go mod tidy

# Tidy up dependencies
# Cleans up any unused dependencies and ensures the "go.mod" file is up to date
tidy:
	go mod tidy

# Declare phony targets
# Indicates that these targets are not files, preventing conflicts if files with these names exist
.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock tidy
