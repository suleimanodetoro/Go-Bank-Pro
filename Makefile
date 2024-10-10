# Start PostgreSQL container
# This target creates and starts a new PostgreSQL container using Docker.
# It maps port 5432 on the host to port 5432 in the container, sets the PostgreSQL root user,
# and the password is "secret". It uses the PostgreSQL 12 image with the Alpine Linux base for a lightweight setup.
postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

# Create the "simple_bank" database in the PostgreSQL container.
# This target runs the `createdb` command inside the PostgreSQL container to create a new database 
# called "simple_bank" with the root user as the owner.
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

# Drop the "simple_bank" database in the PostgreSQL container.
# This target runs the `dropdb` command inside the PostgreSQL container to drop the database
# called "simple_bank".
dropdb:
	docker exec -it postgres12 dropdb simple_bank

# Run database migrations (up direction).
# This target uses the `migrate` tool to apply all new migrations to the database.
# The migrations are located in the "db/migration" directory, and the target database is
# specified by a PostgreSQL connection URL.
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

# Rollback the last database migration (down direction).
# This target undoes the last applied migration using the `migrate` tool.
# It uses the same "db/migration" directory and PostgreSQL connection URL as the migrateup target.
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

# Generate SQL code using sqlc.
# This target runs the `sqlc` code generation tool to generate type-safe Go code
# from SQL queries defined in the "db/sqlc" directory.
sqlc:
	sqlc generate

# Run all Go tests with coverage reporting.
# This target runs `go test` in verbose mode (-v) and includes test coverage information (-cover)
# for all packages in the current directory and subdirectories (`./...`).
test:
	go test -v -cover ./...

server:
	go run main.go

# Declare all targets as "phony" to prevent make from confusing them with actual files.
# This ensures that Make will always run the associated commands when these targets are invoked.
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server
