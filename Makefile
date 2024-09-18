# Start PostgreSQL container
postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

# Create database
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

# Drop database
dropdb:
	docker exec -it postgres12 dropdb simple_bank

# Declare all targets as phony
.PHONY: postgres createdb dropdb

# Explanations:
# 1. The 'postgres' rule sets up a PostgreSQL 12 container with Alpine Linux base.
#    It sets the root user, password, and maps port 5432 to the host.
#
# 2. 'createdb' and 'dropdb' rules remain the same, executing commands in the container.
#
# 3. All targets (postgres, createdb, dropdb) are correctly declared as phony.
#
# 4. The -it flags are kept, which is fine for interactive use. For automated scripts,
#    you might want to remove them.
#
# 5. Consider adding a 'migrateup' and 'migratedown' rules if you plan to use
#    database migrations.